package util

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
}

func CreateASNTable() error {
	db, err := ConnectToDB()
	if err != nil {
		println("Error connecting to the database:", err.Error())
		return err
	}
	defer db.Close()

	_, err = db.Prepare("CREATE TABLE IF NOT EXISTS asn (id INT AUTO_INCREMENT PRIMARY KEY, asn VARCHAR(10), ip_address VARCHAR(15))")
	if err != nil {
		println("Error creating table:", err.Error())
		return err
	}
	println("Table created successfully")
	return nil
}

func InsertASN(asn string, ip string) error {
	db, err := ConnectToDB()
	if err != nil {
		println("Error connecting to the database:", err.Error())
		return err
	}
	defer db.Close()
	CreateASNTable()
	var in bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM asn WHERE asn = ?)", asn).Scan(&in)
	if err != nil {
		println("Error checking ASN existence:", err.Error())
		return err
	}
	if in {
		println("ASN already exists:", asn)
		return nil
	}
	statement, err := db.Prepare("INSERT IGNORE INTO asn (asn, ip_address) VALUES (?, ?)")
	if err != nil {
		println("Error preparing statement:", err.Error())
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(asn, ip)
	if err != nil {
		println("Error executing statement:", err.Error())
		return err
	}
	println("Inserted ASN:", asn, ip)
	return nil
}

func ConnectToDB() (*sql.DB, error) {
	filePath := "config/db.json"
	file, err := os.ReadFile(filePath)
	if err != nil {
		println("Error reading file:", err.Error())
		return nil, err
	}
	var dbConfig DBConfig
	err = json.Unmarshal(file, &dbConfig)
	if err != nil {
		println("Error unmarshalling JSON:", err.Error())
		return nil, err
	}
	dbInf := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)
	db, err := sql.Open("mysql", dbInf)
	if err != nil {
		println("Error connecting to the database:", err.Error())
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		println("Error pinging the database:", err.Error())
		return nil, err
	}
	return db, nil
}

func InsertIP(ip string) error {
	db, err := ConnectToDB()
	if err != nil {
		println("Error connecting to the database:", err.Error())
		return err
	}
	defer db.Close()

	statement, err := db.Prepare("INSERT IGNORE INTO ips (ip_address) VALUES (?)")
	if err != nil {
		println("Error preparing statement:", err.Error())
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(ip)
	if err != nil {
		println("Error executing statement:", err.Error())
		return err
	}
	println("Inserted IP:", ip)
	return nil
}

func InsertOpenPort(ip string, port int) error {
	db, err := ConnectToDB()
	if err != nil {
		println("Error connecting to the database:", err.Error())
		return err
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS ips (id INT AUTO_INCREMENT PRIMARY KEY, ip_address VARCHAR(15), port INT)")
	if err != nil {
		println("Error creating table:", err.Error())
	}

	statement, err := db.Prepare("INSERT IGNORE INTO ips (ip_address, port) VALUES (?, ?)")
	if err != nil {
		println("Error preparing statement:", err.Error())
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(ip, port)
	if err != nil {
		println("Error executing statement:", err.Error())
		return err
	}
	println("Inserted open port:", ip, port)
	return nil
}

func GetWebServerFromDB() ([]map[string]interface{}, error) {
	db, err := ConnectToDB()
	if err != nil {
		println("Error connecting to the database:", err.Error())
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT ip_address, port FROM ips WHERE port IN (443, 80)")
	if err != nil {
		println("Error querying database:", err.Error())
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var ip string
		var port int
		err := rows.Scan(&ip, &port)
		if err != nil {
			println("Error scanning row:", err.Error())
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"ip_address": ip,
			"port":       port,
		})
	}
	return results, nil
}

func GetFTPServerFromDB() ([]map[string]interface{}, error) {
	db, err := ConnectToDB()
	if err != nil {
		println("Error connecting to the database:", err.Error())
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT ip_address, port FROM ips WHERE port IN (21)")
	if err != nil {
		println("Error querying database:", err.Error())
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var ip string
		var port int
		err := rows.Scan(&ip, &port)
		if err != nil {
			println("Error scanning row:", err.Error())
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"ip_address": ip,
			"port":       port,
		})
	}
	return results, nil
}

func GetAllIPsFromDB() ([]string, error) {
	db, err := ConnectToDB()
	if err != nil {
		println("Error connecting to the database:", err.Error())
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT ip_address FROM ips")
	if err != nil {
		println("Error querying database:", err.Error())
		return nil, err
	}
	defer rows.Close()

	var ips []string
	for rows.Next() {
		var ip string
		err := rows.Scan(&ip)
		if err != nil {
			println("Error scanning row:", err.Error())
			return nil, err
		}
		ips = append(ips, ip)
	}
	return ips, nil
}

func GetSSHIPFromDB() ([]string, error) {
	db, err := ConnectToDB()
	if err != nil {
		println("Error connecting to the database:", err.Error())
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT ip_address FROM ips WHERE port = 22")
	if err != nil {
		println("Error querying database:", err.Error())
		return nil, err
	}
	defer rows.Close()

	var ips []string
	for rows.Next() {
		var ip string
		err := rows.Scan(&ip)
		if err != nil {
			println("Error scanning row:", err.Error())
			return nil, err
		}
		ips = append(ips, ip)
	}
	return ips, nil
}

func IsIPChecked(ip string) (bool, error) {
	db, err := ConnectToDB()
	if err != nil {
		println("Error connecting to the database:", err.Error())
		return false, err
	}
	defer db.Close()
	var checked bool
	err = db.QueryRow("SELECT checked FROM ips WHERE ip_address = ?", ip).Scan(&checked)
	if err != nil {
		println("Error querying database:", err.Error())
		return false, err
	}
	return checked, nil
}

func IsIPCheckedInIps(ip string) (bool, error) {
	db, err := ConnectToDB()
	if err != nil {
		println("Error connecting to the database:", err.Error())
		return false, err
	}
	defer db.Close()
	var checked bool
	err = db.QueryRow("SELECT checked FROM ips WHERE ip_address = ?", ip).Scan(&checked)
	if err != nil {
		println("Error querying database:", err.Error())
		return false, err
	}
	return checked, nil
}

func AddChecked(ip string) ([]string, error) {
	db, err := ConnectToDB()
	if err != nil {
		println("Error connecting to the database:", err.Error())
		return nil, err
	}
	defer db.Close()
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS ips (id INT AUTO_INCREMENT PRIMARY KEY, ip_address VARCHAR(15), port INT, checked BOOLEAN DEFAULT false)")
	if err != nil {
		println("Error creating table:", err.Error())
		return nil, err
	}
	query := "UPDATE ips SET checked = true WHERE ip_address = ?"
	results, err := db.Exec(query, ip)
	if err != nil {
		println("Error updating database:", err.Error())
		return nil, err
	}
	rowsChanged, err := results.RowsAffected()
	if err != nil {
		println("Error getting rows affected:", err.Error())
		return nil, err
	}
	println("Rows changed on ips:", rowsChanged)
	return nil, nil
}
func AddCheckedIPS(ip string) ([]string, error) {
	db, err := ConnectToDB()
	if err != nil {
		println("Error connecting to the database:", err.Error())
		return nil, err
	}
	defer db.Close()
	query := "UPDATE ips SET checked = true WHERE ip_address = ?"
	results, err := db.Exec(query, ip)
	if err != nil {
		println("Error updating database:", err.Error())
		return nil, err
	}
	rowsChanged, err := results.RowsAffected()
	if err != nil {
		println("Error getting rows affected:", err.Error())
		return nil, err
	}
	println("Rows changed on ips:", rowsChanged)
	return nil, nil
}
func AddCheckedSSH(ip string) ([]string, error) {
	db, err := ConnectToDB()
	if err != nil {
		println("Error connecting to the database:", err.Error())
		return nil, err
	}
	defer db.Close()
	query := "UPDATE ips SET checked = true WHERE ip_address = ? AND port = 22"
	results, err := db.Exec(query, ip)
	if err != nil {
		println("Error updating database:", err.Error())
		return nil, err
	}
	rowsChanged, err := results.RowsAffected()
	if err != nil {
		println("Error getting rows affected:", err.Error())
		return nil, err
	}
	println("Rows changed on ips:", rowsChanged)
	return nil, nil
}
