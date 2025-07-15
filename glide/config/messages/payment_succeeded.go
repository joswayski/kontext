package messages

import "github.com/brianvoe/gofakeit"

func GeneratePaymentSucceededMessage() interface{} {
	return map[string]interface{}{
		"payment_id":     gofakeit.UUID(),
		"ride_id":        gofakeit.UUID(),
		"rider_id":       gofakeit.UUID(),
		"amount":         gofakeit.Float64Range(5.0, 50.0),
		"currency":       "USD",
		"transaction_id": gofakeit.UUID(),
		"succeeded_at":   gofakeit.Date(),
	}
}
