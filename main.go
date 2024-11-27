package main

import (
	"log"

	"github.com/dimassfeb-09/pestapasta-be/controllers"
	"github.com/dimassfeb-09/pestapasta-be/models"
	"github.com/dimassfeb-09/pestapasta-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func init() {
	var err error
	// Initialize the DB
	db, err = models.InitializeDB()

	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
}

func main() {
	r := gin.Default()

	utils.Cors(r)

	r.POST("/login", func(c *gin.Context) {
		controllers.HandleLogin(c, db)
	})

	r.POST("/checkout", func(c *gin.Context) {
		controllers.HandleCheckout(c, db)
	})

	r.GET("/menus", func(c *gin.Context) {
		controllers.GetMenu(c, db)
	})

	r.GET("/categories", func(c *gin.Context) {
		controllers.GetCategories(c, db)
	})

	r.GET("/payment_methods", func(c *gin.Context) {
		controllers.GetPaymentMethods(c, db)
	})

	r.GET("/orders", func(c *gin.Context) {
		controllers.GetAllOrders(c, db)
	})

	r.GET("/orders/:id", func(c *gin.Context) {
		controllers.GetOrderByID(c, db)
	})

	// Start the server
	if err := r.Run(":8081"); err != nil {
		log.Fatal("Server failed:", err)
	}
}
