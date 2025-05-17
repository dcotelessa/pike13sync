# Synology NAS Deployment Guide for Pike13Sync

This guide walks you through deploying Pike13Sync on your Synology NAS using Terraform.

## Prerequisites

1. Synology NAS with SSH access enabled
2. Docker package installed on your NAS
3. Terraform installed on your local machine
4. SSH key-based authentication set up for your NAS
5. Google Calendar service account credentials
6. Pike13 client ID

## Setup Steps

### 1. Create Configuration Files

1. Navigate to the `deployment` directory
2. Create a `terraform.tfvars` file using the example:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   ```
3. Edit `terraform.tfvars` with your specific values:
   - Update NAS connection details (hostname, username, SSH key path)
   - Enter your Pike13 client ID
   - Add your Google Calendar ID
   - Set your timezone
   - Paste your full Google service account credentials JSON

Example:
```hcl
nas_hostname = "nas.home.local"
nas_username = "admin"
ssh_private_key_path = "~/.ssh/id_rsa"
pike13_client_id = "your_pike13_client_id"
calendar_id = "your_calendar@group.calendar.google.com"
```

### 2. Initialize Terraform

```bash
cd deployment
terraform init
```

### 3. Verify the Deployment Plan

```bash
terraform plan
```

Review the plan to ensure everything looks correct.

### 4. Deploy to NAS

```bash
terraform apply
```

Type `yes` when prompted to confirm the deployment.

### 5. Verify Installation

1. SSH into your NAS:
   ```bash
   ssh admin@nas.home.local
   ```

2. Navigate to the deployment directory:
   ```bash
   cd /volume1/docker/pike13sync
   ```

3. Check that files are in place:
   ```bash
   ls -la
   ls -la credentials/
   ```

4. Run a one-time sync to test:
   ```bash
   docker-compose up
   ```

5. Check the logs:
   ```bash
   ls -la logs/
   cat logs/pike13sync.log
   ```

## Customizing the Deployment

### Changing the Sync Schedule

Edit the `sync_schedule` variable in `terraform.tfvars` using standard cron format:

```hcl
# Run every 6 hours
sync_schedule = "0 */6 * * *"

# Run at 5:30 AM daily
sync_schedule = "30 5 * * *"
```

### Enabling Dry Run Mode

To test without modifying your calendar:

```hcl
dry_run = true
```

### Adjusting Sync Range

Change how many days in the past and future to sync:

```hcl
sync_days_ahead = 30  # Sync 30 days into the future
sync_days_behind = 14 # Sync 14 days in the past
```

## Troubleshooting

### Checking Logs

```bash
cat /volume1/docker/pike13sync/logs/pike13sync.log
cat /volume1/docker/pike13sync/logs/cron_run.log
```

### Checking Docker Container

```bash
docker ps -a | grep pike13sync
docker logs pike13sync
```

### Common Issues

1. **Authentication failures**: Verify your Google credentials and Pike13 client ID.
2. **Permission issues**: Ensure your NAS user has permissions to create and access the deployment directory.
3. **Docker issues**: Make sure Docker is running on your NAS.
4. **Cron not running**: Check your crontab with `crontab -l` and verify the task is present.

## Uninstalling

To remove Pike13Sync from your NAS:

1. Run Terraform destroy:
   ```bash
   terraform destroy
   ```

2. Manually clean up (if needed):
   ```bash
   ssh admin@nas.home.local
   rm -rf /volume1/docker/pike13sync
   crontab -l | grep -v pike13sync | crontab -
   ```

## Updating

To update your deployment:

1. Edit the variables in `terraform.tfvars` as needed
2. Run `terraform apply` again to apply the changes
