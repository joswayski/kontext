package messages

import "github.com/brianvoe/gofakeit"

func GenerateDriverActivatedMessage() interface{} {
	return map[string]interface{}{
		"driver_id":    gofakeit.UUID(),
		"activated_at": gofakeit.Date(),
		"status":       "active",
		"location": map[string]interface{}{
			"latitude":  gofakeit.Latitude(),
			"longitude": gofakeit.Longitude(),
		},
	}
}
