package router

import (
	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/api/handler"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/api/middleware"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/repository"
	"github.com/hoanggtan02/my-app/vat-simple-backend/internal/service"
)

// SetupRouter initializes and configures the Gin router.
func SetupRouter(db *sql.DB) *gin.Engine {
	// Khởi tạo các tầng (dependency injection)
	userRepo := repository.NewUserRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	productRepo := repository.NewProductRepository(db)
	invoiceRepo := repository.NewInvoiceRepository(db)

	authService := service.NewAuthService(userRepo)
	customerService := service.NewCustomerService(customerRepo)
	productService := service.NewProductService(productRepo)
	invoiceService := service.NewInvoiceService(invoiceRepo, customerRepo, productRepo)

	authHandler := handler.NewAuthHandler(authService)
	customerHandler := handler.NewCustomerHandler(customerService)
	productHandler := handler.NewProductHandler(productService)
	invoiceHandler := handler.NewInvoiceHandler(invoiceService, productService)

	// Cài đặt Gin router
	r := gin.Default()

	// *** CẬP NHẬT CẤU HÌNH CORS CHI TIẾT ***
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Gom nhóm các API routes
	api := r.Group("/api/v1")
	{
		// Public Routes
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
		}

		// Protected Routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			userRoutes := protected.Group("/users")
			{
				userRoutes.GET("/me", authHandler.GetUserProfile)
				userRoutes.PUT("/me", authHandler.UpdateUserProfile)
			}

			customerRoutes := protected.Group("/customers")
			{
				customerRoutes.POST("", customerHandler.CreateCustomer)
				customerRoutes.GET("", customerHandler.ListCustomers)
				customerRoutes.GET("/:id", customerHandler.GetCustomer)
				customerRoutes.PUT("/:id", customerHandler.UpdateCustomer)
				customerRoutes.DELETE("/:id", customerHandler.DeleteCustomer)
			}

			productRoutes := protected.Group("/products")
			{
				productRoutes.POST("", productHandler.CreateProduct)
				productRoutes.GET("", productHandler.ListProducts)
			}

			invoiceRoutes := protected.Group("/invoices")
			{
				invoiceRoutes.POST("", invoiceHandler.CreateInvoice)
				invoiceRoutes.GET("", invoiceHandler.ListInvoices)
				invoiceRoutes.GET("/:id", invoiceHandler.GetInvoice)
			}
		}
	}

	return r
}
