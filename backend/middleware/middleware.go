package middleware

import (
	"database/sql"
	"encoding/hex"
	"forum/structs"
	"os"
	"time"

	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CheckDBConnection(db *sql.DB) gin.HandlerFunc { //g.HandlerFunc is gin's middleware denotion
	return func(c *gin.Context) {
		if db == nil {
			fmt.Println("failed to establish a db connection")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to establish a database connection.",
			})
			return
		}
		c.Next() //continues to next middleware/actual function
	}
}

func ValidateToken(c *gin.Context) {
	cookie, err := c.Cookie("jwtToken")
	if err != nil {
		if err == http.ErrNoCookie {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Cookie not found"})
			fmt.Println(err)
			return
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			fmt.Println(err)
			return
		}
	}

	cookieValue := cookie
	fmt.Println("Received cookie value:", cookieValue)

	parsedToken, err := jwt.ParseWithClaims(cookieValue, &structs.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		hexKey := os.Getenv("KEY")
		key, _ := hex.DecodeString(hexKey)
		return key, nil
	})

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error parsing token: %s", err)})
		return
	}

	claims, ok := parsedToken.Claims.(*structs.CustomClaims)
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid token claims"})
		return
	}

	fmt.Println("Received User ID:", claims.UserID)

	if time.Now().After(claims.ExpiresAt.Time) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
		return
	}

	c.Set("UserID", claims.UserID)
	c.Next()
}
