package cache

import (
	"fmt"
	"log"
	"orderinfoservice/internal/model"
	"orderinfoservice/internal/repository"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	c      = cache.New(5*time.Minute, 10*time.Minute)
	cMutex sync.RWMutex
)

func InitCache(repo *repository.Repository) error {
	// Получаем последние 10 заказов из репозитория
	orders, err := repo.GetLast10Orders()
	if err != nil {
		return fmt.Errorf("ошибка при получении последних 10 заказов: %w", err)
	}

	// Кэшируем каждый заказ
	for _, order := range orders {
		key := order.OrderUID
		if _, found := c.Get(key); found {
			continue // Пропускаем кэширование, если ключ уже существует
		}

		// Кэшируем заказ
		cMutex.Lock()
		c.Set(key, order, cache.DefaultExpiration)
		cMutex.Unlock()
	}

	log.Println("Инициализация кэша выполнена")
	return nil
}

func SetOrderToCache(order model.Order) error {
	key := order.OrderUID
	if _, found := c.Get(key); found {
		return nil // Ключ уже существует в кэше, поэтому возвращаем nil
	}

	cMutex.Lock()
	defer cMutex.Unlock()

	if err := c.Add(key, order, cache.DefaultExpiration); err != nil {
		return fmt.Errorf("ошибка добавления заказа в кэш: %w", err) // Возвращаем ошибку, если есть проблемы с кэшированием
	}
	return nil
}

func GetOrderByIDCache(orderUID string) (interface{}, error) {
	if value, found := c.Get(orderUID); found {
		log.Println("Значение из кеша:", value)
		return value, nil
	} else {
		return nil, fmt.Errorf("запись не найдена в кеше")
	}
}
