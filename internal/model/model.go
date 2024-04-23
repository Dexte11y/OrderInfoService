package model

type Order struct {
	ID     int    `json:"id" db:"id"`
	Item   string `json:"item" db:"item"`
	Amount int    `json:"amount" db:"amount"`
}
