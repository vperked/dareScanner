package main

import (
	"encoding/json"
	"fmt"
	"os"
	"playGO/cmd"
	"sync"
	"time"
)

func parse() {
	type ConfigData struct {
		Cidrs []string `json:"cidrBlocks"`
	}
	filepath := "config/cidrlist.json"
	fileReader, err := os.ReadFile(filepath)
	if err != nil {
		println("Error reading file:", err.Error())
		return
	}
	var configData ConfigData
	err = json.Unmarshal(fileReader, &configData)
	if err != nil {
		println("Error unmarshalling JSON:", err.Error())
		return
	}
	var allIps []string
	for _, c := range configData.Cidrs {
		ips, err := cmd.GetIPFromCIDIR(c)
		if err != nil {
			println("Error getting IPs from CIDR:", err.Error())
			return
		}
		allIps = append(allIps, ips...)
	}
	if len(allIps) == 0 {
		println("No IPs found")
		return
	}
	err = cmd.SaveToConfig(allIps)
	if err != nil {
		println("Error saving to config:", err.Error())
		return
	}
	println("IPs saved to config")
}

func main() {
	var answer string
	println("Do you want to parse the CIDR list? (y/n)")
	fmt.Scanln(&answer)
	if answer == "y" {
		parse()
	}
	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cmd.ScanPort()
			println("Port scan completed")
		}(i)
	}
	wg.Wait()
	println("All goroutines finished")
	final := time.Since(start)
	fmt.Printf("Execution time: %s \n", final)
}
