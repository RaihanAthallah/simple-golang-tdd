package main

import (
	"fmt"
	"log"
	"net/http"
	"simple-golang-tdd/controller"
	"simple-golang-tdd/repository"
	"simple-golang-tdd/routes"
	"simple-golang-tdd/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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

	authGroup := router.Group("/v1")
	authGroup.Use()
	{
		routes.SetupAuthRoutes(authGroup, authController)
	}

	port := "8080"
	// Menjalankan server
	fmt.Printf("Server is running on :%s...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))

}
