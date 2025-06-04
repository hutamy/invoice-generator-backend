package main

import (
	"fmt"
	"log"

	"github.com/hutamy/invoice-generator/config"
	"github.com/hutamy/invoice-generator/routes"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.LoadEnv()
	dbUrl := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDB,
		cfg.PostgresPort,
	)
	db := config.InitDB(dbUrl)

	e := echo.New()
	routes.InitRoutes(e, db)

	log.Printf("Starting server on port: %d", cfg.Port)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.Port)))
}
