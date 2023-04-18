package accrualclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	for _, v := range orders {
		select {
		case <-ctx.Done():
		default:
			{
				m, err := c.Send(ctx, v.Number)
				if err != nil {
					log.Printf("send order number %s : %v\n", v.Number, err)
					continue
				}
				if m == nil {
					continue
				}
				if v.Status == m.Status {
					continue
				}
				v.Status = m.Status
				v.Accrual = m.Accrual
				err = c.orderRepository.Update(v)
				if err != nil {
					log.Printf("write to db  %s : %v\n", v.Number, err)
					continue
				}
			}
		}
	}
}

func (c *Client) Send(ctx context.Context, number string) (*models.AccrualOrder, error) {
	res, err := http.Get(fmt.Sprintf("%s/api/orders/%s", c.address, number))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case http.StatusOK:
		var order models.AccrualOrder
		err := json.NewDecoder(res.Body).Decode(&order)
		if err != nil {
			return nil, err
		}
		return &order, nil
	case http.StatusTooManyRequests:
		val, ok := res.Header["Retry-After"]
		if ok && len(val) > 0 {
			s, err := strconv.ParseInt(val[0], 10, 64)
			if err != nil {
				log.Println(err)
			}
			t := time.NewTimer(time.Second * time.Duration(s))
			select {
			case <-ctx.Done():
				return nil, nil
			case <-t.C:
				return c.Send(ctx, number)
			}
		}
	case http.StatusNoContent:
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		log.Println(string(b))
	default:
		return nil, fmt.Errorf("incorect status %s", res.Status)
	}
	return nil, nil
}
