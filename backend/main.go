package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

type user struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	User  user   `json:"user"`
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

func postUsers(c *gin.Context) {
}

func getPosts(c *gin.Context) {

	rows, err := db.Query("SELECT p.id, p.title, u.id AS user_id, u.username FROM posts p JOIN users u ON p.user_id = u.id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	posts := []post{}
	for rows.Next() {
		p := post{}
		err := rows.Scan(&p.ID, &p.Title, &p.User.ID, &p.User.Username)
		if err != nil {
			log.Fatal(err)
		}
		posts = append(posts, p)
	}

	c.IndentedJSON(http.StatusOK, posts)
}

func main() {
	connStr := "user=yonahadev password=u6hzMbvCnJ7B dbname=neondb host=ep-royal-waterfall-74881687.eu-central-1.aws.neon.tech sslmode=verify-full"
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
	router.GET("/posts", CheckDBConnection(), getPosts)
	router.POST("/users", CheckDBConnection(), postUsers)
	router.Run("127.0.0.1:8080")
}
