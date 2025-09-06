package payment

//No need for DTO, cause DTO == model
type Payment struct {
	Transaction   string  `json:"pay_transaction" fake:"skip"`
	Request_id    string  `json:"request_id" fake:"{uuid}"`
	Currency      string  `json:"currency" fake:"{currencylong}"`
	Provider      string  `json:"provider" fake:"{company}"`
	Amount        float64 `json:"amount" fake:"{number:1,10000}"`
	Payment_dt    int     `json:"payment_dt" fake:"{number:1,10000}"`
	Bank          string  `json:"bank" fake:"{bankname}"`
	Delivery_cost float64 `json:"delivery_cost" fake:"{number:1,10000}"`
	Goods_total   float64 `json:"goods_total" fake:"{number:1,10000}"`
	Custom_fee    float64 `json:"custom_fee" fake:"{number:1,10000}"`
}
