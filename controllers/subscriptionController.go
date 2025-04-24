package controllers

import (
	"caaso/models"
	"caaso/services"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetSubscription(c *gin.Context) {
	nusp := c.Param("nusp")

	var user models.User
	services.DB.Where("nusp=?", nusp).Find(&user)

	if user.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": "Usuário não encontrado"})
		return
	}

	if !user.ExpirationDate.IsZero() {
		if float64(time.Now().Unix()) > float64(user.ExpirationDate.Unix()) {
			user.IsSubscribed = false
			services.DB.Save(&user)
		}
	}

	c.JSON(200, gin.H{
		"message":      "Usuário encontrado",
		"displayName":  user.DisplayName,
		"isSubscribed": user.IsSubscribed,
	})
}

func UpdateSubscription(c *gin.Context) {
	uidRaw, ok := c.Get("currentUserId")
	if !ok {
		return
	}
	userId := uidRaw.(string)

	ptRaw, ok := c.Get("planType")
	if !ok {
		return
	}
	planTypeStr := ptRaw.(string)

	fmt.Println(userId, planTypeStr)

	// mark them as subscribed
	if err := services.DB.
		Model(&models.User{}).
		Where("id = ?", userId).
		Update("is_subscribed", true).Error; err != nil {
		return
	}

	// set the expiration
	var expire time.Time
	if planTypeStr == "Monthly" {
		expire = time.Now().AddDate(0, 1, 0)
	} else {
		expire = time.Now().AddDate(1, 0, 0)
	}
	if err := services.DB.
		Model(&models.User{}).
		Where("id = ?", userId).
		Update("expiration_date", expire).Error; err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assinatura atualizada com sucesso"})
}
