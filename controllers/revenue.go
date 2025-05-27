package controllers

import (
	"caaso/models"
	"caaso/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRevenues(c *gin.Context) {
	var revenues []models.Revenue
	result := services.DB.Find(&revenues)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message":  "Receitas encontradas",
		"revenues": revenues,
	})
}