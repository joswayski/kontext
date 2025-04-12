#!/bin/bash

# Exit on error
set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Setting up Kontext...${NC}"

# Build API service
echo -e "\n\n${GREEN}🌐 Building API service...${NC}"
cd services/api
cargo build
cd ../..

# Install Node.js dependencies for web service
echo -e "\n\n${GREEN}🖥️ Installing web service dependencies...${NC}"
cd services/web
npm install
cd ../..

# Shared is not needed, as it is a library and will be rebuilt by the API service

# Start required infrastructure
echo -e "\n\n${GREEN}🐳 Starting infrastructure (Kafka, MySQL, Qdrant)...${NC}"
docker compose up -d

echo -e "\n\n${GREEN}✅ Setup complete!${NC}"
echo -e "\n${BLUE}=======================================================================${NC}"
echo -e "${YELLOW}🚀 NEXT STEP: Run ${GREEN}./scripts/run.sh${YELLOW} to start the services${NC}"
echo -e "${BLUE}=======================================================================${NC}" 
