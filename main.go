package main

import (
	"log"
	"net/http"
	"strings"

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

	// Setup CORS
	utils.Cors(r)

	// Group untuk endpoint publik (tidak memerlukan autentikasi)
	public := r.Group("/")
	{
		public.POST("/login", func(c *gin.Context) {
			controllers.HandleLogin(c, db)
		})

		public.POST("/checkout", func(c *gin.Context) {
			controllers.HandleCheckout(c, db)
		})

		public.GET("/menus", func(c *gin.Context) {
			controllers.GetMenu(c, db)
		})

		public.GET("/menus/:id", func(c *gin.Context) {
			controllers.GetMenuByID(c, db)
		})

		public.GET("/categories", func(c *gin.Context) {
			controllers.GetCategories(c, db)
		})

		public.GET("/categories/:id", func(c *gin.Context) {
			controllers.GetCategoriesByID(c, db)
		})

		public.GET("/payment_methods", func(c *gin.Context) {
			controllers.GetPaymentMethods(c, db)
		})

		public.GET("/orders/:id/status", func(c *gin.Context) {
			controllers.CheckOrderStatusByID(c, db)
		})

		public.GET("/orders/:id", func(c *gin.Context) {
			controllers.GetOrderByID(c, db)
		})

	}

	// Middleware untuk validasi JWT
	authMiddleware := func(ctx *gin.Context) {
		// Ambil header Authorization
		authorization := ctx.GetHeader("Authorization")
		if authorization == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			ctx.Abort()
			return
		}

		// Periksa format Bearer
		if !strings.HasPrefix(authorization, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			ctx.Abort()
			return
		}

		// Ekstrak token
		tokenString := strings.TrimPrefix(authorization, "Bearer ")

		// Validasi token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		// Simpan klaim ke context untuk digunakan pada handler berikutnya
		ctx.Set("claims", claims)
		ctx.Next()
	}

	// Group untuk endpoint yang memerlukan autentikasi
	auth := r.Group("/")
	auth.Use(authMiddleware)
	{
		auth.GET("/orders", func(c *gin.Context) {
			controllers.GetAllOrders(c, db)
		})

		auth.POST("/menus", func(c *gin.Context) {
			controllers.CreateNewProduct(c, db)
		})

		auth.PUT("/menus/:id", func(c *gin.Context) {
			controllers.UpdateProduct(c, db)
		})

		auth.POST("/categories", func(c *gin.Context) {
			controllers.CreateCategory(c, db)
		})

		auth.PUT("/categories/:id", func(c *gin.Context) {
			controllers.UpdateCategory(c, db)
		})
	}

	// Start the server
	if err := r.Run(":8081"); err != nil {
		log.Fatal("Server failed:", err)
	}
}
