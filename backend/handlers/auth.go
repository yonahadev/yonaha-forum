package handlers

import (
	"database/sql"
	"encoding/hex"

	"forum/structs"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func SignOut(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func CreateToken(c *gin.Context) {
	var curUser structs.User
	if err := c.BindJSON(&curUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	var dbUser structs.User
	err := DB.QueryRow("SELECT id, username, password FROM users WHERE username = ($1)", curUser.Username).
		Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	if dbUser.Password != curUser.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect credentials"})
		return
	}

	var (
		key         []byte
		token       *jwt.Token
		signedToken string
	)

	claims := structs.CustomClaims{
		UserID: dbUser.ID,
	}
	claims.Issuer, _ = claims.GetIssuer()
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24))

	hexKey := os.Getenv("KEY")
	key, err = hex.DecodeString(hexKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server key error"})
	}
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	signedToken, err = token.SignedString(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server key error"})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Signed in", "token": signedToken, "tokenExpiry": claims.ExpiresAt.Time})
}
