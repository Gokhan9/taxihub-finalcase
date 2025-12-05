package dto

// NearbyDriverResponse istemciye dönülecek özelleştirilmiş veri yapısıdır.
type NearbyDriverResponse struct {
	FirstName  string  `json:"firstName"`
	LastName   string  `json:"lastName"`
	Plate      string  `json:"plate"`
	DistanceKm float64 `json:"distanceKm"` // km cinsinden mesafe
}
