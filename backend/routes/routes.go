// routes/routes.go

package routes

import (
	"database/sql"

	"forum/handlers"
	"forum/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRoutes initializes and configures the routes for your application.
func SetupRoutes(router *gin.Engine, db *sql.DB) {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://127.0.0.1:3000"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	// Define your API routes and associate them with handlers
	router.GET("/posts", handlers.GetPosts)
	router.POST("/posts", middleware.ValidateToken, handlers.CreatePosts)
	router.DELETE("/posts", middleware.ValidateToken, handlers.DeletePosts)
	router.PATCH("/posts", middleware.ValidateToken, handlers.UpdatePosts)
	router.POST("/users", handlers.PostUsers)
	router.GET("/users", handlers.GetUsers)
	router.POST("/users/signin", handlers.CreateToken)
	router.GET("/users/signout", middleware.ValidateToken, handlers.SignOut)
	router.GET("/auth/getUser", middleware.ValidateToken, handlers.RequestUser)
}
