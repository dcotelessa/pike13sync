#!/bin/bash
# create_configs.sh - Script to create all necessary configuration files

# Determine the project root directory
PROJECT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Function to create a directory if it doesn't exist
create_dir() {
  if [ ! -d "$1" ]; then
    echo "Creating directory: $1"
    mkdir -p "$1"
  else
    echo "Directory already exists: $1"
  fi
}

# Create necessary directories
create_dir "$PROJECT_ROOT/config"
create_dir "$PROJECT_ROOT/credentials"
create_dir "$PROJECT_ROOT/logs"

# Create .env file if it doesn't exist
if [ ! -f "$PROJECT_ROOT/.env" ] && [ -f "$PROJECT_ROOT/.env.example" ]; then
  echo "Creating .env file from .env.example"
  cp "$PROJECT_ROOT/.env.example" "$PROJECT_ROOT/.env"
  echo "Please edit .env file with your actual values"
else
  echo ".env file already exists or .env.example not found"
fi

# Create config.json if it doesn't exist
if [ ! -f "$PROJECT_ROOT/config/config.json" ] && [ -f "$PROJECT_ROOT/config/config.json.example" ]; then
  echo "Creating config/config.json from config.json.example"
  cp "$PROJECT_ROOT/config/config.json.example" "$PROJECT_ROOT/config/config.json"
  echo "Please edit config/config.json with your actual values"
else
  echo "config.json file already exists or config.json.example not found"
fi

# Create pike13_credentials.json if it doesn't exist
if [ ! -f "$PROJECT_ROOT/credentials/pike13_credentials.json" ]; then
  echo "Creating credentials/pike13_credentials.json"
  
  # Read Pike13 client ID from user
  read -p "Enter your Pike13 client ID (or press Enter to skip): " pike13_client_id
  
  if [ -n "$pike13_client_id" ]; then
    # Create the file with user input
    echo "{" > "$PROJECT_ROOT/credentials/pike13_credentials.json"
    echo "  \"client_id\": \"$pike13_client_id\"" >> "$PROJECT_ROOT/credentials/pike13_credentials.json"
    echo "}" >> "$PROJECT_ROOT/credentials/pike13_credentials.json"
    echo "pike13_credentials.json created with provided client ID"
  else
    # Create a template file
    echo "{" > "$PROJECT_ROOT/credentials/pike13_credentials.json"
    echo "  \"client_id\": \"your_pike13_client_id\"" >> "$PROJECT_ROOT/credentials/pike13_credentials.json"
    echo "}" >> "$PROJECT_ROOT/credentials/pike13_credentials.json"
    echo "pike13_credentials.json created with placeholder client ID"
    echo "Please edit credentials/pike13_credentials.json with your actual client ID"
  fi
else
  echo "pike13_credentials.json already exists"
fi

# Check for Google Calendar credentials
if [ ! -f "$PROJECT_ROOT/credentials/credentials.json" ]; then
  echo "WARNING: Google Calendar credentials file not found at credentials/credentials.json"
  echo "Please obtain Google Calendar service account credentials and save them to this location"
fi

echo ""
echo "Configuration setup complete."
echo "Make sure to update the configuration files with your actual values before running pike13sync."
echo "You can use ./run.sh --show-env to check your current configuration."
