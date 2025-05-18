#!/bin/bash
# run.sh - Script to run Pike13Sync

# Parse command line arguments
DRY_RUN=false
DEBUG=false
ENV_FILE=".env"
SHOW_ENV=false
FROM_DATE=""
TO_DATE=""
SAMPLE=false

# Function to show usage
show_usage() {
  echo "Usage: $0 [options]"
  echo "Options:"
  echo "  --dry-run             Run in dry-run mode (no actual changes)"
  echo "  --debug               Enable debug mode with extra logging"
  echo "  --env-file FILE       Use specific .env file"
  echo "  --show-env            Show environment information and exit"
  echo "  --from DATE           Test from date (format: 2025-01-01)"
  echo "  --to DATE             Test to date (format: 2025-01-07)"
  echo "  --sample              Only fetch and display sample events without syncing"
  echo "  --help                Show this help message"
  echo ""
  echo "Examples:"
  echo "  $0 --dry-run                   # Run in dry run mode"
  echo "  $0 --from 2025-05-01 --to 2025-05-07  # Sync specific date range"
  echo "  $0 --sample                    # Show sample events only"
}

# Parse arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --dry-run)
      DRY_RUN=true
      shift
      ;;
    --debug)
      DEBUG=true
      shift
      ;;
    --env-file)
      if [ -n "$2" ]; then
        ENV_FILE="$2"
        shift 2
      else
        echo "Error: --env-file requires a file path"
        exit 1
      fi
      ;;
    --show-env)
      SHOW_ENV=true
      shift
      ;;
    --from)
      if [ -n "$2" ]; then
        FROM_DATE="$2"
        shift 2
      else
        echo "Error: --from requires a date"
        exit 1
      fi
      ;;
    --to)
      if [ -n "$2" ]; then
        TO_DATE="$2"
        shift 2
      else
        echo "Error: --to requires a date"
        exit 1
      fi
      ;;
    --sample)
      SAMPLE=true
      shift
      ;;
    --help)
      show_usage
      exit 0
      ;;
    *)
      echo "Unknown option: $1"
      show_usage
      exit 1
      ;;
  esac
done

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Change to the project root directory
cd "$SCRIPT_DIR" || exit 1

# Check for .env file in the current directory
if [ -f "$ENV_FILE" ]; then
  echo "Using .env file: $ENV_FILE"
  export ENV_FILE
else
  echo "Warning: .env file not found at $ENV_FILE"
  
  # If we're missing the default .env file, check if required env vars are set
  if [ "$ENV_FILE" = ".env" ] && [ -z "$PIKE13_CLIENT_ID" ]; then
    echo "Warning: PIKE13_CLIENT_ID environment variable not set and no .env file found."
    echo "You can create a .env file with the required configuration or set environment variables directly."
  fi
fi

# Build command arguments
ARGS=""

# Add dry-run flag if requested
if [ "$DRY_RUN" = true ]; then
  ARGS="$ARGS --dry-run"
  echo "Running in DRY RUN mode (no actual changes)"
fi

# Add debug flag if requested
if [ "$DEBUG" = true ]; then
  ARGS="$ARGS --debug"
  echo "Debug mode enabled"
fi

# Add show-env flag if requested
if [ "$SHOW_ENV" = true ]; then
  ARGS="$ARGS --show-env"
fi

# Add from date if provided
if [ -n "$FROM_DATE" ]; then
  ARGS="$ARGS --from $FROM_DATE"
fi

# Add to date if provided
if [ -n "$TO_DATE" ]; then
  ARGS="$ARGS --to $TO_DATE"
fi

# Add sample flag if requested
if [ "$SAMPLE" = true ]; then
  ARGS="$ARGS --sample"
  echo "Sample mode enabled (no sync operations)"
fi

# Run the application
echo "Running pike13sync with arguments: $ARGS"
go run cmd/pike13sync/main.go $ARGS
