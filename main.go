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

	e := echo.New()
	routes.Init(e)

	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("API_PORT")
	}
	e.Logger.Fatal(e.Start(":" + port))
}
