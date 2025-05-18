variable "github_token" {
  description = "GitHub personal access token"
  type        = string
  sensitive   = true
}

variable "github_owner" {
  description = "GitHub owner (username or organization)"
  type        = string
}

variable "repository_name" {
  description = "GitHub repository name"
  type        = string
  default     = "pike13sync"
}

variable "repository_branch" {
  description = "GitHub repository branch"
  type        = string
  default     = "main"
}

variable "create_repository" {
  description = "Whether to create a new repository or use an existing one"
  type        = bool
  default     = false
}

variable "repository_visibility" {
  description = "Repository visibility (public or private)"
  type        = string
  default     = "private"
  validation {
    condition     = contains(["public", "private"], var.repository_visibility)
    error_message = "Repository visibility must be either 'public' or 'private'."
  }
}

variable "google_credentials_json" {
  description = "Google service account credentials JSON"
  type        = string
  sensitive   = true
}

variable "pike13_client_id" {
  description = "Pike13 client ID"
  type        = string
  sensitive   = true
}

variable "calendar_id" {
  description = "Google Calendar ID"
  type        = string
}

variable "go_version" {
  description = "Go version to use"
  type        = string
  default     = "1.21"
}

variable "sync_schedule" {
  description = "Cron schedule for the sync workflow"
  type        = string
  default     = "0 2 * * *"  # Daily at 2 AM UTC
}

variable "commit_email" {
  description = "Email to use for commit author"
  type        = string
  default     = "pike13sync-terraform@example.com"
}

# Email notifications
variable "enable_email_notifications" {
  description = "Enable email notifications"
  type        = bool
  default     = false
}

variable "mail_server" {
  description = "SMTP server address"
  type        = string
  default     = ""
  sensitive   = true
}

variable "mail_port" {
  description = "SMTP server port"
  type        = string
  default     = "587"
}

variable "mail_username" {
  description = "SMTP username"
  type        = string
  default     = ""
  sensitive   = true
}

variable "mail_password" {
  description = "SMTP password"
  type        = string
  default     = ""
  sensitive   = true
}

variable "notification_email" {
  description = "Email address to receive notifications"
  type        = string
  default     = ""
}

# Slack notifications
variable "enable_slack_notifications" {
  description = "Enable Slack notifications"
  type        = bool
  default     = false
}

variable "slack_webhook" {
  description = "Slack webhook URL"
  type        = string
  default     = ""
  sensitive   = true
}

variable "slack_channel" {
  description = "Slack channel name"
  type        = string
  default     = "pike13sync-status"
}
