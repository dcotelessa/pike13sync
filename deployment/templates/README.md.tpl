# Pike13Sync: Fitness Studio Calendar Synchronization

Pike13Sync is a Go application that synchronizes class schedules from the Pike13 studio management system to Google Calendar. This repository is configured for automatic synchronization using GitHub Actions.

## Quick Setup

### Prerequisites

1. Google Calendar API credentials
2. Pike13 Client ID
3. Google Calendar ID

### Configuration

All configuration is managed through GitHub repository secrets:

- `GOOGLE_CREDENTIALS`: Your Google Calendar API service account credentials JSON
- `PIKE13_CLIENT_ID`: Your Pike13 API client ID
- `CALENDAR_ID`: The ID of the Google Calendar to sync with
- `PIKE13_URL`: (Optional) Custom Pike13 API endpoint if needed

## Automatic Sync

Classes are automatically synced according to the schedule:
- Current schedule: `${sync_schedule}` (in UTC)

You can also manually trigger a sync by going to the Actions tab and running the "Pike13 to Google Calendar Sync" workflow.

## Getting Your Credentials

### Google Calendar API

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Google Calendar API for your project
4. Create a service account and download the JSON credentials
5. Share your Google Calendar with the service account email (give it "Make changes to events" permission)

### Pike13 Client ID

If you don't have your Pike13 Client ID:
1. Contact Pike13 support or log into your Pike13 admin dashboard
2. Request/locate your Pike13 API client ID

### Google Calendar ID

Find your Calendar ID in the Google Calendar settings under "Integrate calendar" > "Calendar ID".

## Log Management

GitHub Actions automatically manages logs:
- Workflow run logs are retained for 90 days
- Pike13Sync log files are saved as artifacts for 14 days
- Access logs from the Actions tab > specific run > Artifacts section

## Customizing

### Changing the Schedule

To modify the sync schedule, edit the `sync.yml` file:

```yaml
schedule:
  - cron: "${sync_schedule}"  # Update this pattern
```

### Enabling Notifications

To receive notifications for failed syncs:

1. Add the appropriate secrets for email or Slack notifications
2. Enable the notification sections in the workflow file

## Manual Installation

If you prefer to run Pike13Sync locally:

1. Clone this repository
2. Create a `.env` file with the following content:
   ```
   CALENDAR_ID=your_calendar_id@group.calendar.google.com
   PIKE13_CLIENT_ID=your_pike13_client_id
   GOOGLE_CREDENTIALS_FILE=./credentials/credentials.json
   TZ=America/Los_Angeles
   # Optional: PIKE13_URL=your_custom_pike13_url
   ```
3. Save your Google credentials JSON to `./credentials/credentials.json`
4. Run with `./run.sh` or build with `go build`

## Troubleshooting

If you encounter issues:
1. Check the workflow logs in the Actions tab
2. Look at the uploaded log artifacts for detailed error messages
3. Verify your credentials are correct and have proper permissions
4. For API errors, confirm you can access the Pike13 and Google Calendar APIs

## Project Links

- [Source Code](https://github.com/${github_owner}/${repository_name})
- [Actions](https://github.com/${github_owner}/${repository_name}/actions)

---

## Support & Contributors

If you encounter issues or have questions, please [open an issue](https://github.com/${github_owner}/${repository_name}/issues) in the repository.

Pike13Sync is open source under the MIT License.
