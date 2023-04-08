package common

import (
	"encoding/json"
	"fmt"
	"strings"
)

type OrderStatus uint8

const (
	NEW        OrderStatus = 0
	PROCESSING OrderStatus = 1
	INVALID    OrderStatus = 2
	PROCESSED  OrderStatus = 3
)

var (
	Status_name = map[uint8]string{
		0: "NEW",
		1: "PROCESSING",
		2: "INVALID",
		3: "PROCESSED",
	}
	Status_value = map[string]uint8{
		"NEW":        0,
		"PROCESSING": 1,
		"INVALID":    2,
		"PROCESSED":  3,
	}
)

func (o *OrderStatus) UnmarshaJSON(data []byte) error {
	var tmp string
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	value, ok := Status_value[strings.TrimSpace(tmp)]
	if !ok {
		return fmt.Errorf("%q is not a valid status", tmp)
	}
	*o = OrderStatus(value)
	return nil
}

func (s *OrderStatus) MarshalJSON() ([]byte, error) {
	if v, ok := Status_name[uint8(*s)]; ok {
		return json.Marshal(v)
	}
	return json.Marshal(0)
}
