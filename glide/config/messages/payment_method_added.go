package messages

import "github.com/brianvoe/gofakeit/v6"

func GeneratePaymentMethodAddedMessage() interface{} {
	return map[string]interface{}{
		"user_id":           gofakeit.UUID(),
		"payment_method_id": gofakeit.UUID(),
		"card_type":         gofakeit.RandomString([]string{"visa", "mastercard", "amex"}),
		"last_four":         gofakeit.Regex(`\d{4}`),
		"expiry_month":      gofakeit.IntRange(1, 12),
		"expiry_year":       gofakeit.IntRange(2024, 2030),
		"added_at":          gofakeit.Date(),
	}
}
