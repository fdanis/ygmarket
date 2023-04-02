package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/fdanis/yg-loyalsys/internal/app"
	"github.com/fdanis/yg-loyalsys/internal/common"
	"github.com/fdanis/yg-loyalsys/internal/db/entities"
	"github.com/fdanis/yg-loyalsys/internal/db/repositories"
	"github.com/fdanis/yg-loyalsys/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserHandler struct {
	userRepository repositories.UserRepository
	secretKey      string
}

func NewUserHandler(app *app.App, userRepository repositories.UserRepository) UserHandler {
	result := UserHandler{
		userRepository: userRepository,
		secretKey:      app.Config.SecretKey,
	}

	return result
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if !validateContentTypeIsJSON(w, r) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var model models.User
	if err := decodeJSONBody(r.Body, r.Header.Get("Content-Encoding"), &model); err != nil {
		var mr *RequestError
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	user := &entities.User{Login: model.Login, Password: model.Password}
	err := h.userRepository.Add(user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				http.Error(w, "логин уже занят", http.StatusConflict)
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}
	}
	h.addAuthHeader(w, user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if !validateContentTypeIsJSON(w, r) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var model models.User
	if err := decodeJSONBody(r.Body, r.Header.Get("Content-Encoding"), &model); err != nil {
		var mr *RequestError
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	user, err := h.userRepository.GetByLogin(model.Login)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if user == nil || user.Password != model.Password {
		http.Error(w, "неверная пара логин/пароль", http.StatusUnauthorized)
		return
	}

	h.addAuthHeader(w, user)
}

func (h *UserHandler) addAuthHeader(w http.ResponseWriter, user *entities.User) {
	expTime := time.Now().Add(30 * time.Minute)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, common.UserClaims{Login: user.Login, ID: user.ID, RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(expTime)}})
	signedToken, err := token.SignedString([]byte(h.secretKey))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Authorization", signedToken)
}

func (h *UserHandler) Balance(w http.ResponseWriter, r *http.Request) {
	userid := getUserID(r)
	user, err := h.userRepository.GetByID(userid)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	responseJSON(w, struct {
		current   float32
		withdrawn float32
	}{current: user.Balance, withdrawn: user.Withdrawn})
}
