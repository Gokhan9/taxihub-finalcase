package utils

import (
	"math"
	"testing"
)

/*
İstanbul - Ankara arası mesafe hesabı. İki nokta aynı ise mesafe değeri 0 olmalı.
Burada İstanbul ile Ankara arası hata payını da göz önünde bulundurarak mesafeyi doğru hesaplayıp hesaplamadığını test ediyoruz.
*/
func TestHaversineDistance(t *testing.T) {

	tests := []struct {
		name     string
		lat1     float64
		lon1     float64
		lat2     float64
		lon2     float64
		expected float64
		epsilon  float64 // hata payı olma ihtimaline karşı eklenen test değer.
	}{
		{
			name:     "test senaryosu",
			lat1:     41.0082,
			lon1:     28.9784,
			lat2:     41.0082,
			lon2:     28.9784,
			expected: 0.0,
			epsilon:  0.0001,
		},
		{
			name:     "İst - Ank",
			lat1:     41.0082, //ist
			lon1:     28.9784,
			lat2:     39.9334, //ank
			lon2:     32.8597,
			expected: 351.0, //Aşağı yukarı "İstanbul-Ankara" arası mesafe
			epsilon:  5.0,   //Hata payı
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := HaversineDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			if math.Abs(g-tt.expected) > tt.epsilon {
				t.Errorf("Haversine Distance Değeri : %v, beklenen değer: %v (epsilon: %v)", g, tt.expected, tt.epsilon)
			}
		})
	}
}

/*
A'dan B'ye gidilen mesafe ile B'den A'ya gidilen mesafe eşit olmalı
*/
func TestHaversineDistanceSymmetry(t *testing.T) {

	lat1, lon1 := 40.990, 29.025 // ! A Noktası için
	lat2, lon2 := 41.037, 28.974 // ! B Noktası için

	distance1 := HaversineDistance(lat1, lon1, lat2, lon2)
	distance2 := HaversineDistance(lat2, lon2, lat1, lon1)

	if distance1 != distance2 {
		t.Errorf("Mesafe simetrik değil. A-B: %v, B-A: %v", distance1, distance2)
	}
}
