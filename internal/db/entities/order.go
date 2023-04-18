package entities

import (
	"time"

	"github.com/fdanis/yg-loyalsys/internal/common"
)

type Order struct {
	ID      int
	UserID  int
	Number  string
	Accrual float32
	Status  common.OrderStatus
	Created time.Time
}
