package types

import (
	"github.com/guregu/null/v5"
	"time"
)

type Transaction struct {
	ID string `json:"id,omitempty"`

	SenderRequisites   Requisites `json:"sender_requisites,omitempty"`
	ReceiverRequisites Requisites `json:"receiver_requisites,omitempty"`

	Amount   float64  `json:"amount,omitempty"`
	Currency Currency `json:"currency,omitempty"`

	Purpose null.String `json:"purpose,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
