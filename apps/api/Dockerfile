FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .


FROM scratch

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 4000

CMD ["./main"]
