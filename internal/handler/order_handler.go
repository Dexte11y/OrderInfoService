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

func (h *Handler) GetOrderByIDHandler(w http.ResponseWriter, r *http.Request, sc stan.Conn) {
	orderUID := r.URL.Query().Get("orderUID")
	order, err := h.Service.GetOrderByIDService(sc, orderUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	orderJSON, err := json.Marshal(order)
	if err != nil {
		// Обработка ошибки сериализации JSON
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Отправляем данные в виде JSON в http.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(orderJSON)

}
