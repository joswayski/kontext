package messages

import "github.com/brianvoe/gofakeit/v6"

func GeneratePaymentFailedMessage() interface{} {
	return map[string]interface{}{
		"payment_id":    gofakeit.UUID(),
		"ride_id":       gofakeit.UUID(),
		"rider_id":      gofakeit.UUID(),
		"amount":        gofakeit.Float64Range(5.0, 50.0),
		"currency":      "USD",
		"failed_at":     gofakeit.Date(),
		"error_code":    gofakeit.RandomString([]string{"insufficient_funds", "card_declined", "expired_card", "fraud_detected"}),
		"error_message": gofakeit.Sentence(5),
	}
}
