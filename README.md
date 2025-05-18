# Pike13Sync: Fitness Studio Calendar Synchronization

Pike13Sync is a Go application that synchronizes class schedules from the Pike13 studio management system to Google Calendar. This tool helps fitness studio owners and administrators keep their Google Calendar up-to-date with class schedules automatically.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Setup Options](#setup-options)
  - [GitHub Actions Deployment](#github-actions-deployment)
  - [Local Setup](#local-setup)
- [Configuration](#configuration)
- [Automation](#automation)
- [Log Management](#log-management)
- [Usage](#usage)
  - [Command Line Options](#command-line-options)
  - [Example Commands](#example-commands)
- [Monitoring and Notifications](#monitoring-and-notifications)
- [Troubleshooting](#troubleshooting)
- [Uninstalling](#uninstalling)
- [Contributing](#contributing)
- [License](#license)

## Overview

Pike13Sync connects to your Pike13 account API, fetches upcoming class schedule data, and synchronizes it to a specified Google Calendar. It can be run manually or automatically via GitHub Actions.

## Features

- Synchronize Pike13 class schedules to Google Calendar
- Support for class details including instructors, capacity, and waitlist status
- Color-coding for class status (active vs. canceled)
- Dry run mode for testing without making changes
- GitHub Actions integration for serverless operation
- Detailed logging for troubleshooting
- Customizable schedules for automatic synchronization

## Prerequisites

- GitHub account (for GitHub Actions deployment)
- Pike13 account with API access
- Google Cloud account with Calendar API enabled
- Service account credentials for Google Calendar API

## Setup Options

### GitHub Actions Deployment

The recommended deployment method is using GitHub Actions:

1. Fork or clone this repository to your GitHub account
2. Set up required secrets in your repository:
   - `GOOGLE_CREDENTIALS`: Your Google Calendar service account credentials JSON
   - `PIKE13_CLIENT_ID`: Your Pike13 API client ID
   - `CALENDAR_ID`: Your Google Calendar ID

3. Enable workflows in the Actions tab of your repository
4. Run the preflight test workflow to verify connectivity
5. Schedule automatic syncs or run them manually as needed

See [deployment/README.md](deployment/README.md) for detailed GitHub Actions setup instructions.

### Local Setup

For local development or manual operation:

1. Clone the repository
2. Set up Google Calendar API:
   - Create a service account with Calendar API access
   - Download credentials JSON to `./credentials/credentials.json`
   - Share your Google Calendar with the service account

3. Create a `.env` file with your configuration:
   ```
   CALENDAR_ID=your_calendar_id@group.calendar.google.com
   PIKE13_CLIENT_ID=your_pike13_client_id
   GOOGLE_CREDENTIALS_FILE=./credentials/credentials.json
   TZ=America/Los_Angeles
   ```

4. Run with:
   ```bash
   ./run.sh
   ```

## Configuration

Pike13Sync can be configured in several ways, in order of precedence:

1. Command-line flags (highest priority)
2. Environment variables (via `.env` file or set directly)
3. Default values (lowest priority)

### Environment Variables

| Variable | Description | Default Value |
|----------|-------------|---------------|
| `CALENDAR_ID` | Google Calendar ID | "primary" |
| `GOOGLE_CREDENTIALS_FILE` | Path to Google API credentials | "./credentials/credentials.json" |
| `PIKE13_CLIENT_ID` | Pike13 API client ID | (required, no default) |
| `PIKE13_URL` | Pike13 API endpoint URL | "https://herosjourneyfitness.pike13.com/api/v2/front/event_occurrences.json" |
| `TZ` | Time zone for calendar events | "America/Los_Angeles" |
| `LOG_PATH` | Path to log file | "./logs/pike13sync.log" |
| `DRY_RUN` | Whether to run without making changes | "false" |

## Automation

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

You can modify the schedule in the `.github/workflows/sync.yml` file.

### Setting up a Cron Job (Local)

For local deployment, you can set up a cron job:

```bash
# Edit your crontab
crontab -e

# Add a job to run Pike13Sync every hour
0 * * * * cd /path/to/pike13sync && ./run.sh >> ./logs/cron.log 2>&1
```

## Log Management

### GitHub Actions Log Management

When running on GitHub Actions, logs are automatically captured and available in the Actions tab:

1. Go to your repository's Actions tab
2. Click on the specific workflow run
3. View logs for each step

Logs are also saved as artifacts for 14 days and can be downloaded for detailed analysis.

### Local Log Management

Pike13Sync creates logs in the `./logs` directory:

- `pike13sync.log`: Main application log
- `pike13_response.json`: Raw response from Pike13 API (for debugging)
- `cron_run.log`: Output from scheduled runs

To increase logging verbosity, use the `--debug` flag:

```bash
./run.sh --debug
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

## Monitoring and Notifications

### Email Notifications for GitHub Actions

To enable email notifications, add these secrets to your repository:
- `MAIL_SERVER`: SMTP server address
- `MAIL_PORT`: SMTP server port
- `MAIL_USERNAME`: SMTP username
- `MAIL_PASSWORD`: SMTP password
- `NOTIFICATION_EMAIL`: Recipient email address

Then uncomment the email notification section in your workflow file.

### Slack Notifications

For Slack notifications, add the `SLACK_WEBHOOK` secret and uncomment the Slack notification section in your workflow file.

## Troubleshooting

### Common Issues

1. **Authentication Errors**:
   - Verify your Google credentials JSON file is correct
   - Check that the service account has proper permissions on the calendar

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

### Testing Connectivity

Test Google Calendar API connectivity:

```bash
go run tests/test_calendar.go --check-write
```

## Uninstalling

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
