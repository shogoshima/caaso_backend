package controllers

import (
	"caaso/models"
	"caaso/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Plan struct {
	UserType models.UserTypes `json:"userType"`
	PlanType models.PlanTypes `json:"planType"`
	Amount   float64          `json:"amount"`
}

var amounts = map[models.UserTypes]map[models.PlanTypes]float64{
	models.Aloja: {
		models.Monthly: 0.01,
		models.Yearly:  0.01,
	},
	models.Grad: {
		models.Monthly: 0.01,
		models.Yearly:  0.01,
	},
	models.PostGrad: {
		models.Monthly: 0.01,
		models.Yearly:  0.01,
	},
	models.Other: {
		models.Monthly: 0.01,
		models.Yearly:  0.01,
	},
}

func GetPlans(c *gin.Context) {
	var list []Plan
	for u, plans := range amounts {
		for p, amt := range plans {
			list = append(list, Plan{
				UserType: u,
				PlanType: p,
				Amount:   amt,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Planos carregados", "plans": list})

}

func CreatePayment(c *gin.Context) {

	user, _ := c.Get("currentUser")
	User := user.(models.User)

	var input models.PaymentInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Update the User.type column using the String() of the enum:
	userTypeStr := input.UserType.String()
	if err := services.DB.
		Model(&User).
		Where("id = ?", User.ID).
		Update("type", userTypeStr).
		Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Lookup the amount:
	amount := amounts[input.UserType][input.PlanType]

	resource, err := services.CreatePayment(amount, User.Email, User.ID, input.PlanType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// remove any previous payments for this user
	if err := services.DB.
		Where("user_id = ?", User.ID).
		Delete(&models.Payment{}).
		Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Create new payment
	var payment = models.Payment{
		ID:     resource.ID,
		UserId: User.ID,
		QrCode: resource.PointOfInteraction.TransactionData.QRCode,
		Amount: resource.TransactionAmount,
	}

	result := services.DB.Create(&payment)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pagamento criado com sucesso", "payment": payment})

}

func GetPayment(c *gin.Context) {

	user, _ := c.Get("currentUser")
	User := user.(models.User)

	var payment models.Payment
	result := services.DB.Where("user_id=?", User.ID).First(&payment)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
		return
	}

	resource, err := services.GetPaymentFromId(payment.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	switch status := resource.Status; status {
	case "in_process":
		c.JSON(http.StatusForbidden, gin.H{"message": "O seu pagamento está em processamento"})
	case "cancelled", "rejected":
		result := services.DB.Delete(&payment)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
			return
		}
		c.JSON(http.StatusForbidden, gin.H{"message": "O seu pagamento foi cancelado. Recarregue e tente novamente"})
	case "pending":
		// pagamento com "isPaid" false
		c.JSON(http.StatusOK, gin.H{"message": "Pagamento pendente", "payment": payment})
	case "approved":
		// pagamento com "isPaid" true
		result := services.DB.Delete(&payment)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Pagamento aprovado", "payment": payment})
		return
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Status desconhecido."})
		return
	}
}
