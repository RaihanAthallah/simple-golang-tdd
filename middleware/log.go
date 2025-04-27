package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"simple-golang-tdd/model"                          // ganti dengan import path kamu
	historyRepo "simple-golang-tdd/repository/history" // ganti dengan import path kamu
	"simple-golang-tdd/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// responseWriterWrapper is used to capture the status code of the response
type responseWriterWrapper struct {
	gin.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// readRequestBody reads the request body
func readRequestBody(c *gin.Context) ([]byte, error) {
	if c.Request.Method == http.MethodGet || c.Request.Body == nil {
		return nil, nil
	}
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore body for further use
	return bodyBytes, nil
}

// buildHistory creates the History object from the request and response
func buildHistory(c *gin.Context, statusCode int, bodyBytes []byte, username string) model.History {
	details := map[string]interface{}{
		"headers": c.Request.Header,
	}
	if len(bodyBytes) > 0 {
		details["payload"] = json.RawMessage(bodyBytes)
	}

	return model.History{
		ID:         uuid.New().String(),
		CustomerID: username, // Set CustomerID with extracted username
		Status:     http.StatusText(statusCode),
		Action:     c.Request.URL.Path,
		Details:    details,
		Timestamp:  time.Now().Format(time.RFC3339),
	}
}

// HistoryLoggerMiddleware is the middleware wrapper for the HistoryLoggerHandler
func HistoryLoggerMiddleware(repo historyRepo.HistoryRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read the request body (before processing)
		bodyBytes, err := readRequestBody(c)
		if err != nil {
			// Optionally, handle error reading the body
			fmt.Println("Error reading request body:", err)
		}

		// Wrap the ResponseWriter to capture the status code
		rw := &responseWriterWrapper{
			ResponseWriter: c.Writer,
			statusCode:     http.StatusOK,
		}

		// Set the custom ResponseWriter
		c.Writer = rw

		// Call the next middleware/handler
		c.Next()

		// After the request has been processed, log the history
		// Check if the token is present and extract the username
		username, err := utils.ExtractTokenFromHeader(c)
		if err != nil {
			// Token is not available (for unauthenticated requests), log as anonymous or skip logging
			username = "anonymous" // You can choose another way to identify unauthenticated requests
			// Optionally, log the error
			fmt.Println("Error extracting token:", err)
		}

		// Build the history object
		history := buildHistory(c, rw.statusCode, bodyBytes, username)

		// Save the history in the repository
		if _, err := repo.UpdateHistory(history); err != nil {
			// Optional: log error if history cannot be stored
			fmt.Println("Error saving history:", err)
		}
	}
}
