package messages

import "github.com/brianvoe/gofakeit/v6"

func GenerateRefundIssuedMessage() interface{} {
	return map[string]interface{}{
		"refund_id":  gofakeit.UUID(),
		"payment_id": gofakeit.UUID(),
		"ride_id":    gofakeit.UUID(),
		"rider_id":   gofakeit.UUID(),
		"amount":     gofakeit.Float64Range(5.0, 50.0),
		"currency":   "USD",
		"reason":     gofakeit.RandomString([]string{"ride_cancelled", "service_issue", "overcharge", "duplicate_payment"}),
		"issued_at":  gofakeit.Date(),
	}
}
