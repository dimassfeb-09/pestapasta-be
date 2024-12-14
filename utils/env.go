package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Email struct {
	User     string
	Password string
}

// ENV struct holds all environment variables, such as database and Midtrans keys
type ENV struct {
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	SSLMode      string
	MidtransKey  string
	SecretKeyJWT string
	Email        Email
}

// GetENV loads environment variables based on the current environment (production or local)
func GetENV() ENV {
	// Load environment variables from .env file (if present)
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Determine the current environment (production or local)
	env := getEnv("APP_ENV", "local") // Default to 'local' if not set

	// Initialize environment-specific variables
	var dbHost, dbPort, dbUser, dbPassword, dbName, sslMode, midtransKey, secretKeyJWT string
	var email Email

	// Check the environment and load the respective variables
	if env == "production" {
		dbHost = getEnv("DB_HOST_PRODUCTION", "prod-db-host")
		dbPort = getEnv("DB_PORT_PRODUCTION", "5432")
		dbUser = getEnv("DB_USER_PRODUCTION", "prod-user")
		dbPassword = getEnv("DB_PASSWORD_PRODUCTION", "prod-password")
		dbName = getEnv("DB_NAME_PRODUCTION", "prod-db")
		sslMode = getEnv("SSL_MODE_PRODUCTION", "require")
		midtransKey = getEnv("MIDTRANS_SERVER_KEY_PRODUCTION", "default-prod-key")
	} else {
		dbHost = getEnv("DB_HOST_LOCAL", "localhost")
		dbPort = getEnv("DB_PORT_LOCAL", "5432")
		dbUser = getEnv("DB_USER_LOCAL", "local-user")
		dbPassword = getEnv("DB_PASSWORD_LOCAL", "local-password")
		dbName = getEnv("DB_NAME_LOCAL", "local-db")
		sslMode = getEnv("SSL_MODE_LOCAL", "disable")
		midtransKey = getEnv("MIDTRANS_SERVER_KEY_SANDBOX", "default-sandbox-key")
	}

	email.User = getEnv("EMAIL_USER_MAILER", "default-email")
	email.Password = getEnv("EMAIL_PASSWORD_MAILER", "default-password")

	// Get the JWT secret key which is common across environments
	secretKeyJWT = getEnv("SECRET_KEY_JWT", "")

	// Return the populated ENV struct
	return ENV{
		DBHost:       dbHost,
		DBPort:       dbPort,
		DBUser:       dbUser,
		DBPassword:   dbPassword,
		DBName:       dbName,
		SSLMode:      sslMode,
		MidtransKey:  midtransKey,
		SecretKeyJWT: secretKeyJWT,
		Email:        email,
	}
}

// getEnv retrieves the value of an environment variable or returns a default value if not set
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
