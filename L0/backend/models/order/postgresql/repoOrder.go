package postgresql

import (
	"L0/backend/models/order"
	"L0/backend/models/order/storage"
	"L0/backend/pkg/dbClient/postgresql"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
)

type repository struct {
	client postgresql.Client
}

// Create implements storage.Repository.
func (r *repository) Create(ctx context.Context, order *order.Order) error {
	q := `INSERT INTO L0.Orders (order_uid, track_number, entry,
	 							delivery_id, locale, internal_signature,
								customer_id, delivery_service, shard_key,
								sm_id, date_created, oof_shard)
					  VALUES($1, $2, $3,
					  	     $4, $5, $6,
					  	     $7, $8, $9,
					  	     $10, $11, $12)`

	tag, err := r.client.Exec(ctx, q,
		order.Order_uid, order.Track_number, order.Entry,
		order.Delivery_id, order.Locale, order.Internal_signature,
		order.Customer_id, order.Delivery_service, order.Shardkey,
		order.Sm_id, order.Date_created, order.Oof_shard)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf("%s", fmt.Sprintf("sql Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			return newErr
		}
		return err
	}
	fmt.Printf("Rows affected: %d\n", tag.RowsAffected())
	return nil
}

// FindAll implements storage.Repository.
func (r *repository) FindAll(ctx context.Context) ([]order.Order, error) {
	result := make([]order.Order, 0)
	q := `SELECT a.order_uid, a.track_number, a.entry,
							 a.delivery_id,
							 a.locale, a.internal_signature,
							 a.customer_id, a.delivery_service, a.shard_key,
							 a.sm_id, a.date_created, a.oof_shard 
					  FROM L0.Orders a`
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		order := order.Order{}
		rows.Scan(&order.Order_uid, &order.Track_number, &order.Entry,
			&order.Delivery_id,
			&order.Locale, &order.Internal_signature, &order.Customer_id, &order.Delivery_service,
			&order.Shardkey, &order.Sm_id, &order.Date_created, &order.Oof_shard)
		result = append(result, order)
	}
	return result, nil
}

func (r *repository) FindN(ctx context.Context, n int) ([]order.Order, error) {
	result := make([]order.Order, 0)
	q := `SELECT a.order_uid, a.track_number, a.entry,
							 a.delivery_id,
							 a.locale, a.internal_signature,
							 a.customer_id, a.delivery_service, a.shard_key,
							 a.sm_id, a.date_created, a.oof_shard 
					  FROM L0.Orders a
					  LIMIT $1`
	rows, err := r.client.Query(ctx, q, n)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		order := order.Order{}
		rows.Scan(&order.Order_uid, &order.Track_number, &order.Entry,
			&order.Delivery_id,
			&order.Locale, &order.Internal_signature, &order.Customer_id, &order.Delivery_service,
			&order.Shardkey, &order.Sm_id, &order.Date_created, &order.Oof_shard)
		result = append(result, order)
	}
	return result, nil
}

// FindOne implements storage.Repository.
func (r *repository) FindOne(ctx context.Context, id string) (order.Order, error) {
	var order order.Order
	q := `SELECT a.order_uid, a.track_number, a.entry,
							 a.delivery_id,

							 a.locale, a.internal_signature,
							 a.customer_id, a.delivery_service, a.shard_key,
							 a.sm_id, a.date_created, a.oof_shard 
					  FROM L0.Orders a
					  WHERE a.order_uid = $1`

	rows := r.client.QueryRow(ctx, q, id)

	err := rows.Scan(&order.Order_uid, &order.Track_number, &order.Entry,
		&order.Delivery_id,
		&order.Locale, &order.Internal_signature, &order.Customer_id, &order.Delivery_service,
		&order.Shardkey, &order.Sm_id, &order.Date_created, &order.Oof_shard)
	if err != nil {
		return order, err
	}

	return order, nil
}

func NewRepository(client postgresql.Client) storage.Repository {
	return &repository{client: client}
}
