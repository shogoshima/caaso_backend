package controllers

import (
	"caaso/models"
	"caaso/services"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func GetUserProfile(c *gin.Context) {

	user, _ := c.Get("currentUser")

	User := user.(models.User)
	if User.ExpirationDate.Before(time.Now()) && User.IsSubscribed {
		User.IsSubscribed = false
		services.DB.Save(&User)
	}

	c.JSON(200, gin.H{
		"message": "Usuário autenticado",
		"user":    user,
	})

}

func GetAllUsers(c *gin.Context) {

	var users []models.User
	result := services.DB.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Usuários encontrados",
		"users":   users,
	})

}

func GenerateJwtToken(c *gin.Context) {

	userValue, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Usuário não autenticado"})
		return
	}
	user := userValue.(models.User)

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 14).Unix(),
	})

	sessionToken, err := generateToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao gerar token, tente novamente"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token gerado com sucesso",
		"token":   sessionToken,
		"user":    user,
	})

}

func UpdateUserProfile(c *gin.Context) {
	id := c.Param("id")

	var userInput models.UserUpdateInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var user models.User
	// Use First to fetch the user and handle "not found"
	result := services.DB.Where("id= ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Usuário não encontrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
		}
		return
	}

	user.DisplayName = userInput.DisplayName
	user.Email = userInput.Email
	user.IsSubscribed = userInput.IsSubscribed
	user.ExpirationDate = userInput.ExpirationDate

	// Save the updated user
	result = services.DB.Save(&user)
	if result.Error != nil {
		// Handle unique constraint violations (e.g., return 409)
		c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuário atualizado com sucesso", "data": user})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	result := services.DB.Where("id=?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Usuário não encontrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
		}
		return
	}

	result = services.DB.Delete(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuário deletado com sucesso"})
}

func UpdateAllTokens() {
	err := services.DB.
		// temporarily allow an update without WHERE
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Model(&models.User{}).
		UpdateColumns(map[string]any{
			"token": gorm.Expr("gen_random_uuid()"),
		}).Error

	if err != nil {
		fmt.Println("error updating tokens:", err)
		return
	}
	fmt.Println("successfully updated tokens")
}