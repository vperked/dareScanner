package pkg

import (
	"fmt"
	"net"
	"playGO/cmd"
	"playGO/util"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

var ports = []int{22, 80, 443, 8080, 3306, 5432, 21, 25, 110, 135, 139, 445, 3389, 5900, 6379, 27017, 5000, 8000, 9000}

func ScanPort(ip string) {
	for _, port := range ports {
		addr := net.JoinHostPort(ip, fmt.Sprintf("%d", port))
		conn, err := net.DialTimeout("tcp", addr, 2)
		if err != nil {
			println("Port", port, "is closed on", ip)
		} else {
			webhook := "https://discord.com/api/webhooks/1356857871215890452/TUSPalcrGvLv6urWFtTM4mbxHNR34wYeMPwu40nmZjxz3_elHiIlboGfvafO5Ng4OMMm"
			cmd.SendWebhookMessage(webhook, fmt.Sprintf("Port %d is open on %s", port, ip))
			println("Port", port, "is open on", ip)
			conn.Close()
			if err := util.InsertOpenPort(ip, port); err != nil {
				println("Error inserting open port:", err)
			}
		}
	}
}

func ScanWebServer() {
	records, err := util.GetWebServerFromDB()
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Println("Open Ports:")
	for _, record := range records {
		ip := record["ip_address"].(string)
		port := record["port"].(int)
		url := fmt.Sprintf("%s:%d", ip, port)
		println("Scanning web server:", url)
		db, err := util.IsIPChecked(ip)
		if err != nil {
			println("Error checking IP:", err.Error())
			return
		}
		if db {
			println("IP is already checked:", ip)
			continue
		}
		util.AddChecked(ip)
		cmd.IsWebServer(url, ip)
	}
}

type ConfigData struct {
	Cidrs []string `json:"cidrBlocks"`
}

func Scanner(asn int) {
	asns := util.RandomASN(asn)
	var allIPs []string
	for _, asn := range asns {
		err := util.InsertASN(asn, "")
		if err != nil {
			println("Error inserting ASN:", err.Error())
			return
		}
		fmt.Printf("Inserting ASN %s into the database\n", asn)
		ips, err := cmd.GetIPFromASN(asn)
		if err != nil {
			println("Error getting IPs from ASN:", err.Error())
			return
		}
		fmt.Printf("IPs from ASN %s: %v\n", asn, ips)
		allIPs = append(allIPs, ips...)

	}
	var wg sync.WaitGroup
	for _, ip := range allIPs {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			ScanPort(ip)
		}(ip)
		cpuPercent, err := cpu.Percent(1*time.Second, false)
		if err != nil {
			println("Couldnt fetch cpu:", err.Error())
			return
		}
		if cpuPercent[0] > 80 {
			println("CPU usage is high:", fmt.Sprintf("%.2f%%", cpuPercent[0]))
			time.Sleep(3 * time.Second)
			println("Sleeping for 3 seconds")
		} else {
			println("CPU usage is normal:", fmt.Sprintf("%.2f%%", cpuPercent[0]))
		}
	}
	wg.Wait()
	fmt.Println("All goroutines finished")
}

func ScannerIPsInDB() {
	var wg sync.WaitGroup
	allIPs, err := util.GetAllIPsFromDB()
	if err != nil {
		println("Error getting IPs from database:", err.Error())
		return
	}
	for _, ip := range allIPs {
		isChecked, err := util.IsIPCheckedInIps(ip)
		if err != nil {
			println("Error checking IP:", err.Error())
			continue
		}
		if isChecked {
			println("IP is already checked:", ip)
			continue
		}

		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			util.AddCheckedIPS(ip)
			ScanPort(ip)
		}(ip)

		cpuPercent, err := cpu.Percent(1*time.Second, false)
		if err != nil {
			println("Could not fetch CPU usage:", err.Error())
			continue
		}
		if cpuPercent[0] > 80 {
			println("CPU usage is high:", fmt.Sprintf("%.2f%%", cpuPercent[0]))
			time.Sleep(3 * time.Second)
			println("Sleeping for 3 seconds")
		} else {
			println("CPU usage is normal:", fmt.Sprintf("%.2f%%", cpuPercent[0]))
		}
	}
	wg.Wait()
	fmt.Println("All goroutines finished")
}
