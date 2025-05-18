output "repository_url" {
  description = "URL of the GitHub repository"
  value       = var.create_repository ? github_repository.pike13sync[0].html_url : data.github_repository.pike13sync[0].html_url
}

output "actions_url" {
  description = "URL to the GitHub Actions tab"
  value       = "${var.create_repository ? github_repository.pike13sync[0].html_url : data.github_repository.pike13sync[0].html_url}/actions"
}

output "preflight_workflow_path" {
  description = "Path to the preflight workflow file"
  value       = github_repository_file.preflight_workflow.file
}

output "sync_workflow_path" {
  description = "Path to the sync workflow file"
  value       = github_repository_file.sync_workflow.file
}

output "setup_complete_message" {
  description = "Setup completion message"
  value       = <<-EOT
    ✅ Pike13Sync GitHub Actions deployment completed!

    Repository: ${var.create_repository ? github_repository.pike13sync[0].html_url : data.github_repository.pike13sync[0].html_url}
    
    Next steps:
    1. Push your Pike13Sync code to the repository
    2. Go to the Actions tab: ${var.create_repository ? github_repository.pike13sync[0].html_url : data.github_repository.pike13sync[0].html_url}/actions
    3. Run the "Pike13Sync Preflight Test" workflow to verify connectivity
    4. If successful, run the "Pike13 to Google Calendar Sync" workflow

    The sync will run automatically according to the schedule: ${var.sync_schedule}
  EOT
}
