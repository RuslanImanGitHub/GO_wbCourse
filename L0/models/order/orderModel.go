package order

import (
	deliveryaddress "L0/models/deliveryAddress"
	itemdetails "L0/models/itemDetails"
	"L0/models/payment"
	"time"
)

type Order struct {
	Order_uid          string                          `json:"order_uid" fake:"{regex:[a-z]{20}}"`
	Track_number       string                          `json:"track_number" fake:"{regex:[a-z]{20}}"`
	Entry              string                          `json:"entry" fake:"{regex:[a-z]{4}}"`
	Delivery_id        int                             `json:"delivery_id" fake:"skip"`
	Delivery           deliveryaddress.DeliveryAddress `json:"delivery" fake:"skip"`
	Payment            payment.Payment                 `json:"payment" fake:"skip"`
	Items              []itemdetails.ItemDetails       `json:"items" fake:"skip"`
	Locale             string                          `json:"locale" fake:"languageabbreviation"`
	Internal_signature string                          `json:"internal_signature" fake:"{regex:[a-z]{20}}"`
	Customer_id        string                          `json:"customer_id" fake:"{uuid}"`
	Delivery_service   string                          `json:"delivery_service" fake:"{company}"`
	Shardkey           string                          `json:"shardkey" fake:"{number:1,100}"`
	Sm_id              int                             `json:"sm_id" fake:"{number:1,100}"`
	Date_created       time.Time                       `json:"date_created"`
	Oof_shard          string                          `json:"oof_shard" fake:"{number:1,100}"`
}
