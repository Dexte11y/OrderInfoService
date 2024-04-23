package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func ConnectPostgres() (*sqlx.DB, error) {
	// Замените эти данные на ваши параметры подключения
	connStr := "user=postgres dbname=order password=postgres host=localhost port=5432 sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии соединения с базой данных: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ошибка при проверке соединения с базой данных: %v", err)
	}

	return db, nil
}
