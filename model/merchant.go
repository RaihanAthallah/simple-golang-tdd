package model

type Merchant struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	BankAccount string  `json:"bank_account"`
	BankName    string  `json:"bank_name"`
	Balance     float64 `json:"balance"`
}
