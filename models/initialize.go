package models

import (
	"fmt"
	"log"

	"github.com/dimassfeb-09/pestapasta-be/utils"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// InitializeDB initializes the database connection
func InitializeDB() (*gorm.DB, error) {

	env := utils.GetENV()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", env.DBHost, env.DBUser, env.DBPassword, env.DBName, env.DBPort, env.SSLMode)

	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
		log.Fatal("failed to connect to the database")
		return nil, err
	}

	if err := db.AutoMigrate(&User{}, &Category{}, &Menu{}, &Order{}, &OrderDetail{}, &Payment{}, &PaymentMethod{}).Error; err != nil {
		log.Fatal("failed to migrate the database")
		return nil, err
	}

	fmt.Println("Successfully connected to the database")

	return db, nil
}
