package model

type Customer struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Balance  int    `json:"balance"`
}
