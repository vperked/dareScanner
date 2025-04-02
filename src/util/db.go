package util

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectToDB() (*sql.DB, error) {
	dbInf := "root:Perkedishotfr1.@tcp(localhost:3306)/networkScanner"
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

func InsertOpenPort(ip string, port int) error {
	db, err := ConnectToDB()
	if err != nil {
		println("Error connecting to the database:", err.Error())
		return err
	}
	defer db.Close()

	// Create the table if it doesn't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS open_ports (id INT AUTO_INCREMENT PRIMARY KEY, ip_address VARCHAR(15), port INT)")
	if err != nil {
		println("Error creating table:", err.Error())
	}

	statement, err := db.Prepare("INSERT IGNORE INTO open_ports (ip_address, port) VALUES (?, ?)")
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
