package main

import (
	"invitified-go/config"
	"invitified-go/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title University GC API
// @version 1.0
// @description This is a sample server for a university API.
// @contact.name kevinsofyan
// @contact.email kevinsofyan.13@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	e := echo.New()

	// Initialize the database
	config.InitDB()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize routes
	routes.InitRoutes(e)

	// Swagger
	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Get the port from the environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	log.Printf("Starting server on port %s", port)

	e.Logger.Fatal(e.Start(":" + port))
}
