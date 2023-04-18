package routes

import (
	"net/http"

	"github.com/fdanis/yg-loyalsys/internal/app"
	"github.com/fdanis/yg-loyalsys/internal/handlers"
	"github.com/fdanis/yg-loyalsys/internal/middleware"
	"github.com/go-chi/chi"
)

func Routes(app *app.App) http.Handler {
	userHandler := handlers.NewUserHandler(app.Config.SecretKey, app.UserRepository)
	orderHandler := handlers.NewOrderHandler(app.OrderRepository)
	withdrawHandler := handlers.NewWithdrawHandler(app.UserRepository, app.WithdrawRepository)
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
