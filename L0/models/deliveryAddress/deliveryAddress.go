package deliveryAddress

type DeliveryAddress struct {
	Delivery_id int    `json:"delivery_id" fake:"{number:1,2147483646}"`
	Name        string `json:"name" fake:"{firstname}"`
	Phone       string `json:"phone" fake:"{phone}"`
	Zip         string `json:"zip" fake:"{zip}"`
	City        string `json:"city" fake:"{city}"`
	Address     string `json:"address" fake:"{address}"`
	Region      string `json:"region" fake:"{country}"`
	Email       string `json:"email" fake:"{email}"`
}
