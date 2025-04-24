package middlewares

import (
	"caaso/controllers"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GoogleAuth(c *gin.Context) {

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

	tokenString := authToken[1]
	token, err := controllers.VerifyIDToken(c, tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("currentUserID", token.UID)

	c.Next()

}
