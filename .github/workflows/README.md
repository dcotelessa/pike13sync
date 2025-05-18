# GitHub Workflows README

This document contains information specific to the GitHub Actions workflows for Pike13Sync.

## Available Workflows

### 1. preflight.yml

**Purpose**: Test connectivity to Pike13 and Google Calendar APIs without making any changes.

**When to use**: Run this workflow before setting up the main sync workflow to verify that your credentials are correct and the APIs are accessible.

**How to run**:
1. Go to the "Actions" tab in the repository
2. Select "Pike13Sync Preflight Test" from the workflows list
3. Click "Run workflow"

### 2. sync.yml

**Purpose**: Run the full Pike13Sync process to synchronize events between Pike13 and Google Calendar.

**When to use**:
- Manually trigger when you need an immediate sync
- Relies on scheduled runs for regular synchronization (default: daily at 2 AM UTC)

**Options**:
- **dry_run**: Set to "true" to run without making actual changes to Google Calendar
- **debug**: Set to "true" to enable verbose logging

**How to run**:
1. Go to the "Actions" tab in the repository
2. Select "Pike13 to Google Calendar Sync" from the workflows list
3. Click "Run workflow"
4. Choose whether to enable dry-run and/or debug mode
5. Click "Run workflow" again

## Customizing Workflows

### Changing the Schedule

To change when the sync workflow runs automatically, edit the `cron` expression in `sync.yml`:

```yaml
on:
  schedule:
    - cron: '0 2 * * *'  # Currently set to run at 2 AM UTC daily
```

Common cron patterns:
- `0 */6 * * *`: Every 6 hours
- `0 8,20 * * *`: At 8 AM and 8 PM
- `0 10 * * 1-5`: At 10 AM on weekdays
- `0 9 * * 0,6`: At 9 AM on weekends

### Adding Notifications

To add email or Slack notifications, see the Monitoring and Notifications section in the main README.

## Troubleshooting Workflow Issues

1. **Check workflow logs**:
   - Click on the failed workflow run
   - Expand the step that failed
   - Look for error messages

2. **Verify secrets**:
   - Go to Settings > Secrets and variables > Actions
   - Ensure all required secrets are set:
     - `GOOGLE_CREDENTIALS`
     - `PIKE13_CLIENT_ID`
     - `CALENDAR_ID`

3. **Check repository permissions**:
   - GitHub Actions requires proper permissions to run workflows
   - Go to Settings > Actions > General
   - Ensure workflow permissions are properly set

