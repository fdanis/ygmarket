package models

import (
	"encoding/json"
	"fmt"
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

	var accrual string
	if o.Status == common.PROCESSED || o.Accrual > 0 {
		accrual = fmt.Sprintf("%.3f", o.Accrual)
	}
	return json.Marshal(&struct {
		*Alias
		Accrual  string `json:"accrual,omitempty"`
		LastSeen string `json:"uploaded_at"`
	}{
		LastSeen: o.UploadedAt.Format(time.RFC3339),
		Accrual:  accrual,
		Alias:    (*Alias)(o),
	})
}

type AccrualOrder struct {
	Number  string             `json:"order"`
	Status  common.OrderStatus `json:"status"`
	Accrual float32            `json:"accrual,omitempty"`
}
