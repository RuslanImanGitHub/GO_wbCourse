package itemdetails

type ItemDetails struct {
	Chrt_id      int     `json:"chrt_id" validate:"required"`
	Track_number string  `json:"track_number" validate:"required,max=100"`
	Price        float64 `json:"price" fake:"{number:1,1000}" validate:"required"`
	Rid          string  `json:"rid" fake:"{regex:[a-z]{20}}" validate:"required,max=100"`
	Name         string  `json:"name" fake:"{productname}" validate:"required,max=100"`
	Sale         int     `json:"sale" fake:"{number:1,100}" validate:"required"`
	Size         string  `json:"size" fake:"{regex:[1-9]{1}}" validate:"required"`
	Total_price  float64 `json:"total_price" fake:"{number:1,10000}" validate:"required"`
	Nm_id        int     `json:"nm_id" fake:"{number:1,100}" validate:"required"`
	Brand        string  `json:"brand" fake:"{company}" validate:"required,max=100"`
	Status       int     `json:"status" fake:"{number:1,10}" validate:"required"`
}
