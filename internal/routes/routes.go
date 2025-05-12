package routes

import (
	"gold-management-system/internal/config"
	"gold-management-system/internal/controllers"
	"gold-management-system/internal/middleware"
	"gold-management-system/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// Initialize services
	authService := services.NewAuthService(db, cfg)
	adminService := services.NewAdminService(db)
	customerService := services.NewCustomerService(db)
	transactionService := services.NewTransactionService(db)
	storeService := services.NewStoreService(db)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	adminController := controllers.NewAdminController(adminService)
	customerController := controllers.NewCustomerController(customerService)
	transactionController := controllers.NewTransactionController(transactionService)
	storeController := controllers.NewStoreController(storeService)

	// Auth routes
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/register", authController.Register)
	}

	// Admin routes (protected)
	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware(cfg))
	{
		adminGroup.GET("/profile", adminController.GetProfile)
		adminGroup.PUT("/profile", adminController.UpdateProfile)
	}

	// Customer routes (protected)
	customerGroup := router.Group("/customers")
	customerGroup.Use(middleware.AuthMiddleware(cfg))
	{
		customerGroup.POST("/", customerController.CreateCustomer)
		customerGroup.GET("/", customerController.GetAllCustomers)
		customerGroup.GET("/:id", customerController.GetCustomerByID)
		customerGroup.PUT("/:id", customerController.UpdateCustomer)
		customerGroup.DELETE("/:id", customerController.DeleteCustomer)
		customerGroup.GET("/:id/transactions", customerController.GetCustomerTransactions)
		customerGroup.GET("/:id/balance", customerController.GetCustomerGoldBalance)
	}

	// Transaction routes (protected)
	transactionGroup := router.Group("/transactions")
	transactionGroup.Use(middleware.AuthMiddleware(cfg))
	{
		transactionGroup.POST("/", transactionController.CreateTransaction)
		transactionGroup.GET("/", transactionController.GetAllTransactions)
		transactionGroup.GET("/:id", transactionController.GetTransactionByID)
		transactionGroup.GET("/report/daily", transactionController.GetDailyReport)
		transactionGroup.GET("/report/monthly", transactionController.GetMonthlyReport)
		transactionGroup.GET("/report/range", transactionController.GetReport)
	}

	// Store routes (protected)
	storeGroup := router.Group("/stores")
	storeGroup.Use(middleware.AuthMiddleware(cfg))
	{
		storeGroup.POST("/", storeController.CreateStore)
		storeGroup.GET("/", storeController.GetAllStores)
		storeGroup.GET("/:name", storeController.GetStoreByName)
		storeGroup.PUT("/:name", storeController.UpdateStore)
		storeGroup.DELETE("/:name", storeController.DeleteStore)
	}

	return router
}