package entities

import (
	"time"
)

type Withdraw struct {
	ID      int
	UserID  int
	Order   string
	Sum     float32
	Created time.Time
}
