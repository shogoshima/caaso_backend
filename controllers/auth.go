package controllers

import (
	"caaso/models"
	"caaso/services"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		c.Abort()
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		c.Abort()
		return
	}

	idToken := parts[1]

	// Verify the ID token with Firebase
	ctx := c.Request.Context()
	token, err := services.AuthClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Fetch the full Firebase user profile
	userRecord, err := services.AuthClient.GetUser(ctx, token.UID)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to load user record"})
		return
	}

	fmt.Println("User record:", userRecord)

	// Upsert the user in one round-trip
	var user models.User
	result := services.DB.
		WithContext(ctx).
		Where(models.User{ID: userRecord.UID}).
		// only set these fields on INSERT
		Attrs(models.User{
			DisplayName: userRecord.DisplayName,
			Email:       userRecord.Email,
			PhotoUrl:    userRecord.PhotoURL,
		}).
		FirstOrCreate(&user)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// If the record already existed, check for profile changes & save
	if result.RowsAffected == 0 {
		changed := false
		if user.DisplayName != userRecord.DisplayName {
			user.DisplayName = userRecord.DisplayName
			changed = true
		}
		if user.PhotoUrl != userRecord.PhotoURL {
			user.PhotoUrl = userRecord.PhotoURL
			changed = true
		}
		if changed {
			if err := services.DB.WithContext(ctx).Save(&user).Error; err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
