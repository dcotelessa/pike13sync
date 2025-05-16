# Pike13Sync: Fitness Studio Calendar Synchronization

Pike13Sync is a Go application that synchronizes class schedules from the Pike13 studio management system to Google Calendar. This tool helps fitness studio owners and administrators keep their Google Calendar up-to-date with class schedules automatically.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Local Setup](#local-setup)
  - [Setting up Google Calendar API](#setting-up-google-calendar-api)
  - [Setting up Pike13 Access](#setting-up-pike13-access)
  - [Configuration](#configuration)
  - [Running Locally](#running-locally)
- [Deployment](#deployment)
  - [Docker Deployment](#docker-deployment)
  - [Synology NAS Deployment](#synology-nas-deployment)
- [Automation](#automation)
  - [Setting up a Cron Job](#setting-up-a-cron-job)
  - [Docker Scheduling](#docker-scheduling)
- [Usage](#usage)
  - [Command Line Options](#command-line-options)
  - [Example Commands](#example-commands)
- [Troubleshooting](#troubleshooting)
  - [Common Issues](#common-issues)
  - [Logging](#logging)
  - [Testing Connectivity](#testing-connectivity)
- [Uninstalling](#uninstalling)
- [Contributing](#contributing)
- [License](#license)

## Overview

Pike13Sync connects to your Pike13 account API, fetches upcoming class schedule data, and synchronizes it to a specified Google Calendar. It can be run manually, as a scheduled task, or continuously as a Docker container.

## Features

- Synchronize Pike13 class schedules to Google Calendar
- Support for class details including instructors, capacity, and waitlist status
- Color-coding for class status (active vs. canceled)
- Dry run mode for testing without making changes
- Docker support for easy deployment
- Detailed logging for troubleshooting

## Prerequisites

- Go 1.20+ (for local development)
- Docker (for containerized deployment)
- Pike13 account with API access
- Google Cloud account with Calendar API enabled
- Service account credentials for Google Calendar API
- Terraform/OpenTofu (for NAS deployment)

## Local Setup

### Setting up Google Calendar API

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Google Calendar API for your project
4. Create a service account:
   - Navigate to "IAM & Admin" > "Service Accounts"
   - Click "Create Service Account"
   - Enter a name and description
   - Grant this service account access to your project (Project > Editor role)
   - Create a key for the service account (JSON format)
   - Download the JSON key file

5. Share your Google Calendar with the service account:
   - Go to your Google Calendar
   - Click the three dots next to your calendar and select "Settings and sharing"
   - Scroll down to "Share with specific people"
   - Add the service account email (found in the JSON file as `client_email`)
   - Give it "Make changes to events" permission

6. Save the downloaded JSON credentials file to `./credentials/credentials.json`

### Setting up Pike13 Access

1. Contact Pike13 support or log into your Pike13 admin dashboard
2. Request/locate your Pike13 client ID for API access
3. Create a file `./credentials/pike13_credentials.json` with the following content:
   ```json
   {
     "client_id": "YOUR_PIKE13_CLIENT_ID"
   }
   ```
   OR set the `PIKE13_CLIENT_ID` environment variable

### Configuration

1. Create a `.env` file in the project root with the following content:
   ```
   CALENDAR_ID=your_calendar_id@group.calendar.google.com
   PIKE13_CLIENT_ID=your_pike13_client_id
   GOOGLE_CREDENTIALS_FILE=./credentials/credentials.json
   TZ=America/Los_Angeles
   ```

   Replace `your_calendar_id@group.calendar.google.com` with your Google Calendar ID. You can find this in the calendar settings under "Integrate calendar" > "Calendar ID".

2. Ensure you have proper directory structure:
   ```
   pike13sync/
   ├── .env
   ├── credentials/
   │   ├── credentials.json
   │   └── pike13_credentials.json
   ├── config/
   │   └── config.json (optional)
   └── logs/
   ```

3. Test your Google Calendar access:
   ```bash
   go run tests/test_calendar.go --check-write
   ```

### Running Locally

Build and run the application:

```bash
# Using the run script
chmod +x run.sh
./run.sh

# Or manually
go mod tidy
go run cmd/pike13sync/main.go
```

To run in dry-run mode (no actual changes to Google Calendar):

```bash
./run.sh --dry-run
```

## Deployment

### Docker Deployment

1. Build the Docker image:
   ```bash
   docker build -t pike13sync .
   ```

2. Run the container:
   ```bash
   docker run -d \
     --name pike13sync \
     -v $(pwd)/credentials:/app/credentials \
     -v $(pwd)/config:/app/config \
     -v $(pwd)/logs:/app/logs \
     --env-file .env \
     pike13sync
   ```

3. Using docker-compose:
   ```bash
   docker-compose up -d
   ```

### Synology NAS Deployment

You can use Terraform/OpenTofu to deploy Pike13Sync to your Synology NAS:

1. Install Terraform or OpenTofu on your development machine.

2. Set up the deployment environment by filling in the deployment templates:

   - Update `deployment/variables.tf` with your NAS information
   - Update `deployment/terraform.tfvars` with your specific values
   - Configure Docker in DSM on your Synology NAS

3. Initialize and apply the Terraform configuration:
   ```bash
   cd deployment
   terraform init
   terraform plan
   terraform apply
   ```

   The terraform scripts will:
   - Copy the necessary files to your NAS
   - Set up a Docker container on your NAS
   - Configure environment variables
   - Set up a scheduled task for regular execution

## Automation

### Setting up a Cron Job

To run Pike13Sync automatically on a schedule:

1. Edit your crontab:
   ```bash
   crontab -e
   ```

2. Add a job to run Pike13Sync (adjust the path as needed):
   ```
   # Run Pike13Sync every hour
   0 * * * * cd /path/to/pike13sync && ./run.sh >> ./logs/cron.log 2>&1
   ```

For a Synology NAS, you can set up a scheduled task in Control Panel:

1. Open Control Panel > Task Scheduler
2. Create a new Scheduled Task > User-defined script
3. Set your schedule (e.g., hourly)
4. In the Task Settings, enter:
   ```bash
   cd /volume1/docker/pike13sync && docker-compose up
   ```

### Docker Scheduling

For a Docker-based scheduled execution, use Docker's restart policy:

```yaml
# In docker-compose.yml
services:
  pike13sync:
    build: .
    restart: unless-stopped
    # Add this to have the container exit after sync
    command: /app/pike13sync
    # ... other settings ...
```

## Usage

### Command Line Options

Pike13Sync offers several command-line options:

```
--dry-run        Dry run mode - don't actually modify Google Calendar
--from           Test from date (format: 2025-01-01)
--to             Test to date (format: 2025-01-07)
--debug          Enable debug mode with extra logging
--sample         Only fetch and display sample events without syncing
--config         Path to config file
--show-env       Show environment information and exit
```

### Example Commands

```bash
# Run in dry run mode
./run.sh --dry-run

# Sync events for a specific date range
./run.sh --from 2025-05-20 --to 2025-05-27

# Show sample events without syncing
./run.sh --sample

# Debug mode
./run.sh --debug
```

## Troubleshooting

### Common Issues

1. **Authentication Errors**:
   - Verify your Google credentials JSON file is correct
   - Check that the service account has proper permissions on the calendar
   - Run `go run tests/test_calendar.go --check-write` to test connectivity

2. **Pike13 API Errors**:
   - Verify your Pike13 client ID is correct
   - Check internet connectivity
   - Review the Pike13 API status

3. **Calendar Not Updating**:
   - Check that the specified calendar ID is correct
   - Verify the service account has write permissions
   - Check for error messages in the logs

4. **"Application Client ID is required" Error**:
   - Ensure your Pike13 client ID is correctly set in `.env` or credentials file

5. **"Error accessing specified calendar" Error**:
   - Verify your calendar ID is correct
   - Ensure the service account has been granted access to your calendar

### Logging

Pike13Sync creates logs in the `./logs` directory:

- `pike13sync.log`: Main application log
- `pike13_response.json`: Raw response from Pike13 API (for debugging)

To increase logging verbosity, use the `--debug` flag:

```bash
./run.sh --debug
```

### Testing Connectivity

Test Google Calendar API connectivity:

```bash
go run tests/test_calendar.go --check-write
```

This will:
- Verify authentication
- List available calendars
- Create and delete a test event
- Report any permission issues

## Uninstalling

### Docker Cleanup

1. Stop and remove the container:
   ```bash
   docker stop pike13sync
   docker rm pike13sync
   ```

2. With docker-compose:
   ```bash
   docker-compose down
   ```

3. Remove the image:
   ```bash
   docker rmi pike13sync
   ```

### Terraform/NAS Cleanup

1. Run Terraform destroy:
   ```bash
   cd deployment
   terraform destroy
   ```

2. Manually delete any remaining files from your NAS:
   - SSH to your NAS and navigate to the installation directory
   - Remove the application directory: `rm -rf /volume1/docker/pike13sync`

3. If created, remove the scheduled task from Synology Task Scheduler.

### Manual Cleanup

1. Remove the Pike13 events from Google Calendar:
   - You can use the Google Calendar web interface to select and delete events
   - All events created by Pike13Sync have a custom property `pike13_sync: true`

2. Delete the application directory and all its contents.

## Contributing

Contributions to Pike13Sync are welcome! Please feel free to submit a Pull Request.

## License

[MIT License](LICENSE)
