package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	Port int

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTPrivateKey string
	JWTPublicKey  string
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	fmt.Println("env loaded")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8765"
	}
	Port, _ = strconv.Atoi(port)

	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	DBUser = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASS")
	DBName = os.Getenv("DB_NAME")

	JWTPrivateKey = os.Getenv("JWT_PRIVATE_KEY")
	JWTPublicKey = os.Getenv("JWT_PUBLIC_KEY")
}
