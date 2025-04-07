package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"playGO/util"
	"time"

	"github.com/chromedp/chromedp"
)

type WebhookJsonPayload struct {
	Content string `json:"content"`
}

func SendWebhookMessage(webhookURL string, message string) error {
	payload := WebhookJsonPayload{Content: message}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func IsWebServer(ipPort string, ipNoPort string) {
	url := "http://" + ipPort
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	webServer := resp.Header.Get("Server")
	if webServer != "" {
		hostname, err := GetHostname(ipNoPort)
		if err != nil {
			fmt.Println(err)
			return
		}
		message := fmt.Sprintf("Web server found at %s (%s): %s", ipPort, hostname, webServer)
		util.AddChecked(ipNoPort)
		SendWebhookMessage("https://discord.com/api/webhooks/1356857871215890452/TUSPalcrGvLv6urWFtTM4mbxHNR34wYeMPwu40nmZjxz3_elHiIlboGfvafO5Ng4OMMm", message)
		fmt.Printf("Web server found at %s: %s\n", ipPort, webServer)
	} else {
		fmt.Printf("No web server found at %s\n", ipPort)
	}
	httpsURL := "https://" + ipPort
	resp, err = http.Get(httpsURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	webServer = resp.Header.Get("Server")
	hostname, err := GetHostname(ipNoPort)
	if err != nil {
		fmt.Println(err)
		return
	}
	if webServer != "" {
		fmt.Printf("Web server found at %s: %s\n", ipPort, webServer)
		TakeScreenshot(httpsURL, hostname)
	} else {
		fmt.Printf("No web server found at %s\n", ipPort)
	}
}

func TakeScreenshot(url string, hostname string) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var buf []byte
	err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.FullScreenshot(&buf, 90),
	})
	if err != nil {
		fmt.Println("Failed to take screenshot:", err)
		return
	}
	filename := fmt.Sprintf("%s.png", hostname)
	err = os.WriteFile(filename, buf, 0644)
	if err != nil {
		fmt.Println("Failed to save screenshot:", err)
		return
	}
	fmt.Println("Screenshot saved as screenshot.png")
}
