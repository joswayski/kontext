FROM golang:1.24 AS builder

WORKDIR /app

# Copy go module files for api and its local dependencies in pkg
COPY api/go.mod api/go.sum ./api/
COPY pkg/ ./pkg/

# Change to api directory to resolve local dependencies
WORKDIR /app/api
RUN go mod download

# Change back to app root
WORKDIR /app

# Copy all source code
COPY . .

# Build the api executable.
# The `main` package is in the `api` directory
# so we specify that path for the build.
WORKDIR /app/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .


FROM scratch

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 3001

CMD ["./main"]
