FROM golang:1.24 AS builder

WORKDIR /app

# Copy go module files for glide and its local dependencies in pkg
COPY glide/go.mod glide/go.sum ./glide/
COPY pkg/ ./pkg/

# Change to glide directory to resolve local dependencies
WORKDIR /app/glide
RUN go mod download

# Change back to app root
WORKDIR /app

# Copy all source code
COPY . .

# Build the glide executable.
# The `main` package is in the `glide` directory
# so we specify that path for the build.
WORKDIR /app/glide
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .


FROM scratch

WORKDIR /app

COPY --from=builder /app/main .

# EXPOSE 3003 # Not needed for this

CMD ["./main"]
