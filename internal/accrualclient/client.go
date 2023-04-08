package accrualclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/fdanis/yg-loyalsys/internal/db/repositories"
	"github.com/fdanis/yg-loyalsys/internal/models"
)

type Client struct {
	address         string
	orderRepository repositories.OrderRepository
}

func NewClient(address string, orderRepository repositories.OrderRepository) (*Client, error) {
	return &Client{address: address, orderRepository: orderRepository}, nil
}

func (c *Client) Run(ctx context.Context) {
	orders, err := c.orderRepository.GetAllForChecking()
	if err != nil {
		log.Printf("get order for checking %v\n", err)
	}
	fmt.Println(len(orders))
	for _, v := range orders {
		select {
		case <-ctx.Done():
		default:
			{
				m, err := c.Send(v.Number)
				if err != nil {
					log.Printf("send order number %s : %v\n", v.Number, err)
					continue
				}
				if m == nil {
					continue
				}
				if v.Status != m.Status {
					v.Status = m.Status
					v.Accrual = m.Accrual
					c.orderRepository.Update(v)
				}
			}
		}
	}
}

func (c *Client) Send(number string) (*models.Order, error) {
	res, err := http.Get(fmt.Sprintf("%s/api/orders/%s", c.address, number))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	switch res.StatusCode {
	case http.StatusOK:
		{
			var order models.Order
			defer res.Body.Close()
			err := json.NewDecoder(res.Body).Decode(&order)
			if err != nil {
				return nil, err
			}
			return &order, nil
		}
	case http.StatusTooManyRequests:
		{
			val, ok := res.Header["Retry-After"]
			if ok && len(val) > 0 {
				s, err := strconv.ParseInt(val[0], 10, 64)
				if err != nil {
					log.Println(err)
				}
				time.Sleep(time.Second * time.Duration(s))
				return c.Send(number)
			}
		}
	case http.StatusNoContent:
		{
			//do nothing
			log.Println("nocontent")
		}
	default:
		{
			return nil, fmt.Errorf("incorect status %s", res.Status)
		}
	}
	return nil, nil
}
