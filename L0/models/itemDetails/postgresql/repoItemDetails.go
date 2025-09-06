package postgresql

import (
	itemdetails "L0/models/itemDetails"
	"L0/models/itemDetails/storage"
	"L0/models/order"
	"L0/pkg/dbClient/postgresql"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
)

type repository struct {
	client postgresql.Client
}

func (r *repository) Create(ctx context.Context, item *itemdetails.ItemDetails) (err error) {
	q := `INSERT INTO OrderItems (track_number, price, rid,
										item_name, sale, item_size,
										totaL_PRICE, nm_id, brand,
										status)
					VALUES($1, $2, $3,
							$4, $5, $6,
							$7, $8, $9,
							$10)`

	tag, err := r.client.Exec(ctx, q,
		item.Track_number, item.Price, item.Rid,
		item.Name, item.Sale, item.Size,
		item.Total_price, item.Nm_id, item.Brand,
		item.Status)
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

func (r *repository) FindAllFromOrder(ctx context.Context, order order.Order) ([]itemdetails.ItemDetails, error) {
	result := make([]itemdetails.ItemDetails, 0)
	query := `SELECT oi.chrt_id, oi.track_number, oi.price, oi.rid,
				 	 oi.item_name, oi.sale, oi.item_size, oi.total_PRICE,
					 oi.nm_id, oi.brand, oi.status
			  FROM OrderItems oi 
			  WHERE oi.track_number = $1`
	rows, err := r.client.Query(ctx, query, order.Track_number)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		oi := itemdetails.ItemDetails{}
		rows.Scan(&oi.Chrt_id, &oi.Track_number, &oi.Price, &oi.Rid,
			&oi.Name, &oi.Sale, &oi.Size, &oi.Total_price,
			&oi.Nm_id, &oi.Brand, oi.Status)
		result = append(result, oi)
	}
	return result, nil
}

func NewRepository(client postgresql.Client) storage.Repository {
	return &repository{client: client}
}
