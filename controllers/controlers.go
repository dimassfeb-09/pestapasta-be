package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dimassfeb-09/pestapasta-be/models"
	"github.com/dimassfeb-09/pestapasta-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func HandleLogin(c *gin.Context, db *gorm.DB) {
	var loginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Parse JSON input
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	var user models.User
	if err := db.Where("username = ?", loginRequest.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Username or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user"})
		}
		return
	}

	// Verify the password (assuming it's hashed)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Username or password"})
		return
	}

	// Generate a JWT token (example token generation)
	token, err := utils.GenerateJWT(uint(user.ID), user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	// Respond with the token
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

func HandleCheckout(c *gin.Context, db *gorm.DB) {
	var checkoutRequest models.CheckoutRequest

	// Parse JSON input
	if err := c.ShouldBindJSON(&checkoutRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	// Validate input
	for _, product := range checkoutRequest.Products {
		if product.Quantity <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Quantity must be greater than 0"})
			return
		}
	}

	// Fetch products from DB
	var products []models.Menu
	productIDs := make([]int, len(checkoutRequest.Products))
	for i, item := range checkoutRequest.Products {
		productIDs[i] = item.ID
	}

	// Check if all products exist in the database
	if err := db.Where("id IN (?)", productIDs).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching products"})
		return
	}

	// Check if all requested products exist in the fetched products
	productMap := make(map[int]models.Menu)
	for _, product := range products {
		productMap[product.ID] = product
	}

	for _, item := range checkoutRequest.Products {
		if _, exists := productMap[item.ID]; !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Product with ID %d not found", item.ID)})
			return
		}
	}

	// Map quantities from request to products
	productQuantities := make(map[int]int)
	for _, item := range checkoutRequest.Products {
		productQuantities[item.ID] = item.Quantity
	}

	// Calculate total price
	total := 0
	for _, product := range products {
		quantity := productQuantities[product.ID]
		total += product.Price * quantity
	}

	// Create order
	order := models.Order{
		OrderDate:   time.Now().Format("2006-01-02 15:04:05"),
		TotalPrice:  total,
		OrderStatus: "Pending",
		Email:       checkoutRequest.Email,
		Name:        checkoutRequest.Name,
	}

	// Fetch payment method
	var paymentMethod models.PaymentMethod
	if err := db.Where("id = ?", checkoutRequest.PaymentMethodID).First(&paymentMethod).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment method not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payment method"})
		return
	}

	// Create order in database
	if err := db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating order"})
		return
	}

	// Create payment record
	payment := models.Payment{
		OrderID:              order.ID,
		PaymentMethod:        paymentMethod.Name,
		PaymentAccountNumber: paymentMethod.AccountNumber,
		PaymentStatus:        "pending",
		PaymentDate:          time.Now().Format("2006-01-02 15:04:05"),
		TransactionCode:      fmt.Sprintf("TXN%d", order.ID),
		CreatedAt:            time.Now().Format("2006-01-02 15:04:05"),
	}

	// Create order details
	for i, product := range products {
		quantity := productQuantities[product.ID]
		orderDetail := models.OrderDetail{
			OrderID:       order.ID,
			MenuID:        product.ID,
			Quantity:      quantity,
			Notes:         checkoutRequest.Products[i].Note,
			SubtotalPrice: product.Price * quantity,
		}
		if err := db.Create(&orderDetail).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating order details"})
			return
		}
	}

	// Create payment record in database
	if err := db.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating payment"})
		return
	}

	// Build response
	checkoutResponse := models.CheckoutResponse{
		Name:                 checkoutRequest.Name,
		Email:                checkoutRequest.Email,
		ProductDetails:       products,
		Total:                total,
		PaymentAccountNumber: payment.PaymentAccountNumber,
		PaymentMethod:        payment.PaymentMethod,
		PaymentStatus:        payment.PaymentStatus,
		TransactionCode:      payment.TransactionCode,
	}

	// Respond to client
	c.JSON(http.StatusOK, checkoutResponse)
}

func GetMenu(c *gin.Context, db *gorm.DB) {
	category := c.DefaultQuery("category", "")

	var menu []models.Menu

	if category != "" {
		// Fetch menu items with a specific category
		if err := db.Joins("JOIN categories ON categories.id = menus.category_id").
			Where("categories.category_name = ?", category).
			Find(&menu).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch menu items by category"})
			return
		}
	} else {
		// Fetch all menu items
		if err := db.Find(&menu).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch all menu items"})
			return
		}
	}

	c.JSON(http.StatusOK, menu)
}

func GetCategories(c *gin.Context, db *gorm.DB) {
	var categories []models.Category

	if err := db.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching all categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func GetPaymentMethods(c *gin.Context, db *gorm.DB) {
	var paymentMethods []models.PaymentMethod

	if err := db.Find(&paymentMethods).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching all payment methods"})
		return
	}

	c.JSON(http.StatusOK, paymentMethods)
}

func GetAllOrders(c *gin.Context, db *gorm.DB) {
	var orders []models.Order

	// Preload relasi ke User dan OrderDetails
	if err := db.Preload("Payment").Preload("OrderDetails.Menu").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch orders",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func GetOrderByID(c *gin.Context, db *gorm.DB) {
	orderID := c.Param("id")
	var order models.Order
	// Preload relasi ke User dan OrderDetails
	if err := db.Preload("OrderDetails").First(&order, "id = ?", orderID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch order",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, order)
}
