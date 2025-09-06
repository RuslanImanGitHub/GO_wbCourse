package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var (
	host     = viper.GetString("postgreConfig.host")
	port     = viper.GetString("postgreConfig.port")
	user     = viper.GetString("postgreConfig.user")
	password = viper.GetString("postgreConfig.password")
	dbname   = viper.GetString("postgreConfig.dbname")
)

func connectToDB() (sql.DB, error) {
	constr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	db, err := sql.Open("postgres", constr)

	//Check conn
	if err != nil {
		return *db, err
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		return *db, err
	}
	return *db, nil
}

func GetOrders() ([]Order, error) {

	db, err := connectToDB()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	orderQuery := `SELECT a.order_uid, a.track_number, a.entry,
							 a.delivery_id,

							 d.order_name, d.phone, d.zip, d.city,
							 d.address, d.region, d.email,

							 p.pay_transaction, p.request_id, p.currency, p.provider,
							 p.amount, p.payment_dt, p.bank, p.delivery_cost,
							 p.goods_total, p.custom_fee,

							 a.locale, a.internal_signature,
							 a.customer_id, a.delivery_service, a.shard_key,
							 a.sm_id, a.date_created, a.oof_shard 
					  FROM Orders a
					  INNER JOIN Deliveries d ON a.delivery_id = d.delivery_id
					  INNER JOIN Payments p ON a.order_uid = p.pay_transaction`

	result := []Order{}

	rows, err := db.QueryContext(ctx, orderQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var order Order
	for rows.Next() {
		err = rows.Scan(&order.Order_uid, &order.Track_number, &order.Entry,
			&order.Delivery.Delivery_id,
			&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City,
			&order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,
			&order.Payment.Transaction, &order.Payment.Request_id, &order.Payment.Currency, &order.Payment.Provider,
			&order.Payment.Amount, &order.Payment.Payment_dt, &order.Payment.Bank, &order.Payment.Delivery_cost,
			&order.Payment.Goods_total, &order.Payment.Custom_fee,
			&order.Locale, &order.Internal_signature, &order.Customer_id, &order.Delivery_service,
			&order.Shardkey, &order.Sm_id, &order.Date_created, &order.Oof_shard)

		if err != nil {
			return nil, err
		}
		order.Items, err = getOrderItemsRelatedToOrder(tx, &ctx, order.Track_number)
		if err != nil {
			continue
		}
		result = append(result, order)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return result, err
}

func GetOrderById(orderUid string) (Order, error) {
	var order Order
	db, err := connectToDB()
	if err != nil {
		return order, err
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return order, err
	}

	orderQuery := `SELECT a.order_uid, a.track_number, a.entry,
							 a.delivery_id,

							 d.order_name, d.phone, d.zip, d.city,
							 d.address, d.region, d.email,

							 p.pay_transaction, p.request_id, p.currency, p.provider,
							 p.amount, p.payment_dt, p.bank, p.delivery_cost,
							 p.goods_total, p.custom_fee,

							 a.locale, a.internal_signature,
							 a.customer_id, a.delivery_service, a.shard_key,
							 a.sm_id, a.date_created, a.oof_shard 
					  FROM Orders a
					  INNER JOIN Deliveries d ON a.delivery_id = d.delivery_id
					  INNER JOIN Payments p ON a.order_uid = p.pay_transaction
					  
					  WHERE a.order_uid = $1`

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // re-throw panic
		} else if err != nil { // err from operations within the transaction
			tx.Rollback()
		}
	}()

	err = db.QueryRowContext(ctx, orderQuery, orderUid).Scan(&order.Order_uid, &order.Track_number, &order.Entry,
		&order.Delivery.Delivery_id,
		&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City,
		&order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,
		&order.Payment.Transaction, &order.Payment.Request_id, &order.Payment.Currency, &order.Payment.Provider,
		&order.Payment.Amount, &order.Payment.Payment_dt, &order.Payment.Bank, &order.Payment.Delivery_cost,
		&order.Payment.Goods_total, &order.Payment.Custom_fee,
		&order.Locale, &order.Internal_signature, &order.Customer_id, &order.Delivery_service,
		&order.Shardkey, &order.Sm_id, &order.Date_created, &order.Oof_shard)

	if err != nil {
		if err == sql.ErrNoRows {
			return order, err
			//log.Fatalf("No rows with order_uid %s", orderUid)
		}
		return order, err
		//log.Fatal(err)
	}
	order.Items, err = getOrderItemsRelatedToOrder(tx, &ctx, order.Track_number)
	if err != nil {
		return order, err
	}

	err = tx.Commit()
	if err != nil {
		return order, err
		//log.Fatal(err)
	}
	return order, nil
}

func getOrderItemsRelatedToOrder(tx *sql.Tx, ctx *context.Context, track_number string) ([]ItemDetails, error) {
	result := make([]ItemDetails, 0)
	query := `SELECT oi.chrt_id, oi.track_number, oi.price, oi.rid,
				 	 oi.item_name, oi.sale, oi.item_size, oi.total_PRICE,
					 oi.nm_id, oi.brand, oi.status
			  FROM OrderItems oi 
			  WHERE oi.track_number =$1`
	rows, err := tx.QueryContext(*ctx, query, track_number)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		oi := ItemDetails{}
		rows.Scan(&oi.Chrt_id, &oi.Track_number, &oi.Price, &oi.Rid,
			&oi.Name, &oi.Sale, &oi.Size, &oi.Total_price,
			&oi.Nm_id, &oi.Brand, oi.Status)
		result = append(result, oi)
	}
	return result, nil
}

func InsertOrderToDB(order Order) error {

	db, err := connectToDB()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	deliveryQuery := `INSERT INTO Deliveries (order_name, phone, zip, city,
											  address, region, email)
					  VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING delivery_id`

	orderQuery := `INSERT INTO orders (order_uid, track_number, entry,
	 							delivery_id, locale, internal_signature,
								customer_id, delivery_service, shard_key,
								sm_id, date_created, oof_shard)
					  VALUES($1, $2, $3,
					  	     $4, $5, $6,
					  	     $7, $8, $9,
					  	     $10, $11, $12)`

	paymentQuery := `INSERT INTO Payments (pay_transaction, request_id, currency,
										   provider, amount, payment_dt,
										   bank, delivery_cost, goods_total,
										   custom_fee)
					  VALUES($1, $2, $3,
							 $4, $5, $6,
							 $7, $8, $9,
							 $10)`

	ItemsQuery := `INSERT INTO OrderItems (track_number, price, rid,
										   item_name, sale, item_size,
										   totaL_PRICE, nm_id, brand,
										   status)
					  VALUES($1, $2, $3,
							 $4, $5, $6,
							 $7, $8, $9,
							 $10)`
	var newDeliveyId int

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // re-throw panic
		} else if err != nil { // err from operations within the transaction
			tx.Rollback()
		}
	}()

	err = tx.QueryRowContext(ctx, deliveryQuery,
		order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region,
		order.Delivery.Email).Scan(&newDeliveyId)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, orderQuery,
		order.Order_uid, order.Track_number, order.Entry,
		newDeliveyId, order.Locale, order.Internal_signature,
		order.Customer_id, order.Delivery_service, order.Shardkey,
		order.Sm_id, order.Date_created, order.Oof_shard)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, paymentQuery,
		order.Order_uid, order.Payment.Request_id, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.Payment_dt,
		order.Payment.Bank, order.Payment.Delivery_cost, order.Payment.Goods_total,
		order.Payment.Custom_fee)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx, ItemsQuery,
			order.Track_number, item.Price, item.Rid,
			item.Name, item.Sale, item.Size,
			item.Total_price, item.Nm_id, item.Brand,
			item.Status)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
