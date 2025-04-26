package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
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
