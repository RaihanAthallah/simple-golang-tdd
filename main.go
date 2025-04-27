package main

import (
	"fmt"
	"log"
	"net/http"
	AuthController "simple-golang-tdd/controller/auth"
	CustomerController "simple-golang-tdd/controller/customer"
	"simple-golang-tdd/middleware"
	"simple-golang-tdd/routes"

	CustomerRepository "simple-golang-tdd/repository/customer"
	HistoryRepository "simple-golang-tdd/repository/history"
	MerchantRepository "simple-golang-tdd/repository/merchant"

	AuthService "simple-golang-tdd/service/auth"
	CustomerService "simple-golang-tdd/service/customer"

	_ "simple-golang-tdd/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title My Gin API
// @version 1.0
// @description This is a sample server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

func main() {

	// source data path
	const customerDataPath = "./data/customers.json"
	const historyDataPath = "./data/histories.json"
	const merhacntDataPath = "./data/merchants.json"

	// Membuat router Gin
	router := gin.Default()

	// Configure CORS
	router.Use(
		cors.New(
			cors.Config{
				AllowOrigins:     []string{"*"},
				AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowHeaders:     []string{"Content-Type", "Authorization", "token"}, // Add the "token" header here
				AllowCredentials: true,
			},
		),
	)

	customerhRepository, err := CustomerRepository.NewCustomerRepository(customerDataPath)
	if err != nil {
		log.Fatalf("Failed to create customer repository: %v", err)
	}
	merchantRepository, err := MerchantRepository.NewMerchantRepository(merhacntDataPath)
	if err != nil {
		log.Fatalf("Failed to create merchant repository: %v", err)
	}
	historyRepository, err := HistoryRepository.NewHistoryRepository(historyDataPath)
	if err != nil {
		log.Fatalf("Failed to create history repository: %v", err)
	}

	authService := AuthService.NewAuthService(customerhRepository)
	customerService := CustomerService.NewCustomerService(customerhRepository, merchantRepository)

	authController := AuthController.NewAuthController(authService)
	customerController := CustomerController.NewCustomerController(customerService)

	router.Use(middleware.HistoryLoggerMiddleware(historyRepository))
	noAuthGroup := router.Group("/user/v1")
	noAuthGroup.Use()
	{
		routes.SetupAuthRoutes(noAuthGroup, authController)
	}

	// Define the authGroup (authenticated routes)
	authGroup := router.Group("/api/v1/")
	authGroup.Use(middleware.JWTAuthMiddleware()) // Use authentication middleware here
	{

		routes.SetupCustomerRoutes(authGroup, customerController)
		// Add routes that require authentication (e.g., user profile, protected resources)
		// Example:
		// authGroup.GET("/user", userController.GetUser)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	port := "8080"
	// Menjalankan server
	fmt.Printf("Server is running on :%s...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))

}
