package dto

type DriverUpdateRequest struct {
	FirstName *string   `json:"firstName"`
	LastName  *string   `json:"lastName"`
	Plate     *string   `json:"plate"`
	Location  *Location `json:"location"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
