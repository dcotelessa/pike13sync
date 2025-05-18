# GitHub Actions Deployment Guide for Pike13Sync

This guide walks you through deploying Pike13Sync using GitHub Actions for automatic synchronization.

## Overview

Using GitHub Actions for Pike13Sync provides a serverless solution that requires:
- No dedicated server to maintain
- Automatic scheduled runs
- Built-in monitoring and notifications
- Secure credential storage using GitHub Secrets

## Setup Steps

### 1. Create GitHub Repository

If you haven't already:
1. Create a new GitHub repository
2. Push your Pike13Sync code to the repository

### 2. Configure Secrets

In your repository:
1. Go to Settings > Secrets and variables > Actions
2. Add the following secrets:
   - `GOOGLE_CREDENTIALS`: Your entire Google credentials JSON content
   - `PIKE13_CLIENT_ID`: Your Pike13 Client ID
   - `CALENDAR_ID`: Your Google Calendar ID
   - `PIKE13_URL`: (Optional) Custom Pike13 API endpoint if different from default

### 3. Enable Workflows

1. Make sure the workflow files are in the `.github/workflows/` directory:
   - `preflight.yml`: Testing connections
   - `sync.yml`: Actual synchronization process
2. Go to the Actions tab and enable workflows if needed

### 4. First Run

1. Run the preflight test first:
   - Go to Actions > Pike13Sync Preflight Test
   - Click "Run workflow"
   - Confirm both Google Calendar API and Pike13 API connections are successful
   
2. Run the sync with dry-run:
   - Go to Actions > Pike13 to Google Calendar Sync
   - Click "Run workflow"
   - Set dry_run to "true"
   - Verify events are detected correctly
   
3. Run a full sync:
   - Run again with dry_run set to "false"
   - Verify events are created in your Google Calendar

## Terraform Deployment (Optional)

If you want to automate the setup of GitHub Actions workflows and repository secrets, you can use Terraform:

1. Navigate to the `deployment/github_actions` directory
2. Create a `terraform.tfvars` file with your configuration:
   ```hcl
   github_token         = "your_github_token"
   github_owner         = "your_github_username"
   repository_name      = "pike13sync"
   # Add other required values
   ```
3. Run Terraform:
   ```bash
   terraform init
   terraform apply
   ```

This will set up the necessary workflow files and repository secrets automatically.

## Log Management

GitHub Actions automatically manages logs for each workflow run:

1. **Retention policy**: GitHub keeps workflow logs for 90 days by default
2. **Artifacts**: Log files from Pike13Sync are saved as artifacts for 14 days
3. **Downloading logs**: Access logs from the Actions tab > specific run > Artifacts section

## Customizing the Schedule

To change when syncs run automatically:

1. Edit the `sync.yml` file in `.github/workflows/`
2. Modify the `cron` expression in the schedule section:

```yaml
schedule:
  - cron: '0 2 * * *'  # Currently set to run at 2 AM UTC daily
```

Common cron patterns:
- `0 */6 * * *`: Every 6 hours
- `0 8,20 * * *`: At 8 AM and 8 PM
- `0 10 * * 1-5`: At 10 AM on weekdays
- `0 9 * * 0,6`: At 9 AM on weekends

## Notifications

### Email Notifications

To enable email notifications for failed runs:

1. Add these secrets:
   - `MAIL_SERVER`: SMTP server address
   - `MAIL_PORT`: SMTP port (usually 587)
   - `MAIL_USERNAME`: SMTP username
   - `MAIL_PASSWORD`: SMTP password
   - `NOTIFICATION_EMAIL`: Recipient email address

2. Uncomment the email notification section in the workflow file

### Slack Notifications

To enable Slack notifications:

1. Create a Slack webhook
2. Add the `SLACK_WEBHOOK` secret
3. Uncomment the Slack notification section in the workflow file

## Troubleshooting

### Common Issues

1. **Authentication errors**:
   - Verify your Google credentials JSON is correctly formatted
   - Ensure the service account has access to the calendar
   
2. **Pike13 API errors**:
   - Check your Pike13 client ID is correct
   - Verify network connectivity to Pike13
   
3. **Rate limiting**:
   - If runs fail due to API rate limits, adjust the schedule
   
4. **Permission issues**:
   - Ensure the service account has proper permissions on the calendar

### Checking Logs

1. Go to the Actions tab in your repository
2. Click on the specific workflow run
3. View logs for each step
4. Download artifacts for detailed application logs

## Uninstalling

To stop Pike13Sync from running on GitHub Actions:

1. Disable or delete the workflow files
2. Remove the repository secrets
3. Optionally, delete any Pike13Sync events from your calendar

## Need Help?

If you encounter issues:
1. Check the logs in GitHub Actions for error messages
2. Consult the troubleshooting section in the main README
3. Create an issue in the repository if needed
