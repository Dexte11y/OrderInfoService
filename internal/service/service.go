package service

import (
	"fmt"
	"orderinfoservice/internal/cache"
	"orderinfoservice/internal/model"
	"orderinfoservice/internal/nats"
	"orderinfoservice/internal/repository"

	"github.com/nats-io/stan.go"
)

type Service struct {
	Repository *repository.Repository // Зависимость
}

func NewService(Repository *repository.Repository) *Service {
	// Создание экземпляра сервиса и его инициализация
	return &Service{
		Repository: Repository,
	}
}

func (s *Service) CreteOrderService(sc stan.Conn, order model.Order) error {

	// Публикация данных о заказе
	err := nats.PublishOrder(sc, order)
	if err != nil {
		return fmt.Errorf("ошибка публикации заказа: %w", err)
	}

	err = s.Repository.CreateOrderRepository(order)
	if err != nil {
		return fmt.Errorf("ошибка создания заказа в БД: %w", err)
	}

	err = cache.SetOrderToCache(order)
	if err != nil {
		return fmt.Errorf("ошибка добавления заказа в кэш: %w", err)
	}

	return nil
}

func (s *Service) GetOrderByIDService(sc stan.Conn, orderUID string) (interface{}, error) {
	val, err := cache.GetOrderByIDCache(orderUID)
	return val, err
}
