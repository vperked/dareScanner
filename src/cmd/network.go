package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"playGO/util"
)

type RipeResp struct {
	Data struct {
		Prefixes []struct {
			Prefix string `json:"prefix"`
		} `json:"prefixes"`
	} `json:"data"`
}

func GetHostname(ip string) (string, error) {
	parseIP := net.ParseIP(ip)
	if parseIP == nil || parseIP.To4() == nil {
		return "", fmt.Errorf("invalid IP address: %s", ip)
	}
	names, err := net.LookupAddr(ip)
	if err != nil {
		return "", err
	}
	return names[0], nil
}

func GetIPFromASN(asn string) ([]string, error) {
	url := fmt.Sprintf("https://stat.ripe.net/data/announced-prefixes/data.json?resource=%s", asn)
	resp, err := http.Get(url)
	if err != nil {
		println("Error fetching data from RIPE:", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		println("Error reading response body:", err.Error())
		return nil, err
	}
	var ripeResp RipeResp
	err = json.Unmarshal(body, &ripeResp)
	if err != nil {
		println("Error unmarshalling JSON:", err.Error())
		return nil, err
	}
	var allIPs []string
	for _, prefix := range ripeResp.Data.Prefixes {
		ips, err := GetIPFromCIDIR(prefix.Prefix, 10)
		if err != nil {
			println("Error getting IPs from prefix:", err.Error())
			return nil, err
		}
		fmt.Println(ips)
		for _, ip := range ips {
			err := util.InsertIP(ip)
			if err != nil {
				println("Error inserting IP:", err.Error())
				return nil, err
			}
		}
		allIPs = append(allIPs, ips...)
	}
	return allIPs, nil
}

func GetIPFromCIDIR(cidr string, limit int) ([]string, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	ip = ip.To4()
	if ip == nil {
		fmt.Printf("CIDR %s is not IPv4\n", cidr)
		return nil, nil
	}
	var ips []string
	if len(ip) != net.IPv4len {
		fmt.Printf("CIDR %s is not IPv4\n", cidr)
		return nil, nil
	}
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); increaseIP(ip) {
		ips = append(ips, ip.String())
		if len(ips) >= limit {
			break
		}
	}
	return ips, nil
}

func RandomASN(count int) []string {
	// Generate random ASN numbers
	asn := make([]string, count)
	for i := 0; i < count; i++ {
		randomASN := fmt.Sprintf("%d", rand.Intn(100000))
		ips, err := GetIPFromASN(randomASN)
		if err != nil {
			println("Error getting IPs from ASN:", err.Error())
			continue
		}
		if len(ips) > 0 {
			asn[i] = randomASN
			break
		} else {
			fmt.Printf("No IPs found for ASN %s, retrying...\n", randomASN)
		}
	}
	return asn
}

func increaseIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
