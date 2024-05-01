package migrate

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database"
)

// ApplyMigrations применяет миграции к базе данных
func ApplyMigrations(db *sql.DB, driver database.Driver) error {
	// Создание экземпляра объекта миграций
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("ошибка создания экземпляра миграции: %w", err)
	}

	// Применение миграций
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}
		return fmt.Errorf("ошибка применения миграции: %w", err)
	}

	return nil
}
