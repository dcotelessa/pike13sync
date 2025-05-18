terraform {
  required_providers {
    github = {
      source  = "integrations/github"
      version = "~> 5.0"
    }
  }
}

provider "github" {
  token = var.github_token  # Personal access token with repo and workflow permissions
  owner = var.github_owner  # Your GitHub username or organization
}

# Create GitHub repository if it doesn't exist
resource "github_repository" "pike13sync" {
  count       = var.create_repository ? 1 : 0
  name        = var.repository_name
  description = "Pike13 to Google Calendar synchronization"
  visibility  = var.repository_visibility
  auto_init   = true

  lifecycle {
    prevent_destroy = true
  }
}

# Get repository information if it already exists
data "github_repository" "pike13sync" {
  count = var.create_repository ? 0 : 1
  name  = var.repository_name
}

locals {
  repository_name = var.create_repository ? github_repository.pike13sync[0].name : data.github_repository.pike13sync[0].name
}

# Add GitHub Secrets
resource "github_actions_secret" "google_credentials" {
  repository       = local.repository_name
  secret_name      = "GOOGLE_CREDENTIALS"
  plaintext_value  = var.google_credentials_json
}

resource "github_actions_secret" "pike13_client_id" {
  repository       = local.repository_name
  secret_name      = "PIKE13_CLIENT_ID"
  plaintext_value  = var.pike13_client_id
}

resource "github_actions_secret" "calendar_id" {
  repository       = local.repository_name
  secret_name      = "CALENDAR_ID"
  plaintext_value  = var.calendar_id
}

# Optional: Email notification secrets
resource "github_actions_secret" "mail_server" {
  count            = var.enable_email_notifications ? 1 : 0
  repository       = local.repository_name
  secret_name      = "MAIL_SERVER"
  plaintext_value  = var.mail_server
}

resource "github_actions_secret" "mail_port" {
  count            = var.enable_email_notifications ? 1 : 0
  repository       = local.repository_name
  secret_name      = "MAIL_PORT"
  plaintext_value  = var.mail_port
}

resource "github_actions_secret" "mail_username" {
  count            = var.enable_email_notifications ? 1 : 0
  repository       = local.repository_name
  secret_name      = "MAIL_USERNAME"
  plaintext_value  = var.mail_username
}

resource "github_actions_secret" "mail_password" {
  count            = var.enable_email_notifications ? 1 : 0
  repository       = local.repository_name
  secret_name      = "MAIL_PASSWORD"
  plaintext_value  = var.mail_password
}

resource "github_actions_secret" "notification_email" {
  count            = var.enable_email_notifications ? 1 : 0
  repository       = local.repository_name
  secret_name      = "NOTIFICATION_EMAIL"
  plaintext_value  = var.notification_email
}

# Optional: Slack notification secrets
resource "github_actions_secret" "slack_webhook" {
  count            = var.enable_slack_notifications ? 1 : 0
  repository       = local.repository_name
  secret_name      = "SLACK_WEBHOOK"
  plaintext_value  = var.slack_webhook
}

# Create GitHub Actions Workflow Files
resource "github_repository_file" "preflight_workflow" {
  repository          = local.repository_name
  branch              = var.repository_branch
  file                = ".github/workflows/preflight.yml"
  content             = templatefile("${path.module}/../templates/preflight.yml.tpl", {
    go_version = var.go_version
  })
  commit_message      = "Add Pike13Sync preflight workflow"
  commit_author       = "Terraform"
  commit_email        = var.commit_email
  overwrite_on_create = true
}

resource "github_repository_file" "sync_workflow" {
  repository          = local.repository_name
  branch              = var.repository_branch
  file                = ".github/workflows/sync.yml"
  content             = templatefile("${path.module}/../templates/sync.yml.tpl", {
    go_version = var.go_version,
    schedule = var.sync_schedule,
    enable_email_notifications = var.enable_email_notifications,
    enable_slack_notifications = var.enable_slack_notifications,
    slack_channel = var.slack_channel
  })
  commit_message      = "Add Pike13Sync main workflow"
  commit_author       = "Terraform"
  commit_email        = var.commit_email
  overwrite_on_create = true
}

# Create README file
resource "github_repository_file" "readme" {
  repository          = local.repository_name
  branch              = var.repository_branch
  file                = "README.md"
  content             = templatefile("${path.module}/../templates/README.md.tpl", {
    repository_name = local.repository_name,
    github_owner = var.github_owner
  })
  commit_message      = "Update Pike13Sync README"
  commit_author       = "Terraform"
  commit_email        = var.commit_email
  overwrite_on_create = true
}

# Create README for GitHub workflows
resource "github_repository_file" "workflow_readme" {
  repository          = local.repository_name
  branch              = var.repository_branch
  file                = ".github/workflows/README.md"
  content             = templatefile("${path.module}/../templates/workflow_README.md.tpl", {
    sync_schedule = var.sync_schedule
  })
  commit_message      = "Add GitHub Workflows documentation"
  commit_author       = "Terraform"
  commit_email        = var.commit_email
  overwrite_on_create = true
}
