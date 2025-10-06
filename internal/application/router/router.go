package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Register registers all HTTP routes on the given engine.
func Register(router *gin.Engine, db *gorm.DB) {
	// API base group
	api := router.Group("/api")

	// Health check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Initialize services
	// TODO: Initialize services when implementing the actual handlers
	_ = db // Use db to avoid unused variable warning

	// Auth endpoints
	// TODO: Uncomment when auth controller is implemented
	// authCtrl := controller.NewAuthController(authService)
	// api.POST("/auth/register", authCtrl.Register)
	// api.POST("/auth/login", authCtrl.Login)
	// api.POST("/auth/refresh", authCtrl.RefreshToken)
	// api.POST("/auth/password-reset/request", authCtrl.RequestPasswordReset)
	// api.POST("/auth/password-reset/confirm", authCtrl.ConfirmPasswordReset)

	// User endpoints
	api.GET("/users/me", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Get user profile (not implemented)"})
	})
	api.PUT("/users/me", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Update user profile (not implemented)"})
	})
	api.GET("/users/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Get user by ID (not implemented)"})
	})
	api.GET("/users", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Get all users (not implemented)"})
	})
	api.POST("/users", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Create user (not implemented)"})
	})
	api.PUT("/users/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Update user (not implemented)"})
	})
	api.DELETE("/users/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Delete user (not implemented)"})
	})

	// Product endpoints
	// TODO: Uncomment when product controller is implemented
	// productCtrl := controller.NewProductController(productService)
	// api.GET("/products", productCtrl.GetProducts)
	// api.GET("/products/:id", productCtrl.GetProduct)
	// api.POST("/products", productCtrl.CreateProduct)
	// api.PUT("/products/:id", productCtrl.UpdateProduct)
	// api.DELETE("/products/:id", productCtrl.DeleteProduct)

	// Category endpoints
	// TODO: Uncomment when category controller is implemented
	// categoryCtrl := controller.NewCategoryController(categoryService)
	// api.GET("/categories", categoryCtrl.GetCategories)
	// api.GET("/categories/:id", categoryCtrl.GetCategory)
	// api.POST("/categories", categoryCtrl.CreateCategory)
	// api.PUT("/categories/:id", categoryCtrl.UpdateCategory)
	// api.DELETE("/categories/:id", categoryCtrl.DeleteCategory)

	// Cart endpoints
	// TODO: Uncomment when cart controller is implemented
	// cartCtrl := controller.NewCartController(cartService)
	// api.GET("/cart", cartCtrl.GetCart)
	// api.POST("/cart/items", cartCtrl.AddToCart)
	// api.PUT("/cart/items/:id", cartCtrl.UpdateCartItem)
	// api.DELETE("/cart/items/:id", cartCtrl.RemoveFromCart

	// Order endpoints
	// TODO: Uncomment when order controller is implemented
	// orderCtrl := controller.NewOrderController(orderService)
	// api.POST("/orders", orderCtrl.CreateOrder)
	// api.GET("/orders", orderCtrl.GetUserOrders)
	// api.GET("/orders/:id", orderCtrl.GetOrderByID)

	// Payment endpoints
	// TODO: Uncomment when payment controller is implemented
	// paymentCtrl := controller.NewPaymentController(paymentService)
	// api.POST("/payments/create-payment-intent", paymentCtrl.CreatePaymentIntent)
	// api.POST("/payments/webhook", paymentCtrl.HandleWebhook)

	// File uploads
	api.Static("/files", "./uploads")
	api.POST("/upload", func(c *gin.Context) {
		// TODO: Implement file upload handler
		c.JSON(200, gin.H{"message": "File upload endpoint"})
	})

	// Admin routes (protected by admin middleware)
	admin := api.Group("/admin")
	admin.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			// TODO: Implement admin middleware
			c.Next()
		}
	}())

	// Admin user management
	// TODO: Uncomment when user controller is implemented
	// admin.GET("/users", userCtrl.GetAllUsers)
	// admin.GET("/users/:id", userCtrl.GetUserByID)
	// admin.DELETE("/users/:id", userCtrl.DeleteUser)

	// Admin product management
	admin.POST("/products/import", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Product import endpoint (not implemented)"})
	})

	admin.POST("/products/export", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Product export endpoint (not implemented)"})
	})

	// Admin order management
	admin.GET("/all-orders", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "All orders endpoint (not implemented)"})
	})

	admin.PUT("/orders/:id/status", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Order status update endpoint (not implemented)"})
	})

	// Admin statistics
	admin.GET("/stats/orders", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Order stats endpoint (not implemented)"})
	})

	admin.GET("/stats/revenue", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Revenue stats endpoint (not implemented)"})
	})
}
