package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"fullRssAPI/logs"
)

func connectToDB(host string, port int, user string, password string, dbName string) (*sql.DB, error) {
	logs.Info("hogehgoe")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	logs.Info("Connected to DB %s:%d/%s...", host, port, dbName)
	return db, nil
}
