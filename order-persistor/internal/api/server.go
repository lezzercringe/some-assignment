package api

import (
	"log/slog"
	"net"
	"net/http"
	"order-persistor/internal/config"
	"order-persistor/internal/orders"

	gorilla "github.com/gorilla/handlers"
)

type Params struct {
	Logger           *slog.Logger
	OrdersRepository orders.Repository
}

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

	return &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}
}
