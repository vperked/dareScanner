package cmd

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"playGO/util"
	"time"
)

type Data struct {
	Host string `json:"ips"`
}

var ports = []int{22, 80, 443, 8080, 3306, 5432}

func GetConfig() []Data {
	filePath := "config/ips.json"
	data, err := os.ReadFile(filePath)
	if err != nil {
		println("Error reading file:", err)
		return nil
	}
	var config []Data
	err = json.Unmarshal(data, &config)
	if err != nil {
		println("Error unmarshalling JSON:", err)
		return nil
	}
	return config
}

func ScanPort() bool {
	config := GetConfig()
	if config == nil {
		println("No configuration found")
		return false
	}
	for _, entry := range config {
		ip := entry.Host
		for _, port := range ports {
			addr := net.JoinHostPort(ip, fmt.Sprintf("%d", port))
			conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
			if err != nil {
				println("Port", port, "is closed on", ip)
			} else {
				println("Port", port, "is open on", ip)
				conn.Close()
				if err := util.InsertOpenPort(ip, port); err != nil {
					println("Error inserting open port:", err)
				}
			}
		}
	}
	return true
}
