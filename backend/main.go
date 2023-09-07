package main

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

var db *sql.DB
var err error

type user struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type post struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Text_Content string `json:"text_content"`
	User         user   `json:"user"`
}

type customClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func CheckDBConnection() gin.HandlerFunc { //g.HandlerFunc is gin's middleware denotion
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

func validateToken(c *gin.Context) (*int, error) {
	cookie, err := c.Cookie("jwtToken")
	if err != nil {
		if err == http.ErrNoCookie {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Cookie not found"})
			return nil, err
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return nil, err
		}
	}

	cookieValue := cookie
	fmt.Println("Received cookie value:", cookieValue)

	parsedToken, err := jwt.ParseWithClaims(cookieValue, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		hexKey := os.Getenv("KEY")
		key, _ := hex.DecodeString(hexKey)
		return key, nil
	})

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error parsing token: %s", err)})
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*customClaims)
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid token claims"})
		return nil, errors.New("invalid token claims")
	}

	fmt.Println("User ID:", claims.UserID)

	if time.Now().After(claims.ExpiresAt.Time) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
		return nil, errors.New("token has expired")
	}

	return &claims.UserID, nil
}

func deletePosts(c *gin.Context) {
	var curPost post
	if err := c.BindJSON(&curPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	_, err = db.Exec("DELETE FROM posts where posts.id = ($1)", curPost.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't delete post"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Deleted post"})

}

func createPosts(c *gin.Context) {
	userID, err := validateToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}
	var newPost post
	if err := c.BindJSON(&newPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	_, err = db.Exec("INSERT INTO posts (title,text_content,user_id) VALUES ($1,$2,$3)", newPost.Title, newPost.Text_Content, userID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "could not create post"})
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Added Post"})
}

func requestUser(c *gin.Context) {
	userID, err := validateToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	var dbUser user
	err = db.QueryRow("SELECT id, username,password FROM users WHERE id = ($1)", *userID).
		Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Cookie-valid", "user": dbUser})
}

func signOut(c *gin.Context) {
	_, err := validateToken(c)

	if err != nil {
		if err == http.ErrNoCookie {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cookie not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func checkCredentials(c *gin.Context) {
	var curUser user
	if err := c.BindJSON(&curUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	var dbUser user
	err = db.QueryRow("SELECT id, username, password FROM users WHERE username = ($1)", curUser.Username).
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

	claims := customClaims{
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

func postUsers(c *gin.Context) {
	var newUser user
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	_, err := db.Exec("INSERT INTO users (username,password) VALUES ($1,$2)", newUser.Username, newUser.Password)
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

func getUsers(c *gin.Context) {
	rows, err := db.Query("SELECT u.id, u.username FROM users u")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	users := []user{}
	for rows.Next() {
		u := user{}
		err := rows.Scan(&u.ID, &u.Username)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		users = append(users, u)
	}

	c.IndentedJSON(http.StatusOK, users)
}

func getPosts(c *gin.Context) {

	rows, err := db.Query("SELECT p.id, p.title,p.text_content, u.id AS user_id, u.username FROM posts p JOIN users u ON p.user_id = u.id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	posts := []post{}
	for rows.Next() {
		p := post{}
		err := rows.Scan(&p.ID, &p.Title, &p.Text_Content, &p.User.ID, &p.User.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		posts = append(posts, p)
	}

	c.IndentedJSON(http.StatusOK, posts)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbSSLMode := os.Getenv("DB_SSL_MODE")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=%s",
		dbUser, dbPassword, dbName, dbHost, dbSSLMode)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select version()")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var version string
	for rows.Next() {
		err := rows.Scan(&version)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("version=%s\n", version)

	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://127.0.0.1:3000"}
	config.AllowCredentials = true
	router.Use(cors.New(config))
	router.GET("/posts", getPosts)
	router.POST("/posts", createPosts)
	router.DELETE("/posts", deletePosts)
	router.POST("/users", postUsers)
	router.GET("/users", getUsers)
	router.POST("/users/signin", checkCredentials)
	router.GET("/users/signout", signOut)
	router.GET("/auth/getUser", requestUser)
	router.Use(CheckDBConnection())

	router.Run("127.0.0.1:8080")
}
