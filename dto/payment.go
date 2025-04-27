package dto

type PaymentRequest struct {
	MerchantID string  `json:"merchant_id"  binding:"required"`
	Amount     float64 `json:"amount"  binding:"required"`
}
