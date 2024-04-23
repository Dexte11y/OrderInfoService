package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	model "orderinfoservice/internal/model"

	"github.com/nats-io/stan.go"
)

func (h *Handler) CreateOrderHandler(w http.ResponseWriter, r *http.Request, sc stan.Conn) {
	// Парсинг данных из JSON-запроса
	var order model.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Service.CreteOrderService(sc, order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Отправка ответа об успешном создании заказа
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Заказ успешно создан")
}
