package api

import (
	"log/slog"
	"net"
	"net/http"
	"order-persistor/internal/config"
	"order-persistor/internal/orders"

	_ "order-persistor/docs"

	gorilla "github.com/gorilla/handlers"
	swagger "github.com/swaggo/http-swagger"
)

type Params struct {
	Logger           *slog.Logger
	OrdersRepository orders.Repository
}

// @title           Order-persistor API
// @version         1.0

func NewServer(cfg config.API, p Params) *http.Server {
	mux := http.NewServeMux()
	handler := GetOrderHandler{
		Logger:     p.Logger,
		Repository: p.OrdersRepository,
	}

	httpAddr := net.JoinHostPort(cfg.Host, cfg.Port)
	mux.Handle("/order/{id}", stackMiddleware(
		&handler,
		gorilla.RecoveryHandler(),
		gorilla.CORS(),
		NewLogMiddleware(p.Logger),
	))
	mux.Handle("/swagger/", swagger.WrapHandler)

	return &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}
}
