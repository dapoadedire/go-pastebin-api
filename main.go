package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
	"time"
	"fmt"
	"github.com/google/uuid"
	
)

const (
	DB_HOST     = "aws-0-us-west-1.pooler.supabase.com"
	DB_PASSWORD = "pastebin1234:)"
	DB_PORT     = 5432
	DB_USER     = "postgres.pynkrpxqkjoypqdibyru"
	DB_NAME     = "postgres"
)

type Paste struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// var pastes = make(map[string]*Paste)

func main() {
	initDB()
	defer db.Close()
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/pastes", createPaste)
		api.GET("/pastes/:id", getPaste)
		api.DELETE("/pastes/:id", deletePaste)
		api.GET("/pastes", getPastes)
	}

	router.Run(":8080")
}

const UUID_PREFIX = "p_"
func generateUniqueID() string {
	return UUID_PREFIX + uuid.New().String()[0:8]
}

type CreatePasteRequest struct {
	Content string `json:"content" binding:"required"`
}

func createPaste(c *gin.Context) {
	var request CreatePasteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	content := request.Content
	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content is required"})
		return
	}

	id := generateUniqueID()
	createdAt := time.Now()

	_, err := db.Exec("INSERT INTO pastes (id, content, created_at) VALUES ($1, $2, $3)", id, content, createdAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create paste"})
		// print the error
		fmt.Println(err)
		return
	}

	paste := &Paste{
		ID:        id,
		Content:   content,
		CreatedAt: createdAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, paste)
}

func getPaste(c *gin.Context) {
	id := c.Param("id")

	var paste Paste
	err := db.QueryRow("SELECT id, content, created_at FROM pastes WHERE id = $1", id).Scan(&paste.ID, &paste.Content, &paste.CreatedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve paste"})
		return
	}

	c.JSON(http.StatusOK, paste)
}

func getPastes(c *gin.Context) {
	rows, err := db.Query("SELECT id, content, created_at FROM pastes")


	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve pastes"})
		return
	}
	defer rows.Close()

	var pastes []*Paste
	for rows.Next() {
		var paste Paste
		err := rows.Scan(&paste.ID, &paste.Content, &paste.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve pastes"})
			return
		}
		pastes = append(pastes, &paste)
	}

	c.JSON(http.StatusOK, pastes)
}

func deletePaste(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM pastes WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete paste"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("paste %s has been deleted successfully", id)})
}

var db *sql.DB

func initDB() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", DB_HOST, DB_PORT, DB_USER, DB_NAME, DB_PASSWORD)
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS pastes (id TEXT PRIMARY KEY, content TEXT, created_at TIMESTAMP)")
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to the database!")
}
