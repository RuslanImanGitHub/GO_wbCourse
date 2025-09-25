package payment

//No need for DTO, cause DTO == model
type Payment struct {
	Transaction   string  `json:"pay_transaction" fake:"skip" validate:"max=100"`
	Request_id    string  `json:"request_id" fake:"{uuid}" validate:"uuid"`
	Currency      string  `json:"currency" fake:"{currencylong}" validate:"required,max=20"`
	Provider      string  `json:"provider" fake:"{company}" validate:"required,max=100"`
	Amount        float64 `json:"amount" fake:"{number:1,10000}" validate:"required,numeric"`
	Payment_dt    int     `json:"payment_dt" fake:"{number:1,10000}" validate:"required,numeric"`
	Bank          string  `json:"bank" fake:"{bankname}" validate:"required,max=100"`
	Delivery_cost float64 `json:"delivery_cost" fake:"{number:1,10000}" validate:"required,numeric"`
	Goods_total   float64 `json:"goods_total" fake:"{number:1,10000}" validate:"required,numeric"`
	Custom_fee    float64 `json:"custom_fee" fake:"{number:1,10000}" validate:"required,numeric"`
}
