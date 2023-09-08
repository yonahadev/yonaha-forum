package handlers

import (
	"database/sql"
	"fmt"

	"net/http"

	"forum/structs"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func PostUsers(c *gin.Context) {
	var newUser structs.User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	_, err := DB.Exec("INSERT INTO users (username,password) VALUES ($1,$2)", newUser.Username, newUser.Password)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user data."})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Added user"})
}

func GetUsers(c *gin.Context) {
	rows, err := DB.Query("SELECT u.id, u.username FROM users u")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	users := []structs.User{}
	for rows.Next() {
		u := structs.User{}
		err := rows.Scan(&u.ID, &u.Username)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		users = append(users, u)
	}

	c.IndentedJSON(http.StatusOK, users)
}

func RequestUser(c *gin.Context) {
	userID := c.GetInt("UserID")
	fmt.Println("User ID (request user function) is", userID)

	var dbUser structs.User
	err := DB.QueryRow("SELECT id, username,password FROM users WHERE id = ($1)", userID).
		Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User could not be found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Cookie-valid", "user": dbUser})
}
