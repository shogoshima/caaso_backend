package controllers

import (
	"caaso/models"
	"caaso/services"
	"errors"
	"net/http"
	"os"
	"time"

	"firebase.google.com/go/v4/auth"
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

func CreateUser(c *gin.Context) {

	var userInput models.UserInput

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var userFound models.User
	result := services.DB.Where("id=?", userInput.ID).Find(&userFound)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
		return
	}

	user := models.User{
		ID:          userInput.ID,
		Nusp:        userInput.Nusp,
		DisplayName: userInput.DisplayName,
		Email:       userInput.Email,
		PhotoUrl:    userInput.PhotoUrl,
	}

	result = services.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuário criado com sucesso", "user": user})

}

func VerifyIDToken(c *gin.Context, idToken string) (*auth.Token, error) {
	token, err := services.AuthClient.VerifyIDToken(c.Request.Context(), idToken)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func GenerateJwtToken(c *gin.Context) {

	userID, _ := c.Get("currentUserID")

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Hour * 24 * 14).Unix(),
	})

	sessionToken, err := generateToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao gerar token, tente novamente"})
		return
	}

	var user models.User
	result := services.DB.Where("id=?", userID).Find(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
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

	// Check if Nusp is being changed and update NuspModified
	if user.Nusp != userInput.Nusp {
		user.NuspModified = true
	}

	println(userInput.IsSubscribed)

	// Update fields from input
	user.Nusp = userInput.Nusp
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