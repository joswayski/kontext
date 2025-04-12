#!/bin/bash

# Exit on error
set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Initialize run flags (default is to run all tests if no arguments provided)
RUN_API=false
RUN_WEB=false
RUN_SHARED=false

RUN_ALL=true

show_usage() {
  echo -e "Usage: ./scripts/test.sh [service1 service2 ...]"
  echo -e "Run without arguments to test all services."
  echo -e "\nAvailable services:"
  echo -e "  api       Run API tests"
  echo -e "  web       Run web tests"
  echo -e "  shared    Run shared library tests"

  echo -e "\nExamples:"
  echo -e "  ./test.sh           # Run all tests"
  echo -e "  ./test.sh api       # Run only API tests"
  echo -e "  ./test.sh api web   # Run API and web tests"
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
      shared)
        RUN_SHARED=true
        ;;

      *)
        echo -e "${YELLOW}Unknown service: $service${NC}"
        show_usage
        ;;
    esac
  done
fi

echo -e "\n\n${BLUE}🧪 Running Kontext tests...${NC}"

# Initialize test status variables
API_TEST_STATUS=0
WEB_TEST_STATUS=0
SHARED_TEST_STATUS=0


# Run Shared tests if specified or if running all tests
if [ "$RUN_SHARED" = true ] || [ "$RUN_ALL" = true ]; then
    if [ -d "services/shared" ]; then
        echo -e "\n\n${GREEN}📚 Testing Shared library...${NC}"
        cd services/shared
        cargo test
        SHARED_TEST_STATUS=$?
        cd ../..
    else
        echo -e "\n\n${YELLOW}ℹ️ Shared library not found, skipping tests.${NC}"
    fi
else
    echo -e "\n\n${YELLOW}ℹ️ Shared tests skipped.${NC}"
fi

# Run API tests if specified or if running all tests
if [ "$RUN_API" = true ] || [ "$RUN_ALL" = true ]; then
    echo -e "\n\n${GREEN}🌐 Testing API service...${NC}"
    cd services/api
    cargo test
    API_TEST_STATUS=$?
    cd ../..
else
    echo -e "\n\n${YELLOW}ℹ️ API tests skipped.${NC}"
fi

# Run Web tests if specified or if running all tests
if [ "$RUN_WEB" = true ] || [ "$RUN_ALL" = true ]; then
    echo -e "\n\n${GREEN}🖥️ Testing Web service...${NC}"
    cd services/web
    npm run test
    WEB_TEST_STATUS=$?
    cd ../..
else
    echo -e "\n\n${YELLOW}ℹ️ Web tests skipped.${NC}"
fi


echo -e "\n${BLUE}=======================================================================${NC}"

# Check all test statuses
if [ $API_TEST_STATUS -eq 0 ] && [ $WEB_TEST_STATUS -eq 0 ] && [ $SHARED_TEST_STATUS -eq 0 ]; then
    echo -e "${GREEN}✅ All tests passed!${NC}"
    exit 0
else
    echo -e "${YELLOW}❌ Some tests failed!${NC}"
    exit 1
fi
