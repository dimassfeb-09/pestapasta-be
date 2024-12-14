package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dimassfeb-09/pestapasta-be/helpers"
	"github.com/dimassfeb-09/pestapasta-be/models"
	"github.com/dimassfeb-09/pestapasta-be/services"
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
	total := 0.0
	for _, product := range products {
		quantity := productQuantities[product.ID]
		total += product.Price * float64(quantity)
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

	var midtransResponse *models.CreateTransactionMidtransResponse
	if paymentMethod.Code == "qris" {
		var itemDetails []models.ItemDetails
		for _, item := range checkoutRequest.Products {
			product := productMap[item.ID]
			itemDetails = append(itemDetails, models.ItemDetails{
				ID:       fmt.Sprintf("PRODUCTID-%d", product.ID),
				Price:    product.Price,
				Quantity: item.Quantity,
				Name:     product.Name,
			})
		}

		midtransPayload := models.CreateTransactionMidtransPayload{
			PaymentType: "qris",
			ItemDetails: itemDetails,
			TransactionDetails: struct {
				OrderID     string  `json:"order_id"`
				GrossAmount float64 `json:"gross_amount"`
			}{
				OrderID:     fmt.Sprintf("ORDER-%d", time.Now().Unix()), // Generate dynamic OrderID
				GrossAmount: total,
			},
			CustomerDetails: struct {
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Email     string `json:"email"`
			}{
				FirstName: checkoutRequest.Name,
				LastName:  checkoutRequest.Name,
				Email:     checkoutRequest.Email,
			},
			QRIS: struct {
				Acquirer string `json:"acquirer"`
			}{
				Acquirer: "gopay", // Set acquirer to 'gopay', or whatever is needed
			},
		}

		// Create transaction with Midtrans API
		responseCreateTransactionMidtrans, errorResponse, err := services.CreateTransaction(midtransPayload)
		if errorResponse != nil || err != nil {
			statusCode, _ := strconv.Atoi(errorResponse.StatusCode)
			c.JSON(statusCode, gin.H{"error": "Internal Server Error: Payment"})
			return
		}

		midtransResponse = responseCreateTransactionMidtrans
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
		PaymentAccountName:   paymentMethod.AccountName,
		PaymentStatus:        "pending",
		PaymentCreateDate:    time.Now().Format("2006-01-02 15:04:05"),
		PaymentExpiredDate:   midtransResponse.ExpiryTime,
		TransactionCode:      fmt.Sprintf("TXN%d", order.ID),
	}
	if payment.PaymentMethod == "QRIS" {
		payment.PaymentQRCodeURL = midtransResponse.Actions[0].URL
		payment.PaymentTransactionID = midtransResponse.TransactionID
	}

	// Create order details
	for i, product := range products {
		quantity := productQuantities[product.ID]

		orderDetail := models.OrderDetail{
			OrderID:       order.ID,
			MenuID:        product.ID,
			Quantity:      quantity,
			Notes:         checkoutRequest.Products[i].Notes,
			SubtotalPrice: product.Price * float64(quantity),
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

	paymentDetails := models.PaymentDetails{
		PaymentAccountNumber: payment.PaymentAccountNumber,
		PaymentAccountName:   payment.PaymentAccountName,
		PaymentMethod:        payment.PaymentMethod,
		PaymentStatus:        payment.PaymentStatus,
	}

	// Set the PaymentExpiredTime based on the payment method
	if payment.PaymentMethod == "BCA" {
		paymentDetails.PaymentExpiredTime = (time.Minute * 10).Milliseconds()
		paymentDetails.PaymentMethod = "bank"
	} else if payment.PaymentMethod == "QRIS" {
		paymentDetails.PaymentExpiredTime = (time.Minute * 15).Milliseconds()
		paymentDetails.PaymentMethod = "qris"
	}

	if midtransResponse != nil {
		paymentDetails.QRImageURL = midtransResponse.Actions[0].URL
	}

	// Build response
	checkoutResponse := models.CheckoutResponse{
		Name:            checkoutRequest.Name,
		Email:           checkoutRequest.Email,
		ProductDetails:  products,
		Total:           total,
		PaymentDetails:  paymentDetails,
		TransactionCode: payment.TransactionCode,
	}

	// Sending Email
	defer helpers.SendMail(checkoutResponse)

	// Respond to client
	c.JSON(http.StatusOK, models.ResponseSuccessWithData{
		Status:  "OK",
		Message: "Successfully created transaction",
		Code:    http.StatusOK,
		Data: gin.H{
			"order_id":         order.ID,
			"transaction_code": order.Payment.TransactionCode,
		},
	})
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
	orderIDStr := c.Param("id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	// Perbarui status order terlebih dahulu
	_, err = CheckAndUpdateOrderStatus(orderID, db)
	if err != nil {
		log.Println("Error updating order status:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update order status",
			"details": err.Error(),
		})
		return
	}

	// Ambil data order setelah status diperbarui
	var order models.Order
	if err := db.Preload("OrderDetails").Preload("Payment").Preload("OrderDetails.Menu").First(&order, "id = ?", orderID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Order not found",
			})
		} else {
			log.Println("Error fetching order:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to fetch order",
				"details": err.Error(),
			})
		}
		return
	}

	// Kirim respon sukses dengan data order
	c.JSON(http.StatusOK, order)
}

func GetOrderByTransactionCode(c *gin.Context, db *gorm.DB) {
	transactionCodeStr := c.Param("transactionCode")
	var order models.Order

	// Use Join to include Payment and filter by transaction code
	err := db.Preload("OrderDetails.Menu"). // Preload nested relations
						Preload("OrderDetails").
						Preload("Payment").
						Joins("JOIN payments ON payments.order_id = orders.id"). // Explicit join
						Where("payments.transaction_code = ?", transactionCodeStr).
						First(&order).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return 404 if no order is found
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Order not found",
			})
		} else {
			// Handle other database errors
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to fetch order",
				"details": err.Error(),
			})
		}
		return
	}

	// Return the order as a JSON response
	c.JSON(http.StatusOK, order)
}

func CheckOrderStatusByID(c *gin.Context, db *gorm.DB) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	// Perbarui status order
	orderStatus, err := CheckAndUpdateOrderStatus(orderID, db)
	if err != nil {
		log.Println("Error updating order status:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update order status",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Order status updated successfully",
		"order_status": orderStatus,
	})
}

func CheckAndUpdateOrderStatus(orderID int, db *gorm.DB) (string, error) {
	// Cari data payment berdasarkan order_id
	var payment models.Payment
	if err := db.Where("order_id = ?", orderID).First(&payment).Error; err != nil {
		return "", fmt.Errorf("failed to fetch payment: %w", err)
	}

	// Periksa transaksi menggunakan layanan eksternal
	result, errResp, errCheckTrx := services.CheckTransaction(payment.PaymentTransactionID)
	if errResp != nil || errCheckTrx != nil {
		return "", fmt.Errorf("error checking transaction: %v %v", errResp, errCheckTrx)
	}

	// Tentukan status order berdasarkan status transaksi
	var orderStatus string
	switch result.TransactionStatus {
	case "authorize":
		orderStatus = "authorized"
	case "capture":
		orderStatus = "captured"
	case "settlement":
		orderStatus = "success"
	case "deny":
		orderStatus = "denied"
	case "pending":
		orderStatus = "pending"
	case "cancel":
		orderStatus = "canceled"
	case "refund":
		orderStatus = "refunded"
	case "partial_refund":
		orderStatus = "partially_refunded"
	case "chargeback":
		orderStatus = "charged_back"
	case "partial_chargeback":
		orderStatus = "partially_charged_back"
	case "expire":
		orderStatus = "expired"
	case "failure":
		orderStatus = "failed"
	default:
		orderStatus = "unknown"
	}

	// Perbarui status order di database
	if err := db.Model(&models.Order{}).Where("id = ?", orderID).Update("order_status", orderStatus).Error; err != nil {
		return "", fmt.Errorf("failed to update order status: %w", err)
	}

	// Perbarui status payment di database
	if err := db.Model(&models.Payment{}).Where("order_id = ?", orderID).Update("payment_status", orderStatus).Error; err != nil {
		return "", fmt.Errorf("failed to update payment status: %w", err)
	}

	return orderStatus, nil
}
