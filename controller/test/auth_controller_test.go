package controller_test

// Saat password salah, error message akan dikembalikan
// Payload yang dikirimkan ke server adalah JSON
// Format Payload JSON adalah salah maka akan error message akan dikembalikan
// Salah satu field pada payload JSON tidak ada, error message akan dikembalikan
// Value dari field pada payload JSON tidak sesuai dengan tipe data yang diharapkan, error message akan dikembalikan
// Payload JSON yang valid, access token dan refresh token akan dikembalikan

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	controller "simple-golang-tdd/controller"
)

// --- Setup Router ---
func setupRouter(service *MockAuthService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	authCtrl := controller.NewAuthController(service)
	r.POST("/login", authCtrl.Login)

	return r
}

// --- TEST CASES ---

func TestLogin_ValidPayload(t *testing.T) {
	mockService := new(MockAuthService)
	mockService.On("Login", "user", "correct_password").Return("dummy_token", nil)

	router := setupRouter(mockService)

	payload := map[string]string{"username": "user", "password": "correct_password"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "dummy_token")
	mockService.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	mockService := new(MockAuthService)
	mockService.On("Login", "user", "wrong_password").Return("", errors.New("Invalid credentials"))

	router := setupRouter(mockService)

	payload := map[string]string{"username": "user", "password": "wrong_password"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "Invalid credentials")
	mockService.AssertExpectations(t)
}

func TestLogin_MissingField(t *testing.T) {
	mockService := new(MockAuthService)
	router := setupRouter(mockService)

	payload := map[string]string{"username": "user"} // password hilang
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Invalid request")
}

func TestLogin_InvalidFieldType(t *testing.T) {
	mockService := new(MockAuthService)
	router := setupRouter(mockService)

	body := []byte(`{"username": "user", "password": 123}`)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Invalid request")
}

func TestLogin_InvalidJSONFormat(t *testing.T) {
	mockService := new(MockAuthService)
	router := setupRouter(mockService)

	body := []byte(`{"username": "user",`)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Invalid request")
}
