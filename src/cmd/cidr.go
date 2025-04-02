package cmd

import (
	"encoding/json"
	"net"
	"os"
)

type HostData struct {
	Host string `json:"ips"`
}

func GetIPFromCIDIR(cidr string) ([]string, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	var ips []string
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); increaseIP(ip) {
		ips = append(ips, ip.String())
	}
	if len(ips) > 2 {
		ips = ips[1 : len(ips)-1]
	}
	return ips, nil
}

func increaseIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func SaveToConfig(ips []string) error {
	var data []HostData
	filePath := "config/ips.json"
	file, err := os.ReadFile(filePath)
	if err == nil {
		err = json.Unmarshal(file, &data)
		if err != nil {
			println("Error unmarshalling JSON:", err)
			return err
		}
	}
	for _, ip := range ips {
		data = append(data, HostData{Host: ip})
	}
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		println("Error marshalling JSON:", err)
		return err
	}
	return os.WriteFile("config/ips.json", jsonData, 0644)
}
