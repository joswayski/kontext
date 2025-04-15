#!/bin/bash

# Exit on error
set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Initialize run flags (default is to run all services if no arguments provided)
RUN_API=false
RUN_WEB=false

RUN_ALL=true

show_usage() {
  echo -e "Usage: ./scripts/run.sh [service1 service2 ...]"
  echo -e "Run without arguments to start all services.\n"
  echo -e "Available services:"
  echo -e "  api       Start API service"
  echo -e "  web       Start Web service"
  # Add more options as they become available
  # echo -e "  migrator  Start migrator/seeder service"
  # echo -e "  consumer  Start consumer service"
  echo -e "\nExamples:"
  echo -e "  ./run.sh           # Run all services"
  echo -e "  ./run.sh api       # Run only API service"
  echo -e "  ./run.sh api web   # Run API and web services"
  exit 1
}

# Check for help flag
if [[ "$1" == "--help" ]] || [[ "$1" == "-h" ]]; then
  show_usage
fi

# If arguments were provided, set specific services to run
if [ $# -gt 0 ]; then
  RUN_ALL=false
  
  # Process each argument as a service name
  for service in "$@"; do
    case "$service" in
      api)
        RUN_API=true
        ;;
      web)
        RUN_WEB=true
        ;;

      *)
        echo -e "${YELLOW}Unknown service: $service${NC}"
        show_usage
        ;;
    esac
  done
fi

echo -e "\n\n${BLUE}🚀 Starting Kontext services...${NC}"

# Start API service if selected
API_PID=""
if [ "$RUN_API" = true ] || [ "$RUN_ALL" = true ]; then
    echo -e "\n\n${GREEN}🌐 Starting API service...${NC}"
    cd services/api
    bacon run --headless &
    API_PID=$!
    cd ../..
fi

# Start Web service if selected
WEB_PID=""
if [ "$RUN_WEB" = true ] || [ "$RUN_ALL" = true ]; then
    echo -e "\n\n${GREEN}🖥️ Starting Web service...${NC}"
    cd services/web
    npm run dev &
    WEB_PID=$!
    cd ../..
fi


echo -e "\n${BLUE}=======================================================================${NC}"
echo -e "${YELLOW}🚀 Services are now running!${NC}"
echo -e "${YELLOW}Press ${GREEN}Ctrl+C${YELLOW} to stop all services${NC}\n"
echo -e "${BLUE}=======================================================================${NC}\n\n"

# Function to handle cleanup - kill only the services that were started
cleanup() {
    echo -e "\n\n${BLUE}🛑 Stopping services...${NC}"
    # Only attempt to kill processes that exist
    [ -n "$API_PID" ] && kill $API_PID 2>/dev/null || true
    [ -n "$WEB_PID" ] && kill $WEB_PID 2>/dev/null || true

    exit 0
}

# Trap SIGINT and SIGTERM signals and call cleanup
trap cleanup SIGINT SIGTERM

# Keep the script running
wait
