package messages

import "github.com/brianvoe/gofakeit/v6"

func GeneratePaymentMethodRemovedMessage() interface{} {
	return map[string]interface{}{
		"user_id":           gofakeit.UUID(),
		"payment_method_id": gofakeit.UUID(),
		"removed_at":        gofakeit.Date(),
		"reason":            gofakeit.RandomString([]string{"user_requested", "expired", "fraud_detected"}),
	}
}
