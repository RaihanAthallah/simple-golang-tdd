package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// NewJSONRequest creates a new JSON HTTP request with the given payload.
func NewJSONRequest(method, url string, payload interface{}) (*http.Request, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// MarshalJSON safely marshals an object for tests and fails immediately if it fails.
func MarshalJSON(t *testing.T, v interface{}) string {
	t.Helper()
	bytes, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("failed to marshal JSON: %v", err)
	}
	return string(bytes)
}

func LoadJSONFile(filePath string, target interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func SaveJSONFile(filePath string, data interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// extractTokenFromHeader extracts the token from the Authorization header
func ExtractTokenFromHeader(c *gin.Context) (string, error) {
	// Get the "Authorization" header
	authHeader := c.GetHeader("Authorization")

	// Check if the header is valid
	if authHeader == "" {
		return "", fmt.Errorf("Authorization header is missing")
	}

	// Authorization header is expected to be in the format "Bearer <token>"
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", fmt.Errorf("Invalid Authorization header format")
	}

	// Return the token part
	return tokenParts[1], nil
}
