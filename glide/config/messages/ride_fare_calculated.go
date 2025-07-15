package messages

import "github.com/brianvoe/gofakeit/v6"

func GenerateRideFareCalculatedMessage() interface{} {
	return map[string]interface{}{
		"ride_id":          gofakeit.UUID(),
		"fare_amount":      gofakeit.Float64Range(5.0, 50.0),
		"distance_km":      gofakeit.Float64Range(1.0, 25.0),
		"duration_min":     gofakeit.IntRange(5, 45),
		"surge_multiplier": gofakeit.Float64Range(1.0, 2.5),
		"calculated_at":    gofakeit.Date(),
	}
}
