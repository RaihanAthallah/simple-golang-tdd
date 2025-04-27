package controller

import (
	"simple-golang-tdd/dto"
	customerService "simple-golang-tdd/service/customer"
	"simple-golang-tdd/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// CustomerController handles authentication-related operations
type CustomerController struct {
	customerService  customerService.CustomerService
	customervalidate *validator.Validate
}

func NewCustomerController(service customerService.CustomerService) *CustomerController {
	validate := validator.New()
	return &CustomerController{customerService: service, customervalidate: validate}
}

// Paymentgodoc
// @Summary      Customer Payment to Merchant
// @Description  Customer payment reduces balance and send to merchant
// @Tags         Customer
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer Token"
// @Param        body  body  dto.PaymentRequest  true  "Payment Request"
// @Success      200  {object} dto.SuccessResponse  ""payment successful"
// @Failure      400  {object} dto.ErrorResponse
// @Failure      401  {object} dto.ErrorResponse
// @Failure      500  {object} dto.ErrorResponse
// @Router       /api/v1/customer/payment [post]
func (cc *CustomerController) Payment(c *gin.Context) {
	var paymentRequest dto.PaymentRequest
	err := c.BindJSON(&paymentRequest)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid request body")
		return
	}

	if err := cc.customervalidate.Struct(paymentRequest); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	// Get the username from the context (assuming it's set in a middleware)
	username, exists := c.Get("username")
	if !exists {
		utils.ErrorResponse(c, 401, "Unauthorized: Invalid username in token")
		return
	}

	// You should cast the username from context to a string
	strUsername, _ := username.(string)

	customer, err := cc.customerService.Payment(paymentRequest, strUsername)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "payment successful", customer)
}
