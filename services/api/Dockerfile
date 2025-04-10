# Build stage
FROM rust:1.86.0-slim-bullseye AS builder

# Install build dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    build-essential \
    pkg-config \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

# Create workspace structure
WORKDIR /usr/src
RUN mkdir -p api shared

# Copy shared package first
COPY services/shared ./shared

# Copy API manifests
COPY services/api/Cargo.toml services/api/Cargo.lock ./api/

# Create dummy source files to build dependencies
RUN mkdir -p api/src && \
    echo "fn main() {}" > api/src/main.rs && \
    echo "fn lib() {}" > api/src/lib.rs

# Build dependencies first (this will cache them)
WORKDIR /usr/src/api
RUN cargo build --release

# Remove dummy source files
RUN rm -rf src

# Copy the actual source code
COPY services/api/src ./src

# Build for release
RUN cargo build --release

# Final stage
FROM debian:bullseye-slim

# Install runtime dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy the binary from builder
COPY --from=builder /usr/src/api/target/release/api /usr/local/bin/api

# Set the working directory
WORKDIR /usr/local/bin

# Expose the port
EXPOSE 4000

# Run the binary with RUST_LOG set to see all logs TODO remove
ENV RUST_LOG=debug
CMD ["./api"]
