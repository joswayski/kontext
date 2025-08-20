package messages

import "github.com/brianvoe/gofakeit/v6"

func GenerateRideCompletedMessage() interface{} {
	return map[string]interface{}{
		"ride_id":      gofakeit.UUID(),
		"driver_id":    gofakeit.UUID(),
		"rider_id":     gofakeit.UUID(),
		"completed_at": gofakeit.Date(),
		"fare_amount":  gofakeit.Float64Range(5.0, 50.0),
		"distance_km":  gofakeit.Float64Range(1.0, 25.0),
		"duration_min": gofakeit.IntRange(5, 45),
		"dropoff_location": map[string]interface{}{
			"latitude":  gofakeit.Latitude(),
			"longitude": gofakeit.Longitude(),
		},
	}
}
