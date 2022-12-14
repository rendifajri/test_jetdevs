package main

import (
	"log"
	"os"
	"test_jetdevs/database"
	"test_jetdevs/routes"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	database.ShareConnection = database.Connection()
	port := os.Getenv("APP_PORT")
	e := echo.New()
	routes.Build(e)
	e.Logger.Fatal(e.Start(":" + port))
}
