package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"pastebin/database"
	"pastebin/helper"
	"pastebin/model"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)


func CreatePaste(c *gin.Context) {
	var request model.CreatePasteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	content := request.Content
	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content is required"})
		return
	}

	id := helper.GenerateUniqueID()
	createdAt := time.Now()

	_, err := database.DB.Exec("INSERT INTO pastes (id, content, created_at) VALUES ($1, $2, $3)", id, content, createdAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create paste"})
		// print the error
		fmt.Println(err)
		return
	}

	paste := &model.Paste{
		ID:        id,
		Content:   content,
		CreatedAt: createdAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, paste)
}

func GetPaste(c *gin.Context) {
	id := c.Param("id")

	var paste model.Paste
	err := database.DB.QueryRow("SELECT id, content, created_at FROM pastes WHERE id = $1", id).Scan(&paste.ID, &paste.Content, &paste.CreatedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve paste"})
		return
	}

	c.JSON(http.StatusOK, paste)
}

func GetAllPastes(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, content, created_at FROM pastes")


	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve pastes"})
		return
	}
	defer rows.Close()

	var pastes []*model.Paste
	for rows.Next() {
		var paste model.Paste
		err := rows.Scan(&paste.ID, &paste.Content, &paste.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve pastes"})
			return
		}
		pastes = append(pastes, &paste)
	}

	c.JSON(http.StatusOK, pastes)
}

func DeletePaste(c *gin.Context) {
	id := c.Param("id")

	// Check if the paste exists
	var paste model.Paste
	err := database.DB.QueryRow("SELECT id FROM pastes WHERE id = $1", id).Scan(&paste.ID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve paste"})
		return
	}

	// Delete the paste
	_, err = database.DB.Exec("DELETE FROM pastes WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete paste"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("paste %s has been deleted successfully", id)})
}

