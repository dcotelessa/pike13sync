# Variables for Pike13Sync deployment on Synology NAS

variable "nas_hostname" {
  description = "Hostname or IP address of the Synology NAS"
  type        = string
}

variable "nas_ssh_port" {
  description = "SSH port for Synology NAS (default: 22)"
  type        = number
  default     = 22
}

variable "nas_username" {
  description = "SSH username for Synology NAS"
  type        = string
}

variable "ssh_private_key_path" {
  description = "Path to the SSH private key file for authentication"
  type        = string
}

variable "nas_app_path" {
  description = "Path on the NAS where Pike13Sync will be deployed"
  type        = string
  default     = "/volume1/docker/pike13sync"
}

# Pike13Sync specific variables
variable "pike13_client_id" {
  description = "Client ID for Pike13 API"
  type        = string
  sensitive   = true
}

variable "calendar_id" {
  description = "Google Calendar ID for synchronization"
  type        = string
}

variable "timezone" {
  description = "Timezone for calendar events"
  type        = string
  default     = "America/Los_Angeles"
}

variable "sync_schedule" {
  description = "Cron schedule for synchronization (default: run every hour)"
  type        = string
  default     = "0 * * * *"
}

variable "dry_run" {
  description = "Enable dry run mode (no actual changes to calendar)"
  type        = bool
  default     = false
}

variable "google_credentials_content" {
  description = "Content of the Google credentials JSON file"
  type        = string
  sensitive   = true
}

variable "debug_mode" {
  description = "Enable debug mode with extra logging"
  type        = bool
  default     = false
}

variable "sync_days_ahead" {
  description = "Number of days ahead to sync (default: 14 days)"
  type        = number
  default     = 14
}

variable "sync_days_behind" {
  description = "Number of days in the past to sync (default: 7 days)"
  type        = number
  default     = 7
}

variable "pike13_url" {
  description = "Custom Pike13 API URL (optional)"
  type        = string
  default     = "https://herosjourneyfitness.pike13.com/api/v2/front/event_occurrences.json"
}
