package controller_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"simple-golang-tdd/controller"
	"simple-golang-tdd/dto"
	"simple-golang-tdd/utils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Setup Router ---
func setupRouter(service *MockAuthService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	authCtrl := controller.NewAuthController(service)
	r.POST("/v1/customer/login", authCtrl.Login)   // Fixed path
	r.POST("/v1/customer/logout", authCtrl.Logout) // Fixed path
	r.POST("/v1/customer/refresh-token", authCtrl.RefreshToken)

	return r
}

func newRecorderAndRouter(mockService *MockAuthService) (*httptest.ResponseRecorder, http.Handler) {
	return httptest.NewRecorder(), setupRouter(mockService)
}

// --- TEST CASES ---
func TestLogin_Success(t *testing.T) {
	mockService := new(MockAuthService)

	fakeUserData := dto.UserCredentials{
		Username: "user",
		Password: "correct_password",
	}

	fakeAuthResponse := dto.AuthResponse{
		AccessToken:  "dummy_access_token",
		RefreshToken: "dummy_refresh_token",
	}

	fakeResponseMessage := dto.SuccessResponse{
		Status:  200,
		Message: "login successful",
		Data:    fakeAuthResponse,
	}

	mockService.On("Login", fakeUserData).Return(fakeAuthResponse, nil)

	recorder, router := newRecorderAndRouter(mockService)

	payload := dto.UserCredentials{
		Username: "user",
		Password: "correct_password",
	}

	req, err := utils.NewJSONRequest(http.MethodPost, "/v1/customer/login", payload)
	require.NoError(t, err) // use require to fail immediately if request creation fails

	router.ServeHTTP(recorder, req)

	expectedResponseBody := utils.MarshalJSON(t, fakeResponseMessage)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.JSONEq(t, expectedResponseBody, recorder.Body.String())
	mockService.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	mockService := new(MockAuthService)
	userCreds := dto.UserCredentials{Username: "user", Password: "wrong_password"}
	mockService.On("Login", userCreds).Return(dto.AuthResponse{}, errors.New("Invalid credentials"))

	rec, router := newRecorderAndRouter(mockService)

	req, _ := utils.NewJSONRequest("POST", "/v1/customer/login", userCreds)
	router.ServeHTTP(rec, req)

	expected := dto.ErrorResponse{Status: 401, Message: "invalid credentials"}
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.JSONEq(t, utils.MarshalJSON(t, expected), rec.Body.String())
	mockService.AssertExpectations(t)
}

func TestLogin_MissingField(t *testing.T) {
	rec, router := newRecorderAndRouter(new(MockAuthService))

	payload := map[string]string{"username": "user"} // missing password
	req, _ := utils.NewJSONRequest("POST", "/v1/customer/login", payload)
	router.ServeHTTP(rec, req)

	expected := dto.ErrorResponse{Status: 400, Message: "invalid request body"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, utils.MarshalJSON(t, expected), rec.Body.String())
}

func TestLogin_InvalidFieldType(t *testing.T) {
	rec, router := newRecorderAndRouter(new(MockAuthService))

	req, _ := http.NewRequest("POST", "/v1/customer/login", bytes.NewBuffer([]byte(`{"username":"user","password":123}`)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	expected := dto.ErrorResponse{Status: 400, Message: "invalid request body"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, utils.MarshalJSON(t, expected), rec.Body.String())
}

func TestLogin_InvalidJSONFormat(t *testing.T) {
	rec, router := newRecorderAndRouter(new(MockAuthService))

	req, _ := http.NewRequest("POST", "/v1/customer/login", bytes.NewBuffer([]byte(`{"username":"user",`)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	expected := dto.ErrorResponse{Status: 400, Message: "invalid request body"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, utils.MarshalJSON(t, expected), rec.Body.String())
}

func TestLogin_ExtraFields(t *testing.T) {
	mockService := new(MockAuthService)
	userCreds := dto.UserCredentials{Username: "user", Password: "correct_password"}
	authResp := dto.AuthResponse{AccessToken: "dummy_access_token", RefreshToken: "dummy_refresh_token"}

	mockService.On("Login", userCreds).Return(authResp, nil)

	rec, router := newRecorderAndRouter(mockService)

	payload := map[string]interface{}{
		"username":    "user",
		"password":    "correct_password",
		"extra_field": "ignored",
	}
	req, _ := utils.NewJSONRequest("POST", "/v1/customer/login", payload)
	router.ServeHTTP(rec, req)

	expected := dto.SuccessResponse{
		Status:  200,
		Message: "login successful",
		Data:    authResp,
	}
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, utils.MarshalJSON(t, expected), rec.Body.String())
	mockService.AssertExpectations(t)
}

func TestLogin_ValidationError(t *testing.T) {
	rec, router := newRecorderAndRouter(new(MockAuthService))

	payload := map[string]string{"username": "", "password": ""}
	req, _ := utils.NewJSONRequest("POST", "/v1/customer/login", payload)
	router.ServeHTTP(rec, req)

	expected := dto.ErrorResponse{Status: 400, Message: "invalid request body"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, utils.MarshalJSON(t, expected), rec.Body.String())
}

func TestLogout_Success(t *testing.T) {
	mockService := new(MockAuthService)
	mockService.On("Logout", "Bearer valid_token").Return(nil)

	rec, router := newRecorderAndRouter(mockService)

	req, _ := http.NewRequest("POST", "/v1/customer/logout", nil)
	req.Header.Set("Authorization", "Bearer valid_token")
	router.ServeHTTP(rec, req)

	expected := dto.SuccessResponse{Status: 200, Message: "logout successful"}
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, utils.MarshalJSON(t, expected), rec.Body.String())
	mockService.AssertExpectations(t)
}

func TestLogout_MissingToken(t *testing.T) {
	rec, router := newRecorderAndRouter(new(MockAuthService))

	req, _ := http.NewRequest("POST", "/v1/customer/logout", nil)
	router.ServeHTTP(rec, req)

	expected := dto.ErrorResponse{Status: 401, Message: "please provide a token"}
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.JSONEq(t, utils.MarshalJSON(t, expected), rec.Body.String())
}

func TestLogout_ServiceError(t *testing.T) {
	mockService := new(MockAuthService)
	mockService.On("Logout", "Bearer valid_token").Return(errors.New("service error"))

	rec, router := newRecorderAndRouter(mockService)

	req, _ := http.NewRequest("POST", "/v1/customer/logout", nil)
	req.Header.Set("Authorization", "Bearer valid_token")
	router.ServeHTTP(rec, req)

	expected := dto.ErrorResponse{Status: 500, Message: "internal server error"}
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.JSONEq(t, utils.MarshalJSON(t, expected), rec.Body.String())
	mockService.AssertExpectations(t)
}

func TestRefreshToken_Success(t *testing.T) {
	mockService := new(MockAuthService)

	refreshToken, _ := utils.GenerateRefreshToken("user123")

	// Setup payload
	payload := dto.RefreshToken{
		RefreshToken: refreshToken,
	}

	// Setup mock return value
	fakeAccessToken := dto.AccessTokenResponse{
		AccessToken: "new_dummy_token",
	}

	fakeResponseMessage := dto.SuccessResponse{
		Status:  200,
		Message: "token refreshed successfully",
		Data:    fakeAccessToken,
	}

	mockService.On("RefreshToken", payload).Return(fakeAccessToken, nil)

	recorder, router := newRecorderAndRouter(mockService)

	req, err := utils.NewJSONRequest(http.MethodPost, "/v1/customer/refresh-token", payload)
	require.NoError(t, err)

	router.ServeHTTP(recorder, req)

	expectedResponseBody := utils.MarshalJSON(t, fakeResponseMessage)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.JSONEq(t, expectedResponseBody, recorder.Body.String())
	mockService.AssertExpectations(t)
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	mockService := new(MockAuthService)

	payload := dto.RefreshToken{RefreshToken: "invalid_token"}

	// Mocking RefreshToken to return an error
	mockService.On("RefreshToken", payload).Return(dto.AccessTokenResponse{}, errors.New("invalid or expired refresh token"))

	recorder, router := newRecorderAndRouter(mockService)

	req, err := utils.NewJSONRequest(http.MethodPost, "/v1/customer/refresh-token", payload)
	require.NoError(t, err)

	router.ServeHTTP(recorder, req)

	expected := dto.ErrorResponse{
		Status:  500,
		Message: "internal server error", // Sesuai AuthController sekarang, kalo RefreshToken error, dia kasih 500
	}

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.JSONEq(t, utils.MarshalJSON(t, expected), recorder.Body.String())
	mockService.AssertExpectations(t)
}


func TestRefreshToken_MissingToken(t *testing.T) {
	rec, router := newRecorderAndRouter(nil)

	payload := dto.RefreshToken{}
	req, _ := utils.NewJSONRequest("POST", "/v1/customer/refresh-token", payload)
	router.ServeHTTP(rec, req)

	expected := dto.ErrorResponse{Status: 400, Message: "invalid request body"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, utils.MarshalJSON(t, expected), rec.Body.String())
}

func TestRefreshToken_MalformedJSON(t *testing.T) {
	rec, router := newRecorderAndRouter(nil)

	req, _ := http.NewRequest("POST", "/v1/customer/refresh-token", bytes.NewBuffer([]byte(`{"refresh_token":`)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	expected := dto.ErrorResponse{Status: 400, Message: "invalid request body"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, utils.MarshalJSON(t, expected), rec.Body.String())
}
