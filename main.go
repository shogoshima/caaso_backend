package main

import (
	"caaso/controllers"
	"caaso/middlewares"
	"caaso/models"
	"caaso/services"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
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

	// Initialize cron job to reset all user tokens every day at midnight (horário de brasília)
	c := cron.New()
	c.AddFunc("3 0 * * *", controllers.UpdateAllTokens)
	c.Start()

	// Subscription route (to fetch if the user has a plan)
	routes.GET("/go/subscription/:token", controllers.GetSubscription)

	// Benefit route (to fetch which benefits or plans there are)
	routes.GET("/go/benefits", controllers.GetBenefits)
	routes.GET("/go/plans", controllers.GetPlans)

	// Login route
	routes.POST("/go/login", middlewares.GoogleAuth, controllers.GenerateJwtToken)

	// Mercado pago payment checking
	routes.POST("/go/payment/confirm", middlewares.CheckPayment, controllers.UpdateSubscription)

	// Authenticated routes
	userRoutes := routes.Group("/go/auth")
	userRoutes.Use(middlewares.CheckAuth)
	{
		userRoutes.GET("/profile", middlewares.CheckAuth, controllers.GetUserProfile)
		userRoutes.POST("/payment/create", middlewares.CheckAuth, controllers.CreatePayment)
		userRoutes.GET("/payment", middlewares.CheckAuth, controllers.GetPayment)
	}

	// Admin user routes
	adminRoutes := routes.Group("/go/admin")
	adminRoutes.Use(middlewares.OriginWhitelist, middlewares.CheckAdmin)
	{
		adminRoutes.PUT("/user/:id", controllers.UpdateUserProfile)
		adminRoutes.GET("/user/all", controllers.GetAllUsers)
		adminRoutes.POST("/benefit/create", controllers.CreateBenefit)
		adminRoutes.PUT("/benefit/:id", controllers.UpdateBenefit)
		adminRoutes.DELETE("/benefit/:id", controllers.DeleteBenefit)
		adminRoutes.DELETE("/user/:id", controllers.DeleteUser)
	}

	routes.Run()
}
