package nats

import (
	"encoding/json"
	"fmt"
	"orderinfoservice/internal/model"

	"github.com/nats-io/stan.go"
)

func PublishOrder(sc stan.Conn, order model.Order) error {
	// Преобразование данных о заказе в JSON

	orderJSON, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("ошибка при маршалинге заказа в JSON: %w", err)
	}

	// Публикация JSON-данных о заказе в канал NATS
	key := order.OrderUID
	message := append([]byte(key), orderJSON...)

	// Публикация сообщения в канал NATS
	if err := sc.Publish("example-subject", message); err != nil {
		return fmt.Errorf("ошибка при публикации заказа в NATS: %w", err)
	}

	return nil
}
