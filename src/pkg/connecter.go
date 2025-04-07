package pkg

import (
	"database/sql"
	"fmt"
)

func ConnectToDB(username string, password string, ip string, database string) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, ip, 3306, database))
	if err != nil {
		println("Error connecting to the database:", err.Error())
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		println("Error pinging the database:", err.Error())
		return nil, err
	}
	println("Connected to the database successfully")
	return db, nil

}

func MultiDBConnection(usernames []string, passwords []string, ips []string, databases []string) ([]*sql.DB, error) {
	var dbs []*sql.DB
	for i := 0; i < len(usernames); i++ {
		db, err := ConnectToDB(usernames[i], passwords[i], ips[i], databases[i])
		if err != nil {
			println("Error connecting to the database:", err.Error())
			return nil, err
		}
		db.Close()
		dbs = append(dbs, db)
	}
	println("Connected to all databases successfully")
	return dbs, nil
}
