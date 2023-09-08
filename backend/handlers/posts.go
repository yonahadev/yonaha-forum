package handlers

import (
	"fmt"
	"forum/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdatePosts(c *gin.Context) {
	var updatedPost structs.Post
	if err := c.BindJSON(&updatedPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	_, err := DB.Exec("UPDATE posts SET title = ($1), text_content = ($2) WHERE id = ($3)", updatedPost.Title, updatedPost.Text_Content, updatedPost.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't update post"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Updated Post"})

}

func DeletePosts(c *gin.Context) {
	var curPost structs.Post
	if err := c.BindJSON(&curPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	_, err := DB.Exec("DELETE FROM posts where posts.id = ($1)", curPost.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't delete post"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Deleted post"})
}

func CreatePosts(c *gin.Context) {
	userID := c.GetInt("UserID")
	fmt.Println("User ID (create post function) is", userID)

	var newPost structs.Post
	if err := c.BindJSON(&newPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	_, err := DB.Exec("INSERT INTO posts (title,text_content,user_id) VALUES ($1,$2,$3)", newPost.Title, newPost.Text_Content, userID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "could not create post"})
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Added Post"})
}

func GetPosts(c *gin.Context) {

	rows, err := DB.Query("SELECT p.id, p.title,p.text_content, u.id AS user_id, u.username FROM posts p JOIN users u ON p.user_id = u.id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	posts := []structs.Post{}
	for rows.Next() {
		p := structs.Post{}
		err := rows.Scan(&p.ID, &p.Title, &p.Text_Content, &p.User.ID, &p.User.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		posts = append(posts, p)
	}

	c.IndentedJSON(http.StatusOK, posts)
}
