package messages

import "github.com/brianvoe/gofakeit"

func GenerateRideStartedMessage() interface{} {
	return map[string]interface{}{
		"ride_id":    gofakeit.UUID(),
		"driver_id":  gofakeit.UUID(),
		"rider_id":   gofakeit.UUID(),
		"started_at": gofakeit.Date(),
		"pickup_location": map[string]interface{}{
			"latitude":  gofakeit.Latitude(),
			"longitude": gofakeit.Longitude(),
		},
	}
}
