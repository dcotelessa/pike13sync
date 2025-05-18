# Relevant sections of main.tf that need fixing

# Fix for the preflight workflow resource
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

# Fix for the sync workflow resource
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

# Fix for the README resource - adding sync_schedule parameter
resource "github_repository_file" "readme" {
  repository          = local.repository_name
  branch              = var.repository_branch
  file                = "README.md"
  content             = templatefile("${path.module}/../templates/README.md.tpl", {
    repository_name = local.repository_name,
    github_owner = var.github_owner,
    sync_schedule = var.sync_schedule  # Added this parameter
  })
  commit_message      = "Update Pike13Sync README"
  commit_author       = "Terraform"
  commit_email        = var.commit_email
  overwrite_on_create = true
}

# Fix for the workflow README resource - correcting file path and adding parameter
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
