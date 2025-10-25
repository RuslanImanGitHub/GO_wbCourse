package postgresql

import (
	"L0/backend/models/payment"
	"L0/backend/models/payment/storage"
	"L0/backend/pkg/dbClient/postgresql"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
)

type repository struct {
	client postgresql.Client
}

func (r *repository) Create(ctx context.Context, payment *payment.Payment) error {
	q := `INSERT INTO L0.Payments (pay_transaction, request_id, currency,
										   provider, amount, payment_dt,
										   bank, delivery_cost, goods_total,
										   custom_fee)
					  VALUES($1, $2, $3,
							 $4, $5, $6,
							 $7, $8, $9,
							 $10)`

	tag, err := r.client.Exec(ctx, q,
		payment.Transaction, payment.Request_id, payment.Currency,
		payment.Provider, payment.Amount, payment.Payment_dt,
		payment.Bank, payment.Delivery_cost, payment.Goods_total,
		payment.Custom_fee)
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

func (r *repository) FindAll(ctx context.Context) ([]payment.Payment, error) {
	q := `SELECT 
			p.pay_transaction, p.request_id, p.currency, p.provider,
			p.amount, p.payment_dt, p.bank, p.delivery_cost,
			p.goods_total, p.custom_fee

		FROM L0.Payments p`

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	var payments []payment.Payment
	for rows.Next() {
		var payment payment.Payment
		err = rows.Scan(&payment.Transaction, &payment.Request_id, &payment.Currency, &payment.Provider,
			&payment.Amount, &payment.Payment_dt, &payment.Bank, &payment.Delivery_cost,
			&payment.Goods_total, &payment.Custom_fee)
		if err != nil {
			return nil, err
		}

		payments = append(payments, payment)
	}
	return payments, nil
}

func (r *repository) FindOneFromTransaction(ctx context.Context, transaction string) (payment.Payment, error) {
	var payment payment.Payment
	q := `SELECT 
			p.pay_transaction, p.request_id, p.currency, p.provider,
			p.amount, p.payment_dt, p.bank, p.delivery_cost,
			p.goods_total, p.custom_fee

		FROM L0.Payments p 
		WHERE p.pay_transaction = $1`
	rows := r.client.QueryRow(ctx, q, transaction)
	err := rows.Scan(&payment.Transaction, &payment.Request_id, &payment.Currency, &payment.Provider,
		&payment.Amount, &payment.Payment_dt, &payment.Bank, &payment.Delivery_cost,
		&payment.Goods_total, &payment.Custom_fee)
	if err != nil {
		return payment, err
	}

	return payment, nil
}

func NewRepository(client postgresql.Client) storage.Repository {
	return &repository{client: client}
}
