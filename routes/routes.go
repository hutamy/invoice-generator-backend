package routes

import (
	"github.com/hutamy/invoice-generator/controllers"
	"github.com/hutamy/invoice-generator/repositories"
	"github.com/hutamy/invoice-generator/services"
	"github.com/hutamy/invoice-generator/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitRoutes(e *echo.Echo, db *gorm.DB) {
	// Routes for Health Check and Welcome Message
	e.GET("/", func(c echo.Context) error {
		return utils.Response(c, 200, "Welcome to Invoice Generator API", nil)
	})
	e.GET("/health", func(c echo.Context) error {
		return utils.Response(c, 200, "Invoice Generator API is running", nil)
	})

	// These routes are grouped under the "/v1" path
	e.Group("/v1")

	// Authentication Routes
	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo)
	authController := controllers.NewAuthController(authService)
	e.POST("/v1/auth/sign-up", authController.SignUp)
	e.POST("/v1/auth/sign-in", authController.SignIn)
}
