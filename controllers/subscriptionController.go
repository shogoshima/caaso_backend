package controllers

import (
	"caaso/models"
	"caaso/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetSubscription(c *gin.Context) {
	token := c.Param("token")

	var user models.User
	services.DB.Where("token=?", token).Find(&user)

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
		"photoUrl":     user.PhotoUrl,
		"type":         user.Type,
	})
}

func UpdateSubscription(c *gin.Context) {
	uidRaw, ok := c.Get("currentUserId")
	if !ok {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	userId := uidRaw.(string)

	ptRaw, ok := c.Get("planType")
	if !ok {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	planTypeStr := ptRaw.(string)

	// mark them as subscribed
	if err := services.DB.
		Model(&models.User{}).
		Where("id = ?", userId).
		Update("is_subscribed", true).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// set the expiration
	var expire time.Time
	if planTypeStr == models.Monthly.String() {
		expire = time.Now().AddDate(0, 1, 0)
	} else {
		expire = time.Now().AddDate(1, 0, 0)
	}
	if err := services.DB.
		Model(&models.User{}).
		Where("id = ?", userId).
		Update("expiration_date", expire).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
