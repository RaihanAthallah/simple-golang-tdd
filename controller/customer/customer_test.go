package controller

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"simple-golang-tdd/dto"
	"simple-golang-tdd/middleware"
	"simple-golang-tdd/model"
	"simple-golang-tdd/utils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Setup Router ---
func setupRouter(service *MockCustomerService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.Use(middleware.JWTAuthMiddleware())

	customerCtrl := NewCustomerController(service)
	r.POST("/v1/customer/payment", customerCtrl.Payment) // Fixed path

	return r
}

func newRecorderAndRouter(mockService *MockCustomerService) (*httptest.ResponseRecorder, http.Handler) {
	return httptest.NewRecorder(), setupRouter(mockService)
}

// --- TEST CASES ---
func TestPayment_Success(t *testing.T) {
	mockService := new(MockCustomerService)

	// Define fake payment request and expected response
	fakePaymentRequest := dto.PaymentRequest{
		MerchantID: "merchant123",
		Amount:     100.0,
	}

	fakeUsername := "user"

	fakeUserData := model.Customer{
		ID:       "1",
		Username: "user",
		Password: "password", // password matches
		Balance:  1000.0,
	}

	fakeResponseMessage := dto.SuccessResponse{
		Status:  200,
		Message: "payment successful",
		Data:    fakeUserData,
	}

	// Mock the Payment method on the service
	mockService.On("Payment", fakePaymentRequest, fakeUsername).Return(fakeUserData, nil)

	// Create a JWT token for the user
	token, err := utils.GenerateAccessToken(fakeUsername) // Assuming this is your JWT generation function
	require.NoError(t, err)

	// Create a new request with the fake payment data
	rec, router := newRecorderAndRouter(mockService)
	req, err := utils.NewJSONRequest(http.MethodPost, "/v1/customer/payment", fakePaymentRequest)
	require.NoError(t, err)

	// Add the JWT token to the request header
	req.Header.Set("Authorization", "Bearer "+token)

	// Pass the request with the updated context (with token validation)
	ctx := req.Context()
	ctx = context.WithValue(ctx, "username", fakeUsername) // Set the username value in context
	req = req.WithContext(ctx)                             // Update the request with the new context

	// Make the request
	router.ServeHTTP(rec, req)

	// Marshal the expected response body
	expectedResponseBody := utils.MarshalJSON(t, fakeResponseMessage)

	// Assertions
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, expectedResponseBody, rec.Body.String())
	mockService.AssertExpectations(t)
}

func TestPayment_MissingField(t *testing.T) {
	mockService := new(MockCustomerService)
	rec, router := newRecorderAndRouter(mockService)

	payload := map[string]string{"merchant_id": "merchant123"} // missing amount

	// Create a JWT token for the user
	token, err := utils.GenerateAccessToken("user")
	require.NoError(t, err)

	// Create the request with the Authorization header
	req, _ := utils.NewJSONRequest("POST", "/v1/customer/payment", payload)
	req.Header.Set("Authorization", "Bearer "+token)

	// Make the request
	router.ServeHTTP(rec, req)

	expected := dto.ErrorResponse{Status: 400, Message: "invalid request body"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, utils.MarshalJSON(t, expected), rec.Body.String())
}

func TestPayment_InvalidFieldType(t *testing.T) {
	mockService := new(MockCustomerService)
	rec, router := newRecorderAndRouter(mockService)

	payload := map[string]interface{}{"merchant_id": "merchant123", "amount": "invalid_amount"} // invalid amount type

	// Create a JWT token for the user
	token, err := utils.GenerateAccessToken("user")
	require.NoError(t, err)

	// Create the request with the Authorization header
	req, _ := utils.NewJSONRequest("POST", "/v1/customer/payment", payload)
	req.Header.Set("Authorization", "Bearer "+token)

	// Make the request
	router.ServeHTTP(rec, req)

	expected := dto.ErrorResponse{Status: 400, Message: "invalid request body"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, utils.MarshalJSON(t, expected), rec.Body.String())
}

func TestPayment_InvalidToken(t *testing.T) {
	mockService := new(MockCustomerService)
	rec, router := newRecorderAndRouter(mockService)

	payload := map[string]interface{}{"merchant_id": "merchant123", "amount": 100.0}

	// Create an invalid JWT token
	invalidToken := "invalid_token"

	// Create the request with the Authorization header
	req, _ := utils.NewJSONRequest("POST", "/v1/customer/payment", payload)
	req.Header.Set("Authorization", "Bearer "+invalidToken)

	// Make the request
	router.ServeHTTP(rec, req)

	expected := dto.ErrorResponse{Status: 401, Message: "Unauthorized: token contains an invalid number of segments"}
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.JSONEq(t, utils.MarshalJSON(t, expected), rec.Body.String())
}

func TestPayment_EmptyToken(t *testing.T) {
	mockService := new(MockCustomerService)
	rec, router := newRecorderAndRouter(mockService)

	payload := map[string]interface{}{"merchant_id": "merchant123", "amount": 100.0}

	// Create the request without the Authorization header
	req, _ := utils.NewJSONRequest("POST", "/v1/customer/payment", payload)

	// Make the request
	router.ServeHTTP(rec, req)

	expected := dto.ErrorResponse{Status: 401, Message: "Authorization header is missing"}
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.JSONEq(t, utils.MarshalJSON(t, expected), rec.Body.String())
}

func TestPayment_ServiceFailure(t *testing.T) {
	mockService := new(MockCustomerService)
	rec, router := newRecorderAndRouter(mockService)

	// Define valid payload
	fakePaymentRequest := dto.PaymentRequest{
		MerchantID: "merchant123",
		Amount:     100.0,
	}

	fakeUsername := "user"

	// Mock the Payment service call to return an error
	mockService.On("Payment", fakePaymentRequest, fakeUsername).Return(model.Customer{}, errors.New("insufficient balance"))

	// Create a JWT token for the user
	token, err := utils.GenerateAccessToken(fakeUsername)
	require.NoError(t, err)

	// Create the request with the Authorization header
	req, _ := utils.NewJSONRequest("POST", "/v1/customer/payment", fakePaymentRequest)
	req.Header.Set("Authorization", "Bearer "+token)

	// Make the request
	router.ServeHTTP(rec, req)

	expected := dto.ErrorResponse{
		Status:  400,
		Message: "insufficient balance",
	}

	// Assertions
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, utils.MarshalJSON(t, expected), rec.Body.String())
	mockService.AssertExpectations(t)
}
