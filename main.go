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

	router := gin.Default()
	router.SetTrustedProxies([]string{"172.19.0.0/16"})

	// Connecting to database
	if err := services.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations explicitly
	if err := services.DB.AutoMigrate(
		&models.User{},
		&models.Payment{},
		&models.Benefit{},
		&models.AlojaUser{},
		&models.Revenue{},
	); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Initialize cron job to reset all user tokens every day at midnight (horário de brasília)
	c := cron.New()
	c.AddFunc("3 0 * * *", controllers.UpdateAllTokens)
	c.Start()

	// Subscription route (to fetch if the user has a plan)
	router.GET("/go/subscription/:token", controllers.GetSubscription)

	// Benefit route (to fetch which benefits or plans there are)
	router.GET("/go/benefits", controllers.GetBenefits)
	router.GET("/go/plans", controllers.GetPlans)

	// Login route
	router.POST("/go/login", controllers.Login)

	// Mercado pago payment checking
	router.POST("/go/payment/confirm", middlewares.CheckPayment, controllers.UpdateSubscription)

	// Authenticated router
	userRouter := router.Group("/go/auth")
	userRouter.Use(middlewares.CheckAuth)
	{
		userRouter.GET("/profile", controllers.GetUserProfile)
		userRouter.POST("/payment/create", controllers.CreatePayment)
		userRouter.GET("/payment", controllers.GetPayment)
	}

	// Admin user router
	adminRouter := router.Group("/go/admin")
	adminRouter.Use(middlewares.OriginWhitelist, middlewares.CheckAdmin)
	{
		adminRouter.PUT("/user/:id", controllers.UpdateUserProfile)
		adminRouter.GET("/user/all", controllers.GetAllUsers)
		adminRouter.DELETE("/user/:id", controllers.DeleteUser)

		adminRouter.PUT("/benefit/:id", controllers.UpdateBenefit)
		adminRouter.POST("/benefit/create", controllers.CreateBenefit)
		adminRouter.DELETE("/benefit/:id", controllers.DeleteBenefit)

		adminRouter.GET("/alojaUser/all", controllers.GetAlojaUsers)
		adminRouter.POST("/alojaUser/create", controllers.CreateMultipleAlojaUsers)
		adminRouter.DELETE("/alojaUser/:id", controllers.DeleteAlojaUser)

		adminRouter.GET("/revenue/all", controllers.GetRevenues)
	}

	router.Run()
}
