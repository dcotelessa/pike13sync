output "deployment_path" {
  description = "Path where Pike13Sync is deployed on the NAS"
  value       = var.nas_app_path
}

output "cron_schedule" {
  description = "Cron schedule for Pike13Sync"
  value       = var.sync_schedule
}

output "deployment_status" {
  description = "Deployment status message"
  value       = "Pike13Sync has been deployed to ${var.nas_hostname}:${var.nas_app_path}"
}

output "next_steps" {
  description = "Next steps to complete setup"
  value       = <<-EOT
    Pike13Sync has been successfully deployed!
    
    Next steps:
    1. SSH into your NAS: ssh ${var.nas_username}@${var.nas_hostname} -p ${var.nas_ssh_port}
    2. Navigate to the deployment directory: cd ${var.nas_app_path}
    3. Pull the Docker image: docker-compose pull
    4. Start the service for a one-time run: docker-compose up
    5. Check the logs directory for any issues: ls -la logs/
    
    The scheduled sync is configured to run: ${var.sync_schedule}
    
    To uninstall:
    1. Run 'terraform destroy' from this directory
    2. SSH into your NAS and remove the directory: rm -rf ${var.nas_app_path}
  EOT
}
