#!/bin/bash

# Exit on error
set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Starting Kontext services...${NC}"

# Build shared library first
echo -e "${GREEN}📦 Building shared library...${NC}"
cd services/shared
cargo build
cd ../..

# Start API service
echo -e "${GREEN}🌐 Starting API service...${NC}"
cd services/api
cargo run &
API_PID=$!
cd ../..

# Start Web service
echo -e "${GREEN}🖥️ Starting Web service...${NC}"
cd services/web
npm run dev &
WEB_PID=$!
cd ../..

# Function to handle cleanup
cleanup() {
    echo -e "${BLUE}🛑 Stopping services...${NC}"
    kill $API_PID $WEB_PID
    exit 0
}

# Trap SIGINT and SIGTERM signals and call cleanup
trap cleanup SIGINT SIGTERM

echo -e "${GREEN}✅ All services are running!${NC}"
echo -e "${BLUE}Press Ctrl+C to stop all services${NC}"

# Keep the script running
wait
