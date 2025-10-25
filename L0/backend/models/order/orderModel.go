package order

import (
	deliveryaddress "L0/backend/models/deliveryAddress"
	itemdetails "L0/backend/models/itemDetails"
	"L0/backend/models/payment"
	"time"
)

type Order struct {
	Order_uid          string                          `json:"order_uid" fake:"{regex:[a-z]{20}}" validate:"required,max=100"`
	Track_number       string                          `json:"track_number" fake:"{regex:[a-z]{20}}" validate:"required,max=100"`
	Entry              string                          `json:"entry" fake:"{regex:[a-z]{4}}" validate:"required,max=100"`
	Delivery_id        int                             `json:"delivery_id" fake:"skip"`
	Delivery           deliveryaddress.DeliveryAddress `json:"delivery" fake:"skip"`
	Payment            payment.Payment                 `json:"payment" fake:"skip"`
	Items              []itemdetails.ItemDetails       `json:"items" fake:"skip"`
	Locale             string                          `json:"locale" fake:"languageabbreviation" validate:"required,max=20"`
	Internal_signature string                          `json:"internal_signature" fake:"{regex:[a-z]{20}}" validate:"required,max=100"`
	Customer_id        string                          `json:"customer_id" fake:"{uuid}" validate:"required,uuid"`
	Delivery_service   string                          `json:"delivery_service" fake:"{company}" validate:"max=100"`
	Shardkey           string                          `json:"shardkey" fake:"{number:1,100}" validate:"alphanum"`
	Sm_id              int                             `json:"sm_id" fake:"{number:1,100}" validate:"numeric,max=100"`
	Date_created       time.Time                       `json:"date_created" validate:"required"`
	Oof_shard          string                          `json:"oof_shard" fake:"{number:1,100}" validate:"alphanum"`
}
