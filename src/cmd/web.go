package cmd

import (
	"fmt"
	"net/http"
)

func IsWebServer(ip string) {
	url := "http://" + ip
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	webServer := resp.Header.Get("Server")
	if webServer != "" {
		fmt.Printf("Web server found at %s: %s\n", ip, webServer)
	} else {
		fmt.Printf("No web server found at %s\n", ip)
	}
}
