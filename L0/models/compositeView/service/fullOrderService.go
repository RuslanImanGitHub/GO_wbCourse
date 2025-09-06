package service

import (
	"L0/models/compositeView/storage"
	"L0/models/deliveryAddress"
	deliveryService "L0/models/deliveryAddress/service"
	itemdetails "L0/models/itemDetails"
	itemDetailsService "L0/models/itemDetails/service"
	"L0/models/order"
	orderService "L0/models/order/service"
	"L0/models/payment"
	paymentService "L0/models/payment/service"
	"L0/pkg/dbClient/postgresql"
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type repository struct {
	client       postgresql.Client
	deliveryS    deliveryService.Service
	orderS       orderService.Service
	paymentS     paymentService.Service
	itemDetailsS itemDetailsService.Service
}

// Create implements storage.Repository.
func (r *repository) Create(ctx context.Context, pool *pgxpool.Pool, order *order.Order, delivery *deliveryAddress.DeliveryAddress, payment *payment.Payment, itemDetails []itemdetails.ItemDetails) error {
	return pool.BeginFunc(ctx, func(pgx.Tx) error {
		deliveryId, err := r.deliveryS.Create(ctx, delivery)
		if err != nil {
			return err
		}
		if order.Delivery_id != deliveryId {
			return errors.New("order's deliveryId doesn't match with Delivery's Id")
		}
		err = r.orderS.Create(ctx, order)
		if err != nil {
			return err
		}
		err = r.paymentS.Create(ctx, payment)
		if err != nil {
			return err
		}
		for _, value := range itemDetails {
			err = r.itemDetailsS.Create(ctx, &value)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// FindAll implements storage.Repository.
func (r *repository) FindAll(ctx context.Context, pool *pgxpool.Pool) ([]order.Order, error) {
	var result []order.Order
	orders, err := r.orderS.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, order := range orders {
		//Get Delivery (MAY BE NULL)
		delivery, err := r.deliveryS.FindFromId(ctx, order.Delivery_id)
		if err != nil {
			continue
		}
		order.Delivery = delivery

		//Get Payment (CAN'T BE NULL)
		pay, err := r.paymentS.FindOneFromTransaction(ctx, order.Order_uid)
		if err != nil {
			return nil, err
		}
		order.Payment = pay

		//Get OrderItems (CAN'T BE NULL)
		ois, err := r.itemDetailsS.FindAllFromOrder(ctx, order)
		if err != nil {
			return nil, err
		}
		order.Items = ois
		result = append(result, order)
	}
	return result, nil
}

// FindAll implements storage.Repository.
func (r *repository) FindN(ctx context.Context, pool *pgxpool.Pool, n int) ([]order.Order, error) {
	var result []order.Order
	orders, err := r.orderS.FindN(ctx, n)
	if err != nil {
		return nil, err
	}
	for _, order := range orders {
		//Get Delivery (MAY BE NULL)
		delivery, err := r.deliveryS.FindFromId(ctx, order.Delivery_id)
		if err != nil {
			continue
		}
		order.Delivery = delivery

		//Get Payment (CAN'T BE NULL)
		pay, err := r.paymentS.FindOneFromTransaction(ctx, order.Order_uid)
		if err != nil {
			return nil, err
		}
		order.Payment = pay

		//Get OrderItems (CAN'T BE NULL)
		ois, err := r.itemDetailsS.FindAllFromOrder(ctx, order)
		if err != nil {
			return nil, err
		}
		order.Items = ois
		result = append(result, order)
	}
	return result, nil
}

// FindOne implements storage.Repository.
func (r *repository) FindOne(ctx context.Context, pool *pgxpool.Pool, id string) (order.Order, error) {
	order, err := r.orderS.FindOne(ctx, id)
	if err != nil {
		return order, err
	}
	//Get Delivery (MAY BE NULL)
	delivery, err := r.deliveryS.FindFromId(ctx, order.Delivery_id)
	if err != nil {
	}
	order.Delivery = delivery

	//Get Payment (CAN'T BE NULL)
	pay, err := r.paymentS.FindOneFromTransaction(ctx, order.Order_uid)
	if err != nil {
		return order, err
	}
	order.Payment = pay

	//Get OrderItems (CAN'T BE NULL)
	ois, err := r.itemDetailsS.FindAllFromOrder(ctx, order)
	if err != nil {
		return order, err
	}
	order.Items = ois

	return order, nil
}

func NewService(client postgresql.Client, deliveryS deliveryService.Service, orderS orderService.Service, paymentS paymentService.Service, itemDetailsS itemDetailsService.Service) storage.Repository {
	return &repository{client: client, deliveryS: deliveryS, orderS: orderS, paymentS: paymentS, itemDetailsS: itemDetailsS}
}
