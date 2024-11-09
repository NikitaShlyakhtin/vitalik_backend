package types

import (
	"github.com/guregu/null/v5"
	"time"
)

type Transaction struct {
	ID string `json:"id,omitempty"`

	SenderRequisites   Requisites
	ReceiverRequisites Requisites

	Amount   float64  `json:"amount,omitempty"`
	Currency Currency `json:"currency,omitempty"`

	Purpose null.String `json:"purpose,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
