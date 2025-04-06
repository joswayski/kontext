#!/bin/bash

# Exit on error
set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "\n\n${BLUE}🚀 Starting Kontext services...${NC}"

# Build shared library first
echo -e "\n\n${GREEN}📦 Building shared library...${NC}"
cd services/shared
cargo build
cd ../..

# Start API service
echo -e "\n\n${GREEN}🌐 Starting API service...${NC}"
cd services/api
cargo run &
API_PID=$!
cd ../..

# Start Web service
echo -e "\n\n${GREEN}🖥️ Starting Web service...${NC}"
cd services/web
npm run dev &
WEB_PID=$!
cd ../..

echo -e "\n${BLUE}=======================================================================${NC}"
echo -e "${YELLOW}🚀 Services are now running!${NC}"
echo -e "${YELLOW}Press ${GREEN}Ctrl+C${YELLOW} to stop all services${NC}]\n"
echo -e "${BLUE}=======================================================================${NC}\n\n"

# Function to handle cleanup
cleanup() {
    echo -e "\n\n${BLUE}🛑 Stopping services...${NC}"
    kill $API_PID $WEB_PID
    exit 0
}

# Trap SIGINT and SIGTERM signals and call cleanup
trap cleanup SIGINT SIGTERM

# Keep the script running
wait
