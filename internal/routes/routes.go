package routes

import (
	"database/sql"
	"net/http"

	"github.com/fdanis/yg-loyalsys/internal/app"
	"github.com/fdanis/yg-loyalsys/internal/db/repositories"
	"github.com/fdanis/yg-loyalsys/internal/handlers"
	"github.com/fdanis/yg-loyalsys/internal/middleware"
	"github.com/go-chi/chi"
)

func Routes(app *app.App, db *sql.DB) http.Handler {
	userHandler := handlers.NewUserHandler(app, repositories.NewUserRepository(db))
	orderHandler := handlers.NewOrderHandler(app, repositories.NewOrderRepository(db))
	withdrawHandler := handlers.NewWithdrawHandler(app, repositories.NewUserRepository(db), repositories.NewWithdrawRepository(db))
	mux := chi.NewRouter()
	mux.Use(middleware.GzipHandle)
	mux.Post("/api/user/login", userHandler.Login)
	mux.Post("/api/user/register", userHandler.Register)

	mux.Route("/api/user", func(r chi.Router) {
		r.Use(middleware.NewAuthorizeMiddleware(app.Config.SecretKey).Authorize)
		r.Post("/orders", orderHandler.NewOrder)
		r.Get("/orders", orderHandler.GetOrders)
		r.Get("/balance", userHandler.Balance)
		r.Post("/balance/withdraw", withdrawHandler.NewWithdraw)
		r.Get("/withdrawals", withdrawHandler.GetWithdraw)
	})
	return mux
}
