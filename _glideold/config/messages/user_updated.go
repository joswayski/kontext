package messages

import "github.com/brianvoe/gofakeit"

func GenerateUserUpdatedMessage() interface{} {
	return map[string]interface{}{
		"user_id":    gofakeit.UUID(),
		"email":      gofakeit.Email(),
		"first_name": gofakeit.FirstName(),
		"last_name":  gofakeit.LastName(),
		"phone":      gofakeit.Phone(),
		"updated_at": gofakeit.Date(),
		"status":     "active",
	}
}
