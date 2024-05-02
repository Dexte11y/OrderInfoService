package repository

import (
	"fmt"
	"orderinfoservice/internal/model"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (repo *Repository) CreateOrderRepository(order model.Order) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// paymentDt := time.Unix(order.Payment.PaymentDt, 0)

	queryOrder := `
	INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err = tx.Exec(queryOrder, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return fmt.Errorf("ошибка при создании заказа: %w", err)
	}

	queryDelivery := `
	INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = tx.Exec(queryDelivery, order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении информации о доставке: %w", err)
	}

	queryPayment := `
	INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err = tx.Exec(queryPayment, order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении информации о платеже: %w", err)
	}

	queryItems := `
	INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	stmt, err := tx.Prepare(queryItems)
	if err != nil {
		return fmt.Errorf("ошибка при подготовке запроса для товаров: %w", err)
	}
	defer stmt.Close()

	for _, item := range order.Items {
		_, err = stmt.Exec(order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return fmt.Errorf("ошибка при добавлении товара: %w", err)
		}
	}
	return nil
}

func (repo *Repository) GetLast10Orders() ([]model.Order, error) {
	// Запрос к базе данных для получения последних 10 заказов
	query := `
        SELECT * FROM orders
        ORDER BY date_created DESC
        LIMIT 10
    `
	var orders []model.Order
	if err := repo.db.Select(&orders, query); err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}

	// Для каждого заказа получаем информацию о доставке, платеже и товарах
	for i := range orders {
		order := &orders[i]
		if err := repo.fillOrderDetails(order); err != nil {
			return nil, fmt.Errorf("ошибка при получении деталей заказа: %w", err)
		}
	}

	return orders, nil
}

// Заполнение деталей заказа (информации о доставке, платеже и товарах)
func (repo *Repository) fillOrderDetails(order *model.Order) error {
	// Получение информации о доставке
	if err := repo.fillDeliveryInfo(order); err != nil {
		return err
	}

	// Получение информации о платеже
	if err := repo.fillPaymentInfo(order); err != nil {
		return err
	}

	// Получение информации о товарах
	if err := repo.fillItemsInfo(order); err != nil {
		return err
	}

	return nil
}

// Заполнение информации о доставке для заказа
func (repo *Repository) fillDeliveryInfo(order *model.Order) error {
	query := `
        SELECT * FROM delivery
        WHERE order_uid = $1
    `
	if err := repo.db.Get(&order.Delivery, query, order.OrderUID); err != nil {
		return fmt.Errorf("ошибка при получении информации о доставке: %w", err)
	}
	return nil
}

// Заполнение информации о платеже для заказа
func (repo *Repository) fillPaymentInfo(order *model.Order) error {
	query := `
        SELECT * FROM payment
        WHERE order_uid = $1
    `
	if err := repo.db.Get(&order.Payment, query, order.OrderUID); err != nil {
		return fmt.Errorf("ошибка при получении информации о платеже: %w", err)
	}
	return nil
}

// Заполнение информации о товарах для заказа
func (repo *Repository) fillItemsInfo(order *model.Order) error {
	query := `
        SELECT * FROM items
        WHERE order_uid = $1
    `
	var items []model.Item
	if err := repo.db.Select(&items, query, order.OrderUID); err != nil {
		return fmt.Errorf("ошибка при получении информации о товарах: %w", err)
	}
	order.Items = items
	return nil
}
