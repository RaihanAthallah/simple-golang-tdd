package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"simple-golang-tdd/model"
	"time"
)

// logRequest logs a new request with dynamic details
func logRequest(id, customerID, action string, details interface{}) {
	// Convert details to JSON string
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		fmt.Printf("Error marshalling details: %v\n", err)
		return
	}

	// Create the log entry using the JSON-encoded details
	logEntry := model.History{
		ID:         id,
		CustomerID: customerID,
		Action:     action,
		Details:    string(detailsJSON), // Store as a string
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}

	// Convert the log entry to JSON for structured logging
	logData, err := json.Marshal(logEntry)
	if err != nil {
		fmt.Printf("Error marshalling log data: %v\n", err)
		return
	}

	// Log the entry (you can replace this with a more advanced logger if needed)
	log.Println(string(logData))
}
