package models

import (
	"encoding/json"
	"time"
)

type Withdraw struct {
	Number      string    `json:"number"`
	Sum         float32   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

func (o *Withdraw) MarshalJSON() ([]byte, error) {
	type Alias Withdraw

	return json.Marshal(&struct {
		*Alias
		ProcessedAt string `json:"processed_at"`
	}{
		ProcessedAt: o.ProcessedAt.Format(time.RFC3339),
		Alias:       (*Alias)(o),
	})
}
