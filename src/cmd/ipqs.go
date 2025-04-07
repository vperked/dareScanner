package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type IPQSISAbuseResponse struct {
	Success bool `json:"success"`
	Data    struct {
		AbuseConfidenceScore int      `json:"abuse_confidence_score"`
		ISP                  string   `json:"isp"`
		Domain               string   `json:"domain"`
		Hostnames            []string `json:"hostnames"`
		TotalReports         int      `json:"total_reports"`
		NumDistinctUsers     int      `json:"num_distinct_users"`
		LastReportedAt       string   `json:"last_reported_at"`
		Organization         string   `json:"organization"`
		IsVPN                bool     `json:"is_vpn"`
	} `json:"data"`
}

type Cofnig struct {
	IPQSKey string `json:"IPQSKey"`
}

func IPQSISAbuse(ip string) {
	os, err := os.Open("config/config.json")
	if err != nil {
		println("Error opening config file:", err.Error())
		return
	}
	defer os.Close()
	var config Cofnig
	if err := json.NewDecoder(os).Decode(&config); err != nil {
		println("Error decoding config file:", err.Error())
		return
	}
	url := fmt.Sprintf("https://ipqualityscore.com/api/json/ip/%s", ip)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		println("Error creating request:", err.Error())
		return
	}
	req.Header.Add("Key", config.IPQSKey)
	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		println("Error making request:", err.Error())
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)
	if resp.StatusCode != 200 {
		println("Error: received non-200 response code")
		return
	}
	var response IPQSISAbuseResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		println("Error decoding response:", err.Error())
		return
	}
	responseData := response.Data
	fmt.Printf("IP: %s\n", ip)
	fmt.Printf("Abuse Confidence Score: %d\n", responseData.AbuseConfidenceScore)
	fmt.Printf("ISP: %s\n", responseData.ISP)
	fmt.Printf("Domain: %s\n", responseData.Domain)
	fmt.Printf("Hostnames: %v\n", responseData.Hostnames)
	fmt.Printf("Total Reports: %d\n", responseData.TotalReports)
	fmt.Printf("Num Distinct Users: %d\n", responseData.NumDistinctUsers)
	fmt.Printf("Last Reported At: %s\n", responseData.LastReportedAt)
	fmt.Printf("Organization: %s\n", responseData.Organization)
	fmt.Printf("Is VPN: %t\n", responseData.IsVPN)
	fmt.Println("IPQS ISAbuse check completed.")
}
