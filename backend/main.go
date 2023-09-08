// main.go

package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"forum/config"
	"forum/database"
	"forum/handlers"
	"forum/middleware"
	"forum/routes"
)

func main() {
	connStr, err := config.Init()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.Connect(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	handlers.InitializeDB(db)
	router := gin.Default()

	// Configure and set up routes
	routes.SetupRoutes(router, db)

	// Use middleware
	router.Use(middleware.CheckDBConnection(db))

	router.Run("127.0.0.1:8080")
}
