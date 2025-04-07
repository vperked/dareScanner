package util

import (
	"encoding/json"
	"os"
)

type HostData struct {
	Host string `json:"ips"`
}

func GetConfig() []HostData {
	filePath := "config/ips.json"
	data, err := os.ReadFile(filePath)
	if err != nil {
		println("Error reading file:", err)
		return nil
	}
	var config []HostData
	err = json.Unmarshal(data, &config)
	if err != nil {
		println("Error unmarshalling JSON:", err)
		return nil
	}
	return config
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
