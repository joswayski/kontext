package messages

import "github.com/brianvoe/gofakeit/v6"

func GenerateDriverDeactivatedMessage() interface{} {
	return map[string]interface{}{
		"driver_id":      gofakeit.UUID(),
		"deactivated_at": gofakeit.Date(),
		"status":         "inactive",
		"reason":         gofakeit.RandomString([]string{"offline", "break", "end_of_shift"}),
	}
}
