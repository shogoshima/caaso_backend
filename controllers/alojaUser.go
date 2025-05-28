package controllers

import (
	"caaso/models"
	"caaso/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAlojaUsers(c *gin.Context) {

	var alojaUsers []models.AlojaUser
	result := services.DB.Find(&alojaUsers)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Moradores do aloja registrados",
		"users":   alojaUsers,
	})

}

func CreateMultipleAlojaUsers(c *gin.Context) {

	var alojaUsers []models.AlojaUser

	if err := c.ShouldBindJSON(&alojaUsers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	for _, alojaUser := range alojaUsers {

		// Check if the user already exists
		var existingUser models.AlojaUser
		result := services.DB.Where("email = ?", alojaUser.Email).First(&existingUser)
		if result.Error == nil {
			// User already exists, skip creation
			continue
		}

		// Create the new user
		err := services.DB.Create(&alojaUser).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
			return
		}
	}

	c.JSON(200, gin.H{
		"message": "Moradores adicionados com sucesso",
	})

}

func DeleteAlojaUser(c *gin.Context) {

	alojaUserId := c.Param("id")

	var alojaUser models.AlojaUser
	result := services.DB.Where("id=?", alojaUserId).Find(&alojaUser)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"messsage": result.Error.Error()})
		return
	}

	result = services.DB.Delete(&alojaUser)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Morador removido com sucesso",
	})

}
