package messages

import "github.com/brianvoe/gofakeit/v6"

func GenerateRideCancelledMessage() interface{} {
	return map[string]interface{}{
		"ride_id":      gofakeit.UUID(),
		"driver_id":    gofakeit.UUID(),
		"rider_id":     gofakeit.UUID(),
		"cancelled_at": gofakeit.Date(),
		"cancelled_by": gofakeit.RandomString([]string{"rider", "driver", "system"}),
		"reason":       gofakeit.RandomString([]string{"no_driver_available", "rider_cancelled", "driver_cancelled", "payment_failed"}),
	}
}
