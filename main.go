package main

import (
	"caaso/controllers"
	"caaso/middlewares"
	"caaso/models"
	"caaso/services"
	"log"

	"github.com/gin-gonic/gin"
)

func init() {
	services.LoadEnvs()
	services.InitFirebase()
	services.PaymentService()
}

func main() {

	routes := gin.Default()
	routes.SetTrustedProxies([]string{"172.19.0.0/16"})

	// Connecting to database
	if err := services.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations explicitly
	if err := services.DB.AutoMigrate(&models.User{}, &models.Payment{}, &models.Benefit{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Subscription route (to fetch if the user has a plan)
	routes.GET("/go/subscription/:nusp", controllers.GetSubscription)

	// Benefit route (to fetch which benefits there are)
	routes.GET("/go/benefits", controllers.GetBenefits)

	// Plan route (to fetch which plan models there are)
	routes.GET("/go/plans", controllers.GetPlans)

	// User routes
	routes.POST("/go/user/create", controllers.CreateUser)
	routes.POST("/go/user/login", middlewares.GoogleAuth, controllers.GenerateJwtToken)
	routes.GET("/go/user/profile", middlewares.CheckAuth, controllers.GetUserProfile)

	// Payment routes
	routes.POST("/go/payment/create", middlewares.CheckAuth, controllers.CreatePayment)
	routes.POST("/go/payment/confirm", middlewares.CheckPayment, controllers.UpdateSubscription)
	routes.GET("/go/payment", middlewares.CheckAuth, controllers.GetPayment)

	// Admin user routes
	adminRoutes := routes.Group("/go/admin")
	adminRoutes.Use(middlewares.OriginWhitelist, middlewares.CheckAdmin)
	{
		routes.PUT("/go/admin/user/:id", controllers.UpdateUserProfile)
		routes.GET("/go/admin/user/all", controllers.GetAllUsers)
		routes.POST("/go/admin/benefit/create", controllers.CreateBenefit)
		routes.PUT("/go/admin/benefit/:id", controllers.UpdateBenefit)
		routes.DELETE("/go/admin/benefit/:id", controllers.DeleteBenefit)
		routes.DELETE("/go/admin/user/:id", controllers.DeleteUser)
	}

	routes.Run()
}
