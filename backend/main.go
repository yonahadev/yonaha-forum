package main

import (
	"database/sql"
	"encoding/hex"
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
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Password *string `json:"password"`
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

func checkCookie(c *gin.Context) {
	cookie, err := c.Cookie("jwtToken")
	if err != nil {
		if err == http.ErrNoCookie {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Cookie not found"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return // Return early when the cookie is not found or there's an error
	}

	// You can access the cookie value with cookie.Value
	cookieValue := cookie
	fmt.Println("Received cookie value:", cookieValue)

	rawCookie, err := jwt.ParseWithClaims(cookieValue, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		hexKey := os.Getenv("KEY")
		key, _ := hex.DecodeString(hexKey)
		return key, nil
	})

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error passing token %s", err)})
		return
	}

	claims, ok := rawCookie.Claims.(*customClaims)
	if !ok {
		// Handle the case where the claims couldn't be cast to your custom claims struct
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid token claims"})
		return
	}

	// Log the claims or perform any necessary actions with them
	fmt.Println("User ID:", claims.UserID)
	fmt.Println("Expiration Time:", claims.ExpiresAt)

	var dbUser user
	err = db.QueryRow("SELECT id, username FROM users WHERE id = ($1)", claims.UserID).
		Scan(&dbUser.ID, &dbUser.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Cookie-valid", "token": cookie, "user": dbUser})
}

func signOut(c *gin.Context) {
	_, err := c.Cookie("jwtToken")
	if err != nil {
		if err == http.ErrNoCookie {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cookie not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}
	expiredCookie := &http.Cookie{
		Name:     "jwtToken",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // This will instruct the browser to delete the cookie.
		HttpOnly: true,
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully logged out", "cookie": expiredCookie})
}

func checkCredentials(c *gin.Context) {
	var curUser user
	if err := c.BindJSON(&curUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	var dbUser user
	err := db.QueryRow("SELECT id, username, password FROM users WHERE username = ($1)", curUser.Username).
		Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	if dbUser.Password != curUser.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
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
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24))
	claims.Issuer, _ = claims.GetIssuer()

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

	// cookie := http.Cookie{
	// 	Name:     "jwtToken",
	// 	Value:    signedToken,
	// 	Expires:  time.Now().Add(time.Hour * 24), // Same expiration as in claims
	// 	HttpOnly: true,
	// }
	// http.SetCookie(c.Writer, &cookie)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Signed in", "token": signedToken})
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

	// Construct your connection string
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
	config.AllowOrigins = []string{"http://127.0.0.1:3000", "http://localhost:3000"}
	config.AllowCredentials = true
	router.Use(cors.New(config))
	router.GET("/posts", getPosts)
	router.POST("/users", postUsers)
	router.GET("/users", getUsers)
	router.POST("/users/signin", checkCredentials)
	router.GET("/users/signout", signOut)
	router.GET("/cookie-test", checkCookie)
	router.Use(CheckDBConnection())

	router.Run("127.0.0.1:8080")
}
