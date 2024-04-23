package repository

import (
	"fmt"
	"orderinfoservice/internal/model"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB // Пример зависимости, например, база данных
}

func NewRepository(db *sqlx.DB) *Repository {
	// Создание экземпляра репозитория и его инициализация
	return &Repository{
		db: db,
	}
}

func (repo *Repository) CreateOrderRepository(order model.Order) error {
	_, err := repo.db.Exec("INSERT INTO testorder (item, amount) VALUES ($1, $2)", order.Item, order.Amount)
	if err != nil {
		panic(err)
	}
	fmt.Println("Заказ успешно добавлен в БД")
	return nil
}
