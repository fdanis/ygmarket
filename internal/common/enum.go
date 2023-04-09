package common

import (
	"encoding/json"
	"fmt"
	"strings"
)

type OrderStatus uint8

const (
	New        OrderStatus = 0
	Processing OrderStatus = 1
	Invalid    OrderStatus = 2
	Processed  OrderStatus = 3
)

var (
	StatusName = map[uint8]string{
		0: "NEW",
		1: "PROCESSING",
		2: "INVALID",
		3: "PROCESSED",
	}
	StatusValue = map[string]uint8{
		"NEW":        0,
		"PROCESSING": 1,
		"INVALID":    2,
		"PROCESSED":  3,
	}
)

func (o *OrderStatus) UnmarshalJSON(data []byte) error {
	var tmp string
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	value, ok := StatusValue[strings.TrimSpace(tmp)]
	if !ok {
		return fmt.Errorf("%q is not a valid status", tmp)
	}
	*o = OrderStatus(value)
	return nil
}

func (o OrderStatus) MarshalJSON() ([]byte, error) {
	if v, ok := StatusName[uint8(o)]; ok {
		return json.Marshal(v)
	}
	return json.Marshal(0)
}
