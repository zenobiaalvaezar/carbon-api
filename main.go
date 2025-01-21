package main

import (
	"carbon-api/config"
	"carbon-api/routes"
	"context"
	"os"

	"github.com/labstack/echo/v4"
)

// @title Carbon API
// @version 1.0
// @description This is the API for managing carbon ecosystem.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @servers
// @url https://carbon-api-70017640279.us-central1.run.app
// @description Production server

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config.LoadEnv()

	config.ConnectPostgres()
	defer config.ClosePostgres()

	config.ConnectMongo(context.Background())
	defer config.CloseMongo(context.Background())

	config.ConnectRedis()
	defer config.CloseRedis()

	e := echo.New()
	routes.Init(e)

	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("API_PORT")
	}
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
