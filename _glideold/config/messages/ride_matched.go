package messages

import "github.com/brianvoe/gofakeit/v6"

func GenerateRideMatchedMessage() interface{} {
	return map[string]interface{}{
		"ride_id":     gofakeit.UUID(),
		"driver_id":   gofakeit.UUID(),
		"rider_id":    gofakeit.UUID(),
		"matched_at":  gofakeit.Date(),
		"eta_minutes": gofakeit.IntRange(2, 15),
	}
}
