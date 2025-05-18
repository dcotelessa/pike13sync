#!/bin/bash
# run_tests.sh - A script to run all tests for pike13sync

# Directories to test
DIRECTORIES=(
  "./internal/util"
  "./internal/config"
  "./internal/pike13"
  "./internal/calendar"
  "./internal/sync"
  "./cmd/pike13sync"
)

# Set color outputs
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Running pike13sync tests...${NC}"
echo -e "${YELLOW}=============================${NC}"

# Set up test environment
setup_test_env() {
  # Store original environment variables to restore later
  export ORIGINAL_PIKE13_CLIENT_ID="$PIKE13_CLIENT_ID"
  export ORIGINAL_PIKE13_URL="$PIKE13_URL"
  export ORIGINAL_CALENDAR_ID="$CALENDAR_ID"
  
  # Set test environment variables if not already set
  if [ -z "$PIKE13_CLIENT_ID" ]; then
    export PIKE13_CLIENT_ID="test_client_id"
    echo -e "${YELLOW}Setting test PIKE13_CLIENT_ID environment variable${NC}"
  fi
  
  if [ -z "$PIKE13_URL" ]; then
    export PIKE13_URL="https://test.pike13.com/api/v2/front/event_occurrences.json"
    echo -e "${YELLOW}Setting test PIKE13_URL environment variable${NC}"
  fi
  
  if [ -z "$CALENDAR_ID" ]; then
    export CALENDAR_ID="test_calendar@group.calendar.google.com"
    echo -e "${YELLOW}Setting test CALENDAR_ID environment variable${NC}"
  fi
  
  # Set test mode flag
  export TEST_MODE=true
  
  # Create necessary directories for tests
  mkdir -p ./config
  mkdir -p ./credentials
  mkdir -p ./logs
  
  echo -e "${YELLOW}Test environment set up${NC}"
}

# Restore original environment
restore_env() {
  export PIKE13_CLIENT_ID="$ORIGINAL_PIKE13_CLIENT_ID"
  export PIKE13_URL="$ORIGINAL_PIKE13_URL"
  export CALENDAR_ID="$ORIGINAL_CALENDAR_ID"
  unset TEST_MODE
  unset ORIGINAL_PIKE13_CLIENT_ID
  unset ORIGINAL_PIKE13_URL
  unset ORIGINAL_CALENDAR_ID
  
  echo -e "${YELLOW}Original environment restored${NC}"
}

# Create a coverage profile directory if it doesn't exist
mkdir -p ./coverage

# Function to run tests in a directory
run_tests() {
  local dir=$1
  echo -e "${YELLOW}Testing $dir...${NC}"
  
  # Run tests with coverage
  go test -v -coverprofile=./coverage/coverage_$(basename $dir).out $dir
  
  # Check if tests passed
  if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Tests in $dir passed!${NC}"
    return 0
  else
    echo -e "${RED}✗ Tests in $dir failed!${NC}"
    return 1
  fi
}

# Set up test environment
setup_test_env

# Run all tests
FAILED=0

for dir in ${DIRECTORIES[@]}; do
  run_tests $dir
  if [ $? -ne 0 ]; then
    FAILED=1
  fi
done

# Run integration tests separately
echo -e "${YELLOW}Running tests in ./tests (with consistent package)...${NC}"
go test -v ./tests
if [ $? -ne 0 ]; then
  FAILED=1
  echo -e "${RED}✗ Tests in ./tests failed!${NC}"
else
  echo -e "${GREEN}✓ Tests in ./tests passed!${NC}"
fi

# Restore original environment
restore_env

# Merge coverage profiles
echo -e "${YELLOW}Merging coverage reports...${NC}"
go tool covdata merge -i ./coverage -o ./coverage/merged.out

# Generate HTML coverage report
echo -e "${YELLOW}Generating coverage report...${NC}"
go tool cover -html=./coverage/merged.out -o ./coverage/coverage.html

# Display coverage summary
go tool cover -func=./coverage/merged.out

# Check if any tests failed
if [ $FAILED -eq 1 ]; then
  echo -e "${RED}Some tests failed!${NC}"
  exit 1
else
  echo -e "${GREEN}All tests passed!${NC}"
  echo -e "${YELLOW}HTML coverage report generated at ./coverage/coverage.html${NC}"
  exit 0
fi
