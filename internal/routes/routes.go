package routes

import (
	"gold-management-system/internal/config"
	"gold-management-system/internal/controllers"
	"gold-management-system/internal/middleware"
	"gold-management-system/internal/services"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// Set up CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{
			"Origin",
			"Content-Type",
			"Authorization",
			"Accept",
			"X-Requested-With",
		},
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:5173"
		},
	}))

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
	storeController := controllers.NewStoreController(storeService, adminService)

	// Auth routes
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/logout", authController.Logout)
		authGroup.GET("/ping", authController.Ping)
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
	}

	// Transaction routes (protected)
	transactionGroup := router.Group("/transactions")
	transactionGroup.Use(middleware.AuthMiddleware(cfg))
	{
		transactionGroup.POST("/", transactionController.CreateTransaction)
		transactionGroup.GET("/", transactionController.GetAllTransactions)
		transactionGroup.GET("/:id", transactionController.GetTransactionByID)
		transactionGroup.POST("/filter", transactionController.GetTransactionsByFilters)
		transactionGroup.GET("/report/daily", transactionController.GetDailyReport)
		transactionGroup.GET("/report/monthly", transactionController.GetMonthlyReport)
		transactionGroup.GET("/report/range", transactionController.GetReport)
	}

	// Store routes (protected)
	storeGroup := router.Group("/stores")
	storeGroup.Use(middleware.AuthMiddleware(cfg))
	{
		storeGroup.POST("/", storeController.CreateStore)
		storeGroup.GET("/all-stores", storeController.GetAllStores)
		storeGroup.GET("/", storeController.GetAllStoresByAdmin)
		storeGroup.GET("/:id", storeController.GetStoreByName)
		storeGroup.PUT("/:id", storeController.UpdateStore)
		storeGroup.DELETE("/:id", storeController.DeleteStore)
		storeGroup.PUT("/mangedBy", storeController.UpdateStoreManagedBy)
	}

	// Dashboard routes (protected)
	dashboardGroup := router.Group("/dashboard")
	dashboardGroup.Use(middleware.AuthMiddleware(cfg))
	{
		dashboardGroup.GET("/", transactionController.GetDashboard)
	}

	return router
}