package handlers

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/fdanis/yg-loyalsys/internal/app"
	"github.com/fdanis/yg-loyalsys/internal/db/entities"
	"github.com/fdanis/yg-loyalsys/internal/db/repositories"
	"github.com/fdanis/yg-loyalsys/internal/helpers"
	"github.com/fdanis/yg-loyalsys/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type OrderHandler struct {
	orderRepository repositories.OrderRepository
}

func NewOrderHandler(app *app.App, orderRepository repositories.OrderRepository) OrderHandler {
	result := OrderHandler{
		orderRepository: orderRepository,
	}
	return result
}

func (h *OrderHandler) NewOrder(w http.ResponseWriter, r *http.Request) {
	number, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if !helpers.CheckNumber(string(number)) {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	userid := getUserID(r)

	err = h.orderRepository.Add(entities.Order{UserID: userid, Number: string(number), Status: 0})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				{
					order, err := h.orderRepository.GetByNumber(string(number))
					if err != nil || order == nil {
						log.Println(err)
						http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
						return
					}

					if order.UserID != userid {
						http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
						w.Write([]byte("this number has already been registered"))
						return
					} else {
						return
					}
				}
			default:
			}
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			w.Write([]byte("server error"))
			return
		}
	}
	//todo отправить запрос  в черный ящик
	w.WriteHeader(http.StatusAccepted)
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	userid := getUserID(r)
	orders, err := h.orderRepository.GetAllByUser(userid)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		w.Write([]byte("server error"))
		return
	}

	if len(orders) > 0 {
		model := make([]models.Order, 0, len(orders))
		for _, o := range orders {
			model = append(model, models.Order{
				Number:     o.Number,
				Status:     o.Status,
				UploadedAt: o.Created,
				Accrual:    o.Accrual,
			})
		}
		responseJSON(w, model)
	}
}
