module github.com/joswayski/kontext/pkg/kafka

go 1.24.4

require (
	github.com/joswayski/kontext/pkg/config v0.0.0-00010101000000-000000000000
	github.com/twmb/franz-go v1.19.5
	github.com/twmb/franz-go/pkg/kadm v1.16.0
)

require (
	github.com/brianvoe/gofakeit v3.18.0+incompatible // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/twmb/franz-go/pkg/kmsg v1.11.2 // indirect
	golang.org/x/crypto v0.38.0 // indirect
)

replace github.com/joswayski/kontext/pkg/config => ../config
