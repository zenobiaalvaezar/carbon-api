package main

import (
	"carbon-api/config"
	"carbon-api/routes"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	config.LoadEnv()

	config.ConnectPostgres()
	defer config.ClosePostgres()

	config.ConnectRedis()
	defer config.CloseRedis()

	e := echo.New()
	routes.Init(e)

	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("8080")
	}
	e.Logger.Fatal(e.Start(":" + port))
}
