package types

import "time"

type Wallet struct {
	Requisites Requisites

	Currency Currency `json:"currency"`
	Balance  float64  `json:"balance"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Requisites struct {
	UserID  string `json:"user_id"`
	Address string `json:"address"`
}
