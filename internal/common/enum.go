package common

type OrderStatus uint8

const (
	NEW        OrderStatus = 0
	PROCESSING OrderStatus = 1
	INVALID    OrderStatus = 2
	PROCESSED  OrderStatus = 3
)
