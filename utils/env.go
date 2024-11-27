package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ENV struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	SSLMode    string
}

func GetENV() ENV {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	return ENV{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		SSLMode:    os.Getenv("SSL_MODE"),
	}
}
