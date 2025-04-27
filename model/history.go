package model

type History struct {
	ID         string      `json:"id"`
	CustomerID string      `json:"customer_id"`
	Status     string      `json:"status"`
	Action     string      `json:"action"`
	Details    interface{} `json:"details"` // Dynamic details using a map
	Timestamp  string      `json:"timestamp"`
}
