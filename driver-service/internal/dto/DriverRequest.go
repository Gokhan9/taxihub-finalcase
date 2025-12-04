package dto

/*
Yeni bir sürücü kaydı oluşturulurken, API'ye gönderilen kişisel, araç ve konum bilgierini içeren istek modeli.
*/
type DriverCreateRequest struct {
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Plate     string  `json:"plate"`
	TaxiType  string  `json:"taxiType"`
	CarBrand  string  `json:"carBrand"`
	CarModel  string  `json:"carModel"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
}
