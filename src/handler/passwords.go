package handler

import (
	"encoding/json"
	"fmt"
	"os"
)

type Password struct {
	ID       int    `json:"id"`
	Password string `json:"passwords"`
}

func LoadPasswords() ([]string, error) {
	file, err := os.Open("config/passwords.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	var data []Password
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}
	var passwords []string
	for _, password := range data {
		passwords = append(passwords, password.Password)
	}
	return passwords, nil
}
