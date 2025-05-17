# Configure log rotation on NAS
resource "null_resource" "setup_logrotate" {
  depends_on = [null_resource.deploy_files]

  # Create logrotate config file from template
  provisioner "local-exec" {
    command = "mkdir -p ${path.module}/generated"
  }

  provisioner "local_file" {
    content = templatefile("${path.module}/templates/pike13sync.logrotate.tpl", {
      app_path = var.nas_app_path
    })
    filename        = "${path.module}/generated/pike13sync.logrotate"
    file_permission = "0644"
  }

  # Copy and install logrotate config
  provisioner "file" {
    connection {
      type        = "ssh"
      user        = var.nas_username
      host        = var.nas_hostname
      port        = var.nas_ssh_port
      private_key = file(var.ssh_private_key_path)
    }

    source      = "${path.module}/generated/pike13sync.logrotate"
    destination = "/tmp/pike13sync.logrotate"
  }

  provisioner "remote-exec" {
    connection {
      type        = "ssh"
      user        = var.nas_username
      host        = var.nas_hostname
      port        = var.nas_ssh_port
      private_key = file(var.ssh_private_key_path)
    }

    inline = [
      "sudo mv /tmp/pike13sync.logrotate /etc/logrotate.d/pike13sync",
      "sudo chmod 644 /etc/logrotate.d/pike13sync"
    ]
  }
}# Main Terraform configuration file for Pike13Sync deployment

terraform {
  required_version = ">= 1.0.0"
  required_providers {
    null = {
      source  = "hashicorp/null"
      version = "~> 3.2.0"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 2.4.0"
    }
  }
}

# Local file for environment variables
resource "local_file" "env_file" {
  content = templatefile("${path.module}/templates/.env.tpl", {
    calendar_id         = var.calendar_id
    pike13_client_id    = var.pike13_client_id
    pike13_url          = var.pike13_url
    timezone            = var.timezone
    dry_run             = var.dry_run ? "true" : "false"
    debug_mode          = var.debug_mode ? "true" : "false"
  })
  filename        = "${path.module}/generated/.env"
  file_permission = "0644"
}

# Local file for Google credentials
resource "local_file" "google_credentials" {
  content         = var.google_credentials_content
  filename        = "${path.module}/generated/credentials.json"
  file_permission = "0600"
}

# Docker Compose file
resource "local_file" "docker_compose" {
  content = templatefile("${path.module}/templates/docker-compose.yml.tpl", {
    app_path = var.nas_app_path
  })
  filename        = "${path.module}/generated/docker-compose.yml"
  file_permission = "0644"
}

# Cron job file
resource "local_file" "cron_file" {
  content = templatefile("${path.module}/templates/crontab.tpl", {
    schedule = var.sync_schedule
    app_path = var.nas_app_path
  })
  filename        = "${path.module}/generated/pike13sync.cron"
  file_permission = "0644"
}

# Create directory structure on NAS
resource "null_resource" "create_directories" {
  provisioner "remote-exec" {
    connection {
      type        = "ssh"
      user        = var.nas_username
      host        = var.nas_hostname
      port        = var.nas_ssh_port
      private_key = file(var.ssh_private_key_path)
    }

    inline = [
      "mkdir -p ${var.nas_app_path}/config",
      "mkdir -p ${var.nas_app_path}/credentials",
      "mkdir -p ${var.nas_app_path}/logs",
    ]
  }
}

# Copy files to NAS
resource "null_resource" "deploy_files" {
  depends_on = [null_resource.create_directories, local_file.env_file, local_file.google_credentials, local_file.docker_compose, local_file.cron_file]

  provisioner "file" {
    connection {
      type        = "ssh"
      user        = var.nas_username
      host        = var.nas_hostname
      port        = var.nas_ssh_port
      private_key = file(var.ssh_private_key_path)
    }

    source      = "${path.module}/generated/.env"
    destination = "${var.nas_app_path}/.env"
  }

  provisioner "file" {
    connection {
      type        = "ssh"
      user        = var.nas_username
      host        = var.nas_hostname
      port        = var.nas_ssh_port
      private_key = file(var.ssh_private_key_path)
    }

    source      = "${path.module}/generated/credentials.json"
    destination = "${var.nas_app_path}/credentials/credentials.json"
  }

  provisioner "file" {
    connection {
      type        = "ssh"
      user        = var.nas_username
      host        = var.nas_hostname
      port        = var.nas_ssh_port
      private_key = file(var.ssh_private_key_path)
    }

    source      = "${path.module}/generated/docker-compose.yml"
    destination = "${var.nas_app_path}/docker-compose.yml"
  }
}

# Configure cron job on NAS
resource "null_resource" "setup_cron" {
  depends_on = [null_resource.deploy_files]

  provisioner "file" {
    connection {
      type        = "ssh"
      user        = var.nas_username
      host        = var.nas_hostname
      port        = var.nas_ssh_port
      private_key = file(var.ssh_private_key_path)
    }

    source      = "${path.module}/generated/pike13sync.cron"
    destination = "/tmp/pike13sync.cron"
  }

  provisioner "remote-exec" {
    connection {
      type        = "ssh"
      user        = var.nas_username
      host        = var.nas_hostname
      port        = var.nas_ssh_port
      private_key = file(var.ssh_private_key_path)
    }

    inline = [
      "cat /tmp/pike13sync.cron | crontab -",
      "rm /tmp/pike13sync.cron"
    ]
  }
}
