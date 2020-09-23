package config

import (
	"os"
	"strconv"
)

var DBPort, _ = strconv.Atoi(os.Getenv("POSTGRES_PORT"))
var DBContainer string = os.Getenv("DB_CONTAINER")
var DBName string = os.Getenv("POSTGRES_DB")
var DBUser string = os.Getenv("POSTGRES_USER")
var DBPassword string = os.Getenv("POSTGRES_PASSWORD")

var DBInfo = struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}{
	Host:     DBContainer,
	Port:     DBPort,
	User:     DBName,
	Password: DBPassword,
	DBName:   DBUser,
}
