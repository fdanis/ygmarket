package models

import (
	"encoding/json"
	"time"

	"github.com/fdanis/yg-loyalsys/internal/common"
)

type Order struct {
	Number     string             `json:"number"`
	Status     common.OrderStatus `json:"status"`
	Accrual    float32            `json:"accrual,omitempty"`
	UploadedAt time.Time          `json:"uploaded_at"`
}

func (o *Order) MarshalJSON() ([]byte, error) {
	type Alias Order

	return json.Marshal(&struct {
		*Alias
		LastSeen string `json:"uploaded_at"`
	}{
		LastSeen: o.UploadedAt.Format(time.RFC3339),
		Alias:    (*Alias)(o),
	})
}

type AccrualOrder struct {
	Number  string             `json:"order"`
	Status  common.OrderStatus `json:"status"`
	Accrual float32            `json:"accrual,omitempty"`
}
