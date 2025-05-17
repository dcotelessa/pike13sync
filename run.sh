#!/bin/bash
# run.sh - Script to run Pike13Sync

# Handle script arguments as env file
if [ -n "$1" ] && [ "$1" == "--env-file" ]; then
  if [ -n "$2" ]; then
    export ENV_FILE="$2"
    shift 2
  else
    echo "Error: --env-file requires a file path"
    exit 1
  fi
fi

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Change to the project root directory
cd "$SCRIPT_DIR"

# Check for .env file in the current directory
if [ -f ".env" ] && [ -z "$ENV_FILE" ]; then
  echo "Using .env file from project root"
  export ENV_FILE=".env"
fi

# Pass through all script arguments to the application
go run cmd/pike13sync/main.go "$@"
