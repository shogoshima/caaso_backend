package middlewares

import (
	"caaso/models"
	"caaso/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckAuth(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is missing"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	authToken := strings.Split(authHeader, " ")
	if len(authToken) != 2 || authToken[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token format"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	idToken := authToken[1]

	// Verify the ID token with Firebase
	ctx := c.Request.Context()
	token, err := services.AuthClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Upsert the user in one round-trip
	var user models.User
	err = services.DB.Where("id = ?", token.UID).Find(&user).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	c.Set("currentUser", user)

	c.Next()
}
