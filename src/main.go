package main

import (
	"fmt"
	"playGO/pkg"
)

func main() {
	var answer string
	fmt.Println("Welcome to the Scanner, to scan a ASN enter: pScan")
	fmt.Println("Else, you can enter: ftp, SSH, Web to scan all services.")
	fmt.Scanln(&answer)
	switch answer {
	case "pscan":
		println("Starting the scan...")
		println("How many asns are you scanning?")
		var asn int
		fmt.Scanln(&asn)
		pkg.Scanner(asn)
	case "ftp":
		println("Starting the FTP scan...")
		pkg.ConnectToFTP()
	case "ssh":
		println("Starting the SSH scan...")
		pkg.ConnectToServer(22)
	case "web":
		println("Starting the Web scan...")
		pkg.ScanWebServer()
	case "pscanDB":
		pkg.ScannerIPsInDB()
	case "no":
		println("Scan aborted.")
	default:
		println("Invalid input. Please enter 'yes' or 'no'.")
	}
}
