package controllers

import (
	"caaso/models"
	"caaso/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBenefits(c *gin.Context) {

	var benefits []models.Benefit
	result := services.DB.Find(&benefits)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message":  "Benefícios encontrados",
		"benefits": benefits,
	})

}

func CreateBenefit(c *gin.Context) {

	var benefitInput models.BenefitInput

	if err := c.ShouldBindJSON(&benefitInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	benefit := models.Benefit{
		Title:       benefitInput.Title,
		Description: benefitInput.Description,
		PhotoUrl:    benefitInput.PhotoUrl,
	}

	result := services.DB.Create(&benefit)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Benefício criado com sucesso",
		"benefit": benefit,
	})

}

func UpdateBenefit(c *gin.Context) {

	benefitId := c.Param("id")

	var benefit models.Benefit
	result := services.DB.Where("id=?", benefitId).Find(&benefit)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"messsage": result.Error.Error()})
		return
	}

	var benefitInput models.BenefitInput
	if err := c.ShouldBindJSON(&benefitInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	benefit.Title = benefitInput.Title
	benefit.Description = benefitInput.Description
	benefit.PhotoUrl = benefitInput.PhotoUrl

	result = services.DB.Save(&benefit)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Benefício atualizado com sucesso",
		"benefit": benefit,
	})

}

func DeleteBenefit(c *gin.Context) {

	benefitId := c.Param("id")

	var benefit models.Benefit
	result := services.DB.Where("id=?", benefitId).Find(&benefit)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
		return
	}

	result = services.DB.Delete(&benefit)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Benefício deletado com sucesso",
	})

}
