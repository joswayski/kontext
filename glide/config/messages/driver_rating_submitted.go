package messages

import "github.com/brianvoe/gofakeit/v6"

func GenerateDriverRatingSubmittedMessage() interface{} {
	return map[string]interface{}{
		"driver_id":    gofakeit.UUID(),
		"rider_id":     gofakeit.UUID(),
		"ride_id":      gofakeit.UUID(),
		"rating":       gofakeit.IntRange(1, 5),
		"comment":      gofakeit.Sentence(10),
		"submitted_at": gofakeit.Date(),
	}
}
