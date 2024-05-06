package main

import (
	"github.com/gin-gonic/gin"
	
	
	"pastebin/database"
	"pastebin/controller"

	
	
)


func main() {
	database.InitDB()
	defer database.DB.Close()
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/pastes", controller.CreatePaste)
		api.GET("/pastes/:id", controller.GetPaste)
		api.DELETE("/pastes/:id", controller.DeletePaste)
		api.GET("/pastes", controller.GetAllPastes)
	}

	router.Run(":8080")
}

