package routes

import (
	"github.com/hutamy/invoice-generator/controllers"
	_ "github.com/hutamy/invoice-generator/docs"
	"github.com/hutamy/invoice-generator/middleware"
	"github.com/hutamy/invoice-generator/repositories"
	"github.com/hutamy/invoice-generator/services"
	"github.com/hutamy/invoice-generator/utils"
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

	clientRoutes := protected.Group("/clients")
	clientRoutes.POST("", clientController.CreateClient)
	clientRoutes.GET("", clientController.GetAllClients)
	clientRoutes.GET("/:id", clientController.GetClientByID)
	clientRoutes.PUT("/:id", clientController.UpdateClient)
	clientRoutes.DELETE("/:id", clientController.DeleteClient)

	protectedInvoiceRoutes := protected.Group("/invoices")
	protectedInvoiceRoutes.POST("", invoiceController.CreateInvoice)
	protectedInvoiceRoutes.GET("/:id", invoiceController.GetInvoiceByID)
	protectedInvoiceRoutes.PATCH("/:id", invoiceController.UpdateInvoice)
	protectedInvoiceRoutes.GET("", invoiceController.ListInvoicesByUserID)
	protectedInvoiceRoutes.POST("/:id/pdf", invoiceController.DownloadInvoicePDF)
}
