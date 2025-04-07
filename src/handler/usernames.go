package handler

import (
	"encoding/json"
	"fmt"
	"os"
)

type Username struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func LoadUsernames() ([]string, error) {
	file, err := os.Open("config/usernames.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	var data []Username
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}
	var usernames []string
	for _, username := range data {
		usernames = append(usernames, username.Username)
	}
	return usernames, nil
}
