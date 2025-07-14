package messages

import "github.com/brianvoe/gofakeit/v6"

func GenerateDriverOnboardedMessage() interface{} {
	return map[string]interface{}{
		"driver_id":     gofakeit.UUID(),
		"vehicle_id":    gofakeit.UUID(),
		"license_plate": gofakeit.Regex(`[A-Z]{2}\d{2}[A-Z]{2}\d{4}`),
		"vehicle_model": gofakeit.Car().Model,
		"vehicle_color": gofakeit.Color(),
		"onboarded_at":  gofakeit.Date(),
		"status":        "pending_verification",
	}
}
