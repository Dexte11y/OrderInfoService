// cmd/main.go
package main

import (
	"log"
	"log/slog"
	"net/http"
	cache "orderinfoservice/internal/cache"
	database "orderinfoservice/internal/database"
	subscriber "orderinfoservice/internal/nats"
	repository "orderinfoservice/internal/repository"

	handler "orderinfoservice/internal/handler"
	migrate "orderinfoservice/internal/migrate"
	nats "orderinfoservice/internal/nats"
	service "orderinfoservice/internal/service"

	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

func main() {
	logger := slog.Default()

	//TODO добавить clean_env

	sc, err := nats.ConnectToNats()
	if err != nil {
		logger.Error("Ошибка подключения к NATS: ", err)
		return
	}
	defer sc.Close()

	logger.Info("Успешное подключение к NATS")

	sub, err := subscriber.NewSubscribe(sc)
	if err != nil {
		logger.Error("Ошибка создания подписчика: ", err)
		return
	}
	defer sub.Close()

	logger.Info("Подписчик запущен. Ожидание сообщений...")

	// Создание подключения к базе данных
	db, err := database.ConnectPostgres()
	if err != nil {
		logger.Error("Ошибка подключения к DB: ", err)
		return
	}
	defer db.Close()

	logger.Info("Успешное подключение к DB")

	// Инициализация драйвера миграций
	dbx := db.DB
	driver, err := postgres.WithInstance(dbx, &postgres.Config{})
	if err != nil {
		logger.Error("Ошибка при создании экземпляра драйвера PostgreSQL:", err)
		return
	}

	// Вызов логики миграций
	if err := migrate.ApplyMigrations(dbx, driver); err != nil {
		logger.Error("Ошибка применения миграции:", err)
		return
	}

	logger.Info("Миграции применины успешно")

	// Инициализация маршрутов
	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	handler := handler.NewHandler(service)
	if err := cache.InitCache(repository); err != nil {
		log.Fatalf("ошибка при инициализации кэша: %v", err)
	}

	logger.Info("Инициализация кэша выполнена")
	handler.InitRoutes(sc)

	// Запуск веб-сервера на порте 8080
	logger.Info("Сервер запущен на порте :8080")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
