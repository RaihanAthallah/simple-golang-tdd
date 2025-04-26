package main

import (
	"fmt"
	"log"
	"net/http"
	"simple-golang-tdd/controller"
	"simple-golang-tdd/repository"
	"simple-golang-tdd/routes"
	"simple-golang-tdd/service"

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
// @BasePath /api/v1

func main() {
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

	authRepository := repository.NewAuthRepository()
	authService := service.NewAuthService(authRepository)
	authController := controller.NewAuthController(authService)

	authGroup := router.Group("/api/v1")
	authGroup.Use()
	{
		routes.SetupAuthRoutes(authGroup, authController)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	port := "8080"
	// Menjalankan server
	fmt.Printf("Server is running on :%s...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))

}
