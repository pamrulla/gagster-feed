package database

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB_USERNAME = "root"
var DB_PASSWORD = "root"
var DB_NAME = "gagster"
var DB_HOST = "127.0.0.1"
var DB_PORT = "3306"

var DB_CONNNAME = ""

var Db *gorm.DB

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func InitDb() *gorm.DB {
	DB_USERNAME = getEnv("DB_USERNAME", DB_USERNAME)
	DB_PASSWORD = getEnv("DB_PASSWORD", DB_PASSWORD)
	DB_CONNNAME = getEnv("DB_CONNNAME", DB_CONNNAME)
	Db = connectDB()
	return Db
}

func connectDB() *gorm.DB {
	var err error
	dsn := ""
	if len(DB_CONNNAME) == 0 {
		dsn = DB_USERNAME + ":" + DB_PASSWORD + "@tcp" + "(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME + "?" + "parseTime=true&loc=Local"
	} else {
		dsn = DB_USERNAME + ":" + DB_PASSWORD + "@unix(/cloudsql/" + DB_CONNNAME + ")/" + DB_NAME + "?" + "parseTime=true&loc=Local"
	}
	fmt.Println("dsn : ", dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Printf("Error connecting to database : error=%v\n", err)
		return nil
	}

	return db
}
