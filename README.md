# Pike13Sync Documentation

This artifact contains two README files for the Pike13Sync project:

1. `README.md` - The main project README with GitHub Actions integration
2. `.github/workflows/README.md` - Specific documentation for GitHub Actions workflows

---

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
  - [GitHub Actions Deployment](#github-actions-deployment)
- [Automation](#automation)
  - [Setting up a Cron Job](#setting-up-a-cron-job)
  - [Docker Scheduling](#docker-scheduling)
  - [GitHub Actions Scheduling](#github-actions-scheduling)
- [Log Management](#log-management)
- [Usage](#usage)
  - [Command Line Options](#command-line-options)
  - [Example Commands](#example-commands)
- [Monitoring and Notifications](#monitoring-and-notifications)
- [Troubleshooting](#troubleshooting)
  - [Common Issues](#common-issues)
  - [Logging](#logging)
  - [Testing Connectivity](#testing-connectivity)
- [Uninstalling](#uninstalling)
- [Contributing](#contributing)
- [License](#license)

## Overview

Pike13Sync connects to your Pike13 account API, fetches upcoming class schedule data, and synchronizes it to a specified Google Calendar. It can be run manually, as a scheduled task, or automatically via GitHub Actions.

## Features

- Synchronize Pike13 class schedules to Google Calendar
- Support for class details including instructors, capacity, and waitlist status
- Color-coding for class status (active vs. canceled)
- Dry run mode for testing without making changes
- Docker support for easy deployment
- GitHub Actions integration for serverless operation
- Detailed logging for troubleshooting
- Automatic log rotation to prevent filesystem bloat

## Prerequisites

- Go 1.20+ (for local development)
- GitHub account (for GitHub Actions deployment)
- Pike13 account with API access
- Google Cloud account with Calendar API enabled
- Service account credentials for Google Calendar API

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

### GitHub Actions Deployment

Pike13Sync can be deployed as a fully serverless solution using GitHub Actions. This method requires no server management and runs on GitHub's infrastructure.

#### Setting Up GitHub Actions Deployment

1. **Create a GitHub Repository**
   - Create a new repository or use an existing one
   - Push your Pike13Sync code to the repository

2. **Add GitHub Secrets**
   - Go to your repository Settings > Secrets and variables > Actions
   - Add the following repository secrets:
     - `GOOGLE_CREDENTIALS`: Your entire Google credentials.json content
     - `PIKE13_CLIENT_ID`: Your Pike13 Client ID
     - `CALENDAR_ID`: Your Google Calendar ID

3. **Create Workflow Files**
   - Create a directory `.github/workflows/` in your repository
   - Add the following workflow files:
     - `preflight.yml`: For testing connections
     - `sync.yml`: For the actual synchronization process

4. **Preflight Workflow Example**
   ```yaml
   name: Pike13Sync Preflight Test

   on:
     workflow_dispatch:

   jobs:
     preflight:
       runs-on: ubuntu-latest
       steps:
         - name: Check out repository
           uses: actions/checkout@v3

         - name: Set up Go
           uses: actions/setup-go@v4
           with:
             go-version: '1.21'
           
         - name: Create directory structure
           run: |
             mkdir -p config
             mkdir -p credentials
             mkdir -p logs

         - name: Set up credentials
           run: |
             echo '${{ secrets.GOOGLE_CREDENTIALS }}' > credentials/credentials.json
             chmod 600 credentials/credentials.json
             echo '{"client_id": "${{ secrets.PIKE13_CLIENT_ID }}"}' > credentials/pike13_credentials.json

         - name: Test API connections
           run: |
             go mod tidy
             go run cmd/pike13sync/main.go --sample
   ```

5. **Full Sync Workflow Example**
   ```yaml
   name: Pike13 to Google Calendar Sync

   on:
     workflow_dispatch:
       inputs:
         dry_run:
           description: 'Run in dry-run mode (no actual changes)'
           required: false
           default: 'false'
           type: choice
           options:
             - 'true'
             - 'false'
     schedule:
       - cron: '0 2 * * *'  # Run daily at 2 AM UTC

   jobs:
     sync:
       runs-on: ubuntu-latest
       steps:
         - name: Check out repository
           uses: actions/checkout@v3

         - name: Set up Go
           uses: actions/setup-go@v4
           with:
             go-version: '1.21'
           
         - name: Create directory structure
           run: |
             mkdir -p config
             mkdir -p credentials
             mkdir -p logs

         - name: Set up credentials
           run: |
             echo '${{ secrets.GOOGLE_CREDENTIALS }}' > credentials/credentials.json
             chmod 600 credentials/credentials.json
             echo '{"client_id": "${{ secrets.PIKE13_CLIENT_ID }}"}' > credentials/pike13_credentials.json

         - name: Run Pike13Sync
           id: sync
           run: |
             ARGS=""
             if [ "${{ github.event.inputs.dry_run || 'false' }}" == "true" ]; then
               ARGS="$ARGS --dry-run"
             fi
             go run cmd/pike13sync/main.go $ARGS
   ```

6. **Run the Workflows**
   - Go to the "Actions" tab in your repository
   - First, run the preflight test to verify connectivity
   - Then, run the full sync workflow with dry-run mode enabled to test
   - Finally, run the full sync workflow without dry-run to perform the actual sync

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

### GitHub Actions Scheduling

GitHub Actions provides built-in scheduling through the workflow file:

```yaml
on:
  schedule:
    # Run daily at 2 AM UTC
    - cron: '0 2 * * *'
    
    # Alternative: Run every 6 hours
    # - cron: '0 */6 * * *'
    
    # Alternative: Run at 8 AM and 8 PM UTC
    # - cron: '0 8,20 * * *'
```

Considerations for GitHub Actions scheduling:
- Times are in UTC
- Schedule precision is not guaranteed (may delay by up to 15 minutes)
- GitHub will disable scheduled workflows if the repo has no activity for 60 days

## Log Management

### Local and Docker Log Management

Pike13Sync creates logs in the `./logs` directory:

- `pike13sync.log`: Main application log
- `pike13_response.json`: Raw response from Pike13 API (for debugging)
- `cron_run.log`: Output from scheduled runs

To increase logging verbosity, use the `--debug` flag:

```bash
./run.sh --debug
```

### GitHub Actions Log Management

When running on GitHub Actions, logs are automatically captured and available in the Actions tab:

1. Go to your repository's Actions tab
2. Click on the specific workflow run
3. View logs for each step

For long-term log storage, you can add this to your workflow:

```yaml
- name: Upload logs as artifacts
  uses: actions/upload-artifact@v3
  with:
    name: pike13sync-logs
    path: |
      logs/
      sync_output.txt
    retention-days: 14
```

This will save logs as artifacts for 14 days (adjustable).

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

## Monitoring and Notifications

### Email Notifications for GitHub Actions

Add this to your workflow for email notifications:

```yaml
- name: Send email notification on failure
  if: failure()
  uses: dawidd6/action-send-mail@v3
  with:
    server_address: ${{ secrets.MAIL_SERVER }}
    server_port: ${{ secrets.MAIL_PORT }}
    username: ${{ secrets.MAIL_USERNAME }}
    password: ${{ secrets.MAIL_PASSWORD }}
    subject: Pike13Sync Failed
    body: The Pike13Sync job has failed. See details at https://github.com/${{ github.repository }}/actions/runs/${{ github.run_id }}
    to: ${{ secrets.NOTIFICATION_EMAIL }}
    from: Pike13Sync Notifications
```

### Slack Notifications

For Slack notifications:

```yaml
- name: Slack Notification
  uses: rtCamp/action-slack-notify@v2
  env:
    SLACK_WEBHOOK: ${{ secrets.SLACK_WEBHOOK }}
    SLACK_CHANNEL: pike13-sync-status
    SLACK_COLOR: ${{ job.status == 'success' && 'good' || 'danger' }}
    SLACK_TITLE: Pike13Sync Result
    SLACK_MESSAGE: Pike13Sync completed with status: ${{ job.status }}
```

### Job Summary

```yaml
- name: Create job summary
  run: |
    echo "## Pike13Sync Summary" >> $GITHUB_STEP_SUMMARY
    echo "" >> $GITHUB_STEP_SUMMARY
    echo "Events created: ${{ steps.sync.outputs.created }}" >> $GITHUB_STEP_SUMMARY
    echo "Events updated: ${{ steps.sync.outputs.updated }}" >> $GITHUB_STEP_SUMMARY
    echo "Events deleted: ${{ steps.sync.outputs.deleted }}" >> $GITHUB_STEP_SUMMARY
    echo "Events unchanged: ${{ steps.sync.outputs.unchanged }}" >> $GITHUB_STEP_SUMMARY
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

4. **GitHub Actions Issues**:
   - Verify secrets are correctly set up
   - Check workflow logs for errors
   - Make sure Go version is compatible (1.20 or newer)

### Logging

To increase logging verbosity, use the `--debug` flag:

```bash
./run.sh --debug
```

In GitHub Actions, check the detailed logs in the workflow run view.

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

### GitHub Actions Cleanup

1. Go to the repository settings
2. Delete the GitHub repository secrets
3. Delete or disable the workflow files in `.github/workflows/`

### Manual Cleanup

1. Remove the Pike13 events from Google Calendar:
   - You can use the Google Calendar web interface to select and delete events
   - All events created by Pike13Sync have a custom property `pike13_sync: true`

2. Delete the application directory and all its contents.

## Contributing

Contributions to Pike13Sync are welcome! Please feel free to submit a Pull Request.

## License

[MIT License](LICENSE)
