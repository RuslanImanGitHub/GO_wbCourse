package deliveryAddress

type DeliveryAddress struct {
	Delivery_id int    `json:"delivery_id" fake:"{number:1,2147483646}" validate:"required"`
	Name        string `json:"name" fake:"{firstname}" validate:"alpha,max=100"`
	Phone       string `json:"phone" fake:"{phone}" validate:"max=100"`
	Zip         string `json:"zip" fake:"{zip}" validate:"max=100"`
	City        string `json:"city" fake:"{city}" validate:"max=100"`
	Address     string `json:"address" fake:"{address}" validate:"required,max=1000"`
	Region      string `json:"region" fake:"{country}" validate:"max=100"`
	Email       string `json:"email" fake:"{email}" validate:"required,email,max=100"`
}
