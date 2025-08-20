package messages

import "github.com/brianvoe/gofakeit"

func GeneratePaymentInitiatedMessage() interface{} {
	return map[string]interface{}{
		"payment_id":        gofakeit.UUID(),
		"ride_id":           gofakeit.UUID(),
		"rider_id":          gofakeit.UUID(),
		"amount":            gofakeit.Float64Range(5.0, 50.0),
		"currency":          "USD",
		"payment_method_id": gofakeit.UUID(),
		"initiated_at":      gofakeit.Date(),
	}
}
