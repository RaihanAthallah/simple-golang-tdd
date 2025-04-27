package routes

import (
	controller "simple-golang-tdd/controller/auth"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.RouterGroup, authController *controller.AuthController) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/logout", authController.Logout)
		authGroup.POST("/refresh-token", authController.RefreshToken)
	}
}
