package main

import (
	"carbon-api/config"
	"carbon-api/routes"
	"context"
	"os"

	"github.com/labstack/echo/v4"
)

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
