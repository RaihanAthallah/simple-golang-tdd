package routes

import (
	controller "simple-golang-tdd/controller/customer"

	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(router *gin.RouterGroup, customerController *controller.CustomerController) {
	customerGroup := router.Group("/customer")
	{
		customerGroup.POST("/payment", customerController.Payment)
	}
}
