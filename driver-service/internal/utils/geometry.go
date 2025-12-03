package utils

import (
	"math"
)

const R = 6371 // dünyanın km cinsinden radius değeri

// dünya üzerinde ki iki noktayı konumsal olarak(lat ve lon) enlem, boylam şeklinde hesaplama işlemi.
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	var (
		dLat = (lat2 - lat1) * math.Pi / 180.0
		dLon = (lon2 - lon1) * math.Pi / 180.0

		a = math.Sin(dLat/2)*math.Sin(dLat/2) +
			math.Cos(lat1*math.Pi/180.0)*math.Cos(lat2*math.Pi/180.0)*
				math.Sin(dLon/2)*math.Sin(dLon/2)
		c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	)

	return R * c
}
