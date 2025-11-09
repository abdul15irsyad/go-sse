package util

import (
	"fmt"
	"go-sse/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	host := config.DBHost
	user := config.DBUser
	password := config.DBPassword
	dbname := config.DBName
	port := config.DBPort

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB = database
	fmt.Println("database connected successfully")
}
