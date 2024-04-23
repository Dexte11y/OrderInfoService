package handler

import (
	"net/http"
	service "orderinfoservice/internal/service"

	"github.com/nats-io/stan.go"
)

type Handler struct {
	Service *service.Service // Зависимость
}

func NewHandler(Service *service.Service) *Handler {
	// Создание экземпляра обработчика и его инициализация
	return &Handler{
		Service: Service,
	}
}

// InitRoutes инициализирует маршруты
func (h *Handler) InitRoutes(sc stan.Conn) {
	http.HandleFunc("/create_order", func(w http.ResponseWriter, r *http.Request) {
		h.CreateOrderHandler(w, r, sc)
	})

	http.HandleFunc("/about", aboutHandler)
	// Добавьте другие маршруты здесь по мере необходимости
}

// Обработчик для URL /about
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Это страница 'О нас'."))
}
