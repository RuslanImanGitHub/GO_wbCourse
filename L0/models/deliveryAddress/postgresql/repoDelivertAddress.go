package postgresql

import (
	"L0/models/deliveryAddress"
	"L0/models/deliveryAddress/storage"
	"L0/pkg/dbClient/postgresql"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
)

type repository struct {
	client postgresql.Client
}

func (r *repository) Create(ctx context.Context, delivery *deliveryAddress.DeliveryAddress) (delivery_id int, err error) {
	q := `INSERT INTO Deliveries (delivery_id, order_name, phone, zip, city, address, region, email)
		  VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING delivery_id`

	if err = r.client.QueryRow(ctx, q,
		delivery.Delivery_id, delivery.Name, delivery.Phone, delivery.Zip,
		delivery.City, delivery.Address, delivery.Region,
		delivery.Email).Scan(&delivery_id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf("%s", fmt.Sprintf("sql Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			return 0, newErr
		}
		return 0, err
	}
	return delivery_id, nil
}

func (r *repository) FindAll(ctx context.Context) ([]deliveryAddress.DeliveryAddress, error) {
	q := `SELECT d.delivery_id,
				d.order_name, d.phone, d.zip, d.city,
				d.address, d.region, d.email
		FROM deliveries d`

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	var deliveries []deliveryAddress.DeliveryAddress
	for rows.Next() {
		var delivery deliveryAddress.DeliveryAddress
		err = rows.Scan(&delivery.Delivery_id, &delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City,
			&delivery.Address, &delivery.Region, &delivery.Email)
		if err != nil {
			return nil, err
		}

		deliveries = append(deliveries, delivery)
	}
	return deliveries, nil
}

func (r *repository) FindFromId(ctx context.Context, id int) (deliveryAddress.DeliveryAddress, error) {
	var delivery deliveryAddress.DeliveryAddress
	q := `SELECT d.delivery_id,
				d.order_name, d.phone, d.zip, d.city,
				d.address, d.region, d.email
		FROM Deliveries d
		WHERE d.delivery_id = $1`

	rows := r.client.QueryRow(ctx, q, id)

	err := rows.Scan(&delivery.Delivery_id, &delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City,
		&delivery.Address, &delivery.Region, &delivery.Email)
	if err != nil {
		return delivery, err
	}

	return delivery, nil
}

func NewRepository(client postgresql.Client) storage.Repository {
	return &repository{client: client}
}
