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

	var status string
	var accrual string

	switch o.Status {
	case common.NEW:
		{
			status = "NEW"
		}
	case common.PROCESSING:
		{
			status = "PROCESSING"
		}
	case common.PROCESSED:
		{
			status = "PROCESSED"
			accrual = fmt.Sprintf("%.3f", o.Accrual)
		}
	case common.INVALID:
		{
			status = "INVALID"
		}
	default:
		status = "NEW"
	}

	return json.Marshal(&struct {
		*Alias
		Status   string `json:"status"`
		Accrual  string `json:"accrual,omitempty"`
		LastSeen string `json:"uploaded_at"`
	}{
		LastSeen: o.UploadedAt.Format(time.RFC3339),
		Accrual:  accrual,
		Status:   status,
		Alias:    (*Alias)(o),
	})
}
