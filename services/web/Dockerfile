# Build stage
FROM node:22-slim AS builder

# Set working directory
WORKDIR /usr/src

# Copy package files
COPY services/web/package*.json ./

# Install dependencies
RUN npm ci

# Copy source files
COPY services/web/ ./

# Build the application
RUN npm run build

# Production stage
FROM node:22-slim

# Set working directory
WORKDIR /usr/src

# Copy built files from builder
COPY --from=builder /usr/src/.output ./.output

# Expose the port
EXPOSE 3000

# Start the application
CMD ["node", ".output/server/index.mjs"]
