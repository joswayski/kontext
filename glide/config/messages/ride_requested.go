package messages

import "github.com/brianvoe/gofakeit/v6"

func GenerateRideRequestedMessage() interface{} {
	return map[string]interface{}{
		"ride_id":  gofakeit.UUID(),
		"rider_id": gofakeit.UUID(),
		"pickup": map[string]interface{}{
			"latitude":  gofakeit.Latitude(),
			"longitude": gofakeit.Longitude(),
			"address":   gofakeit.Address().Address,
		},
		"destination": map[string]interface{}{
			"latitude":  gofakeit.Latitude(),
			"longitude": gofakeit.Longitude(),
			"address":   gofakeit.Address().Address,
		},
		"requested_at": gofakeit.Date(),
		"ride_type":    gofakeit.RandomString([]string{"standard", "premium", "pool"}),
	}
}
