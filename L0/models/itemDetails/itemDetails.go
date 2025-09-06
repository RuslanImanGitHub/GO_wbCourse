package itemdetails

type ItemDetails struct {
	Chrt_id      int     `json:"chrt_id"`
	Track_number string  `json:"track_number"`
	Price        float64 `json:"price" fake:"{number:1,1000}"`
	Rid          string  `json:"rid" fake:"{regex:[a-z]{20}}"`
	Name         string  `json:"name" fake:"{productname}"`
	Sale         int     `json:"sale" fake:"{number:1,100}"`
	Size         string  `json:"size" fake:"{regex:[1-9]{1}}"`
	Total_price  float64 `json:"total_price" fake:"{number:1,10000}"`
	Nm_id        int     `json:"nm_id" fake:"{number:1,100}"`
	Brand        string  `json:"brand" fake:"{company}"`
	Status       int     `json:"status" fake:"{number:1,10}"`
}
