package messages

import "github.com/brianvoe/gofakeit"

func GenerateDriverLocationUpdatedMessage() interface{} {
	return map[string]interface{}{
		"driver_id": gofakeit.UUID(),
		"location": map[string]interface{}{
			"latitude":  gofakeit.Latitude(),
			"longitude": gofakeit.Longitude(),
		},
		"timestamp": gofakeit.Date(),
		"speed":     gofakeit.Float64Range(0, 80),
		"heading":   gofakeit.Float64Range(0, 360),
	}
}
