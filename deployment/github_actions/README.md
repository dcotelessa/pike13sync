# GitHub Actions Deployment with Terraform/OpenTofu

This directory contains Terraform/OpenTofu configuration for deploying Pike13Sync to GitHub Actions.

## Overview

This configuration:
1. Creates or uses an existing GitHub repository
2. Sets up GitHub repository secrets for credentials
3. Creates GitHub Actions workflow files
4. Creates README files with documentation

## Prerequisites

1. Terraform or OpenTofu installed on your local machine
2. GitHub personal access token with repo and workflow permissions
3. Pike13 Client ID
4. Google Calendar service account credentials
5. Google Calendar ID

## Setup

1. Copy the example variables file:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   ```

2. Edit `terraform.tfvars` with your specific values:
   - GitHub token and owner
   - Repository settings
   - Pike13 and Google credentials
   - Notification settings (optional)

3. Initialize Terraform/OpenTofu:
   ```bash
   terraform init
   # or
   tofu init
   ```

4. Preview the changes:
   ```bash
   terraform plan
   # or
   tofu plan
   ```

5. Apply the configuration:
   ```bash
   terraform apply
   # or
   tofu apply
   ```

6. Push your Pike13Sync code to the repository

## Customization

### Scheduling

The default schedule is to run daily at 2 AM UTC. You can customize this by changing the `sync_schedule` variable in your `terraform.tfvars` file:

```hcl
sync_schedule = "0 */6 * * *"  # Every 6 hours
```

### Notifications

The configuration supports both email and Slack notifications. To enable them, update the corresponding variables in your `terraform.tfvars` file:

```hcl
# Email Notifications
enable_email_notifications = true
mail_server          = "smtp.example.com"
mail_port            = "587"
mail_username        = "your_email@example.com"
mail_password        = "your_email_app_password"
notification_email   = "recipient@example.com"

# Slack Notifications
enable_slack_notifications = true
slack_webhook        = "https://hooks.slack.com/services/your/webhook/url"
slack_channel        = "pike13sync-status"
```

## Key Differences from Local Development

When using GitHub Actions:

1. **Credentials as Secrets**:
   - Local development uses `.env` files and local credential files
   - GitHub Actions uses repository secrets (more secure)

2. **No Local Files**:
   - No need to store credentials in files
   - Configuration is managed through GitHub secrets
   - `.env` file is only created during workflow execution and not committed

3. **Automated Scheduling**:
   - Local development relies on cron jobs or manual execution
   - GitHub Actions has built-in scheduling

4. **Logging**:
   - Local logs are stored in files
   - GitHub Actions logs are available in the GitHub UI and as artifacts

## Uninstalling

To remove the GitHub Actions deployment:

```bash
terraform destroy
# or
tofu destroy
```

This will remove the workflow files and secrets, but will not delete the repository itself unless you set `prevent_destroy = false` in the repository resource.

## Terraform Configuration Files

### main.tf

The main configuration file that:
- Sets up the GitHub provider
- Creates or references the GitHub repository
- Configures repository secrets
- Creates workflow files from templates

### variables.tf

Defines all the input variables for the configuration, including:
- GitHub access information
- Repository settings
- Credentials for Pike13 and Google Calendar
- Notification settings

### outputs.tf

Defines the outputs after applying the configuration:
- Repository URL
- Actions URL
- Setup completion message

### terraform.tfvars.example

An example file showing how to configure the variables:
- GitHub token and owner settings
- Repository configuration
- Pike13 and Google Calendar credentials
- Notification settings

## Templates

The templates directory contains:
- `preflight.yml.tpl`: Template for the preflight test workflow
- `sync.yml.tpl`: Template for the main sync workflow
- `README.md.tpl`: Template for the repository README

## Security Considerations

1. **GitHub Secrets**: All sensitive information is stored as GitHub Secrets, which are encrypted and only exposed to the GitHub Actions workflow during runtime.

2. **Personal Access Token**: The GitHub token used for Terraform has access to create/update repositories and secrets. Store this token securely and consider using a short-lived token.

3. **Notifications**: If using email or Slack notifications, be aware that the notifications could potentially contain information about your calendar events.

## Troubleshooting

1. **Permission Issues**: Ensure your GitHub token has the correct permissions (repo, workflow).

2. **Secret Size Limitations**: GitHub Secrets have a size limitation. If your Google credentials JSON is very large, ensure it doesn't exceed GitHub's limits.

3. **Rate Limiting**: Be aware of GitHub API rate limits if you're applying changes frequently.

4. **Workflow Debugging**: If workflows fail, check the GitHub Actions logs in the repository's Actions tab.

## Support

If you encounter issues with the Terraform/OpenTofu configuration:

1. Check the Terraform logs for detailed error messages
2. Verify your variables are correctly set
3. Ensure your GitHub token has not expired
4. Reference the GitHub provider documentation for specific errors

## Example Usage

```bash
# Initialize the configuration
terraform init

# Plan the changes
terraform plan -out=plan.out

# Apply the changes
terraform apply plan.out

# To remove everything
terraform destroy
```

This provides a complete infrastructure-as-code solution for deploying Pike13Sync to GitHub Actions.
