package cmd

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
)

type AbuseIPDBResponse struct {
	Data struct {
		IPAddress            string   `json:"ipAddress"`
		IsPublic             bool     `json:"isPublic"`
		AbuseConfidenceScore int      `json:"abuseConfidenceScore"`
		UsageType            string   `json:"usageType"`
		ISP                  string   `json:"isp"`
		Domain               string   `json:"domain"`
		Hostnames            []string `json:"hostnames"`
		TotalReports         int      `json:"totalReports"`
		NumDistinctUsers     int      `json:"numDistinctUsers"`
		LastReportedAt       string   `json:"lastReportedAt"`
	} `json:"data"`
}

type Config struct {
	Key string `json:"AbuseKey"`
}

func IsAbuse(ip string) {
	os, err := os.Open("config/config.json")
	if err != nil {
		println("Error opening config file:", err.Error())
		return
	}
	defer os.Close()
	var config Config
	if err := json.NewDecoder(os).Decode(&config); err != nil {
		println("Error decoding config file:", err.Error())
		return
	}
	baseURL := "https://api.abuseipdb.com/api/v2/check"
	params := url.Values{}
	params.Set("ipAddress", ip)
	params.Set("maxAgeInDays", "90")
	params.Set("verbose", "true")
	fullURL := baseURL + "?" + params.Encode()
	println(fullURL)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		println("Error creating request:", err.Error())
		return
	}
	req.Header.Add("Key", config.Key)
	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		println("Error making request:", err.Error())
		return
	}
	if resp.StatusCode != 200 {
		println("Error: received non-200 response code:", resp.StatusCode)
		return
	}
	defer resp.Body.Close()
	var response AbuseIPDBResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		println("Error decoding response:", err.Error())
		return
	}
	println("IP Address:", response.Data.IPAddress)
	println("Abuse Confidence Score:", response.Data.AbuseConfidenceScore)
	println("ISP:", response.Data.ISP)
	println("Domain:", response.Data.Domain)
	println("Total Reports:", response.Data.TotalReports)
	println("Num Distinct Users:", response.Data.NumDistinctUsers)

}
