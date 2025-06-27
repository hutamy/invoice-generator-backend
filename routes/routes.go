package routes

import (
	"github.com/hutamy/invoice-generator-backend/controllers"
	_ "github.com/hutamy/invoice-generator-backend/docs"
	"github.com/hutamy/invoice-generator-backend/middleware"
	"github.com/hutamy/invoice-generator-backend/repositories"
	"github.com/hutamy/invoice-generator-backend/services"
	"github.com/hutamy/invoice-generator-backend/utils"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

func InitRoutes(e *echo.Echo, db *gorm.DB) {
	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo)
	authController := controllers.NewAuthController(authService)

	clientRepo := repositories.NewClientRepository(db)
	clientService := services.NewClientService(clientRepo)
	clientController := controllers.NewClientController(clientService)

	invoiceRepo := repositories.NewInvoiceRepository(db)
	invoiceService := services.NewInvoiceService(invoiceRepo, clientRepo, authRepo)
	invoiceController := controllers.NewInvoiceController(invoiceService)

	// Routes for Health Check and Welcome Message
	e.GET("/", func(c echo.Context) error {
		return utils.Response(c, 200, "Welcome to Invoice Generator API", nil)
	})
	e.GET("/health", func(c echo.Context) error {
		return utils.Response(c, 200, "Invoice Generator API is running", nil)
	})
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// These routes are grouped under the "/v1" path
	v1 := e.Group("/v1")
	public := v1.Group("/public")

	// Public Routes
	authRoutes := public.Group("/auth")
	authRoutes.POST("/sign-up", authController.SignUp)
	authRoutes.POST("/sign-in", authController.SignIn)

	publicInvoiceRoutes := public.Group("/invoices")
	publicInvoiceRoutes.POST("/generate-pdf", invoiceController.GeneratePublicInvoice)

	protected := v1.Group("/protected")
	protected.Use(middleware.JWTMiddleware)

	protected.GET("/me", authController.Me)
	protected.PUT("/me", authController.UpdateUser)

	authPrivateRoutes := protected.Group("/auth")
	authPrivateRoutes.POST("/refresh-token", authController.RefreshToken)

	clientRoutes := protected.Group("/clients")
	clientRoutes.POST("", clientController.CreateClient)
	clientRoutes.GET("", clientController.GetAllClients)
	clientRoutes.GET("/:id", clientController.GetClientByID)
	clientRoutes.PUT("/:id", clientController.UpdateClient)
	clientRoutes.DELETE("/:id", clientController.DeleteClient)

	protectedInvoiceRoutes := protected.Group("/invoices")
	protectedInvoiceRoutes.GET("/summary", invoiceController.InvoiceSummary)
	protectedInvoiceRoutes.POST("", invoiceController.CreateInvoice)
	protectedInvoiceRoutes.GET("/:id", invoiceController.GetInvoiceByID)
	protectedInvoiceRoutes.PUT("/:id", invoiceController.UpdateInvoice)
	protectedInvoiceRoutes.DELETE("/:id", invoiceController.DeleteInvoice)
	protectedInvoiceRoutes.GET("", invoiceController.ListInvoicesByUserID)
	protectedInvoiceRoutes.PATCH("/:id/status", invoiceController.UpdateInvoiceStatus)
	protectedInvoiceRoutes.POST("/:id/pdf", invoiceController.DownloadInvoicePDF)
}
