package dto

// İsteğe bağlı "driver" bilgileri üzerinde update.
type DriverUpdateRequest struct {
	FirstName *string   `json:"firstName"`
	LastName  *string   `json:"lastName"`
	Plate     *string   `json:"plate"`
	Location  *Location `json:"location"`
}

//"Driver" konumuna ait lokasyon bilgisi.
type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type DriverStatusRequest struct {
	Status string `json:"status"`
}
