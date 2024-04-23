package service

import (
	"fmt"
	"orderinfoservice/internal/model"
	"orderinfoservice/internal/publisher"
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
	err := publisher.PublishOrder(sc, order)
	if err != nil {
		return fmt.Errorf("ошибка публикации заказа: %w", err)
	}

	err = s.Repository.CreateOrderRepository(order)
	if err != nil {
		return fmt.Errorf("ошибка создания заказа в БД: %w", err)
	}

	return nil
}
