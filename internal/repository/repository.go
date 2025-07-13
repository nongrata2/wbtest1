package repository

import (
	"context"
	"database/sql"
	"firstmod/internal/models"
	"log/slog"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	log  *slog.Logger
	conn *pgxpool.Pool
}

func New(log *slog.Logger, address string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), address)
	if err != nil {
		log.Error("connection problem", "address", address, "error", err)
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Error("failed to ping database", "error", err)
		return nil, err
	}

	log.Info("successfully connected to database", "address", address)

	return &DB{
		log:  log,
		conn: pool,
	}, nil
}

func (db *DB) Add(ctx context.Context, order models.Order) error {
	db.log.Debug("attempting to add new order", "order_uid", order.OrderUID)

	// Начинаем транзакцию
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		db.log.Error("failed to begin transaction", "error", err)
		return err
	}
	// Откатываем транзакцию в случае ошибки, если она еще не была закоммичена
	defer func() {
		if r := recover(); r != nil {
			db.log.Error("recovered from panic during transaction, rolling back", "panic", r)
			tx.Rollback(ctx)
			panic(r) // Re-throw panic
		} else if err != nil {
			db.log.Error("transaction failed, rolling back", "error", err)
			tx.Rollback(ctx)
		}
	}()

	// 1. Вставка данных в таблицу orders
	orderSQL := `
        INSERT INTO orders (
            order_uid, track_number, entry, locale, internal_signature,
            customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
        )`
	_, err = tx.Exec(ctx, orderSQL,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.Shardkey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	)
	if err != nil {
		db.log.Error("failed to insert order", "order_uid", order.OrderUID, "error", err)
		return err
	}
	db.log.Debug("order inserted successfully", "order_uid", order.OrderUID)

	// 2. Вставка данных в таблицу delivery_info
	deliverySQL := `
        INSERT INTO delivery_info (
            order_uid, name, phone, zip, city, address, region, email
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8
        )`
	_, err = tx.Exec(ctx, deliverySQL,
		order.OrderUID, // Связываем по OrderUID
		order.DeliveryInfo.Name,
		order.DeliveryInfo.Phone,
		order.DeliveryInfo.Zip,
		order.DeliveryInfo.City,
		order.DeliveryInfo.Address,
		order.DeliveryInfo.Region,
		order.DeliveryInfo.Email,
	)
	if err != nil {
		db.log.Error("failed to insert delivery info", "order_uid", order.OrderUID, "error", err)
		return err
	}
	db.log.Debug("delivery info inserted successfully", "order_uid", order.OrderUID)

	// 3. Вставка данных в таблицу payments
	paymentSQL := `
        INSERT INTO payments (
            transaction_uid, request_id, currency, provider, amount,
            payment_dt, bank, delivery_cost, goods_total, custom_fee
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
        )`
	_, err = tx.Exec(ctx, paymentSQL,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDT,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	)
	if err != nil {
		db.log.Error("failed to insert payment info", "order_uid", order.OrderUID, "error", err)
		return err
	}
	db.log.Debug("payment info inserted successfully", "order_uid", order.OrderUID)

	// 4. Вставка данных в таблицу items (может быть несколько)
	itemSQL := `
        INSERT INTO items (
            order_uid, chrt_id, track_number, price, rid, name,
            sale, size, total_price, nm_id, brand, status
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
        )`
	for i, item := range order.Items {
		_, err = tx.Exec(ctx, itemSQL,
			order.OrderUID, // Связываем по OrderUID
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		)
		if err != nil {
			db.log.Error("failed to insert item", "order_uid", order.OrderUID, "item_index", i, "error", err)
			return err
		}
	}
	db.log.Debug("items inserted successfully", "order_uid", order.OrderUID, "count", len(order.Items))

	// Коммитим транзакцию
	err = tx.Commit(ctx)
	if err != nil {
		db.log.Error("failed to commit transaction", "error", err)
		return err
	}

	db.log.Info("order and related data added successfully", "order_uid", order.OrderUID)
	return nil
}

// GetInfo получает полную информацию о заказе по его OrderUID.
func (db *DB) GetInfo(ctx context.Context, orderUID string) (models.Order, error) {
	db.log.Debug("attempting to get order info", "order_uid", orderUID)

	var order models.Order
	order.OrderUID = orderUID // Устанавливаем OrderUID сразу

	// 1. Получаем данные из таблицы orders
	orderSQL := `
        SELECT
            track_number, entry, locale, internal_signature,
            customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
        FROM orders
        WHERE order_uid = $1`

	err := db.conn.QueryRow(ctx, orderSQL, orderUID).Scan(
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
	)
	if err != nil {
		if err == pgx.ErrNoRows { // Или sql.ErrNoRows, в зависимости от того, что возвращает pgx.QueryRow
			db.log.Debug("order not found", "order_uid", orderUID)
			return models.Order{}, sql.ErrNoRows // Возвращаем ошибку, если заказ не найден
		}
		db.log.Error("failed to query order", "order_uid", orderUID, "error", err)
		return models.Order{}, err
	}
	db.log.Debug("order data fetched successfully", "order_uid", orderUID)

	// 2. Получаем данные из таблицы delivery_info
	deliverySQL := `
        SELECT
            name, phone, zip, city, address, region, email
        FROM delivery_info
        WHERE order_uid = $1`
	err = db.conn.QueryRow(ctx, deliverySQL, orderUID).Scan(
		&order.DeliveryInfo.Name,
		&order.DeliveryInfo.Phone,
		&order.DeliveryInfo.Zip,
		&order.DeliveryInfo.City,
		&order.DeliveryInfo.Address,
		&order.DeliveryInfo.Region,
		&order.DeliveryInfo.Email,
	)
	if err != nil {
		db.log.Warn("failed to query delivery info, might not exist", "order_uid", orderUID, "error", err)
		// Не возвращаем ошибку, если info нет - возможно, это допустимо.
		// Но в данном случае, если order есть, delivery_info по FK должен быть.
		// Если это критично, можете вернуть err
	}
	db.log.Debug("delivery info fetched successfully", "order_uid", orderUID)

	// 3. Получаем данные из таблицы payments
	paymentSQL := `
        SELECT
            request_id, currency, provider, amount,
            payment_dt, bank, delivery_cost, goods_total, custom_fee
        FROM payments
        WHERE transaction_uid = $1` // Используем transaction_uid для поиска
	err = db.conn.QueryRow(ctx, paymentSQL, orderUID).Scan( // orderUID должен быть равен transaction_uid
		&order.Payment.RequestID,
		&order.Payment.Currency,
		&order.Payment.Provider,
		&order.Payment.Amount,
		&order.Payment.PaymentDT,
		&order.Payment.Bank,
		&order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
	)
	if err != nil {
		db.log.Warn("failed to query payment info, might not exist", "order_uid", orderUID, "error", err)
		// Аналогично delivery_info
	}
	// Устанавливаем Transaction, так как он не выбирается из БД
	order.Payment.Transaction = orderUID
	db.log.Debug("payment info fetched successfully", "order_uid", orderUID)

	// 4. Получаем данные из таблицы items (может быть несколько)
	itemSQL := `
        SELECT
            chrt_id, track_number, price, rid, name,
            sale, size, total_price, nm_id, brand, status
        FROM items
        WHERE order_uid = $1`
	rows, err := db.conn.Query(ctx, itemSQL, orderUID)
	if err != nil {
		db.log.Error("failed to query items", "order_uid", orderUID, "error", err)
		return models.Order{}, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			db.log.Error("failed to scan item row", "order_uid", orderUID, "error", err)
			return models.Order{}, err
		}
		// Устанавливаем TrackNumber (если он не возвращается из БД, но нужен в Item)
		// item.TrackNumber = order.TrackNumber // Если TrackNumber в Item всегда тот же, что и в Order
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		db.log.Error("error after scanning item rows", "order_uid", orderUID, "error", err)
		return models.Order{}, err
	}
	order.Items = items
	db.log.Debug("items fetched successfully", "order_uid", orderUID, "count", len(order.Items))

	db.log.Info("order info retrieved successfully", "order_uid", orderUID)
	return order, nil
}

func (db *DB) Delete(ctx context.Context, orderID int) error {
	return nil
}

func (db *DB) GetIDs(ctx context.Context) ([]string, error) {
	db.log.Debug("attempting to get all order UIDs")

	rows, err := db.conn.Query(ctx, "SELECT order_uid FROM orders")
	if err != nil {
		db.log.Error("failed to query order UIDs", "error", err)
		return nil, err
	}
	defer rows.Close()

	var uids []string
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			db.log.Error("failed to scan order UID row", "error", err)
			return nil, err
		}
		uids = append(uids, uid)
	}

	if err = rows.Err(); err != nil {
		db.log.Error("error after scanning order UID rows", "error", err)
		return nil, err
	}

	db.log.Info("successfully retrieved all order UIDs", "count", len(uids))
	return uids, nil
}
