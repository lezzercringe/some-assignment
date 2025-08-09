package api

import (
	"errors"
	"log/slog"
	"net/http"
	"order-persistor/internal/orders"
)

type GetOrderHandler struct {
	Logger     *slog.Logger
	Repository orders.Repository
}

// GetOrder godoc
// @Summary Get order by ID
// @Description Returns the order object for the specified ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} orders.Order
// @Failure 400 {object} Error "Invalid request"
// @Failure 404 {object} Error "Order not found"
// @Failure 500 {object} Error "Internal server error"
// @Router /order/{id} [get]
func (h *GetOrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := h.Logger.With("url", r.URL)

	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	orderID := r.PathValue("id")
	if orderID == "" {
		log.ErrorContext(r.Context(), "handler was called with empty path parameter value")
		responseInternalError.Write(w)
		return
	}

	order, err := h.Repository.GetByID(r.Context(), orderID)
	if err != nil {
		if errors.Is(err, orders.ErrNotFound) {
			newErrorResponse(404, "order with that id was not found").Write(w)
			return
		}

		log.ErrorContext(r.Context(), "retrieving order from repository", "err", err)
		responseInternalError.Write(w)
		return
	}

	if err := respondJSON(order, w); err != nil {
		log.ErrorContext(r.Context(), "sending http response", "err", err)
		responseInternalError.Write(w)
	}
}
