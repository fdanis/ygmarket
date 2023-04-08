package handlers

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/fdanis/yg-loyalsys/internal/db/entities"
	"github.com/fdanis/yg-loyalsys/internal/db/repositories"
	"github.com/fdanis/yg-loyalsys/internal/helpers"
	"github.com/fdanis/yg-loyalsys/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type WithdrawHandler struct {
	userRepository     repositories.UserRepository
	withdrawRepository repositories.WithdrawRepository
}

func NewWithdrawHandler(userRepository repositories.UserRepository, withdrawRepository repositories.WithdrawRepository) WithdrawHandler {
	result := WithdrawHandler{
		userRepository:     userRepository,
		withdrawRepository: withdrawRepository,
	}
	return result
}

func (h *WithdrawHandler) NewWithdraw(w http.ResponseWriter, r *http.Request) {
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

	err = h.withdrawRepository.Add(entities.Withdraw{UserID: userid, Number: string(number), Sum: 10})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				{
					http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
					w.Write([]byte("this number has already been registered"))
					return
				}
			case pgerrcode.CheckViolation:
				{
					http.Error(w, http.StatusText(http.StatusPaymentRequired), http.StatusPaymentRequired)
					w.Write([]byte("not enough money"))
					return
				}
			default:
			}
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			w.Write([]byte("server error"))
			return
		}
	}
}

func (h *WithdrawHandler) GetWithdraw(w http.ResponseWriter, r *http.Request) {
	userid := getUserID(r)
	withdraw, err := h.withdrawRepository.GetAllByUser(userid)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		w.Write([]byte("server error"))
		return
	}

	if len(withdraw) > 0 {
		model := make([]models.Withdraw, 0, len(withdraw))
		for _, o := range withdraw {
			model = append(model, models.Withdraw{
				Number:      o.Number,
				Sum:         o.Sum,
				ProcessedAt: o.Created,
			})
		}
		responseJSON(w, model)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
