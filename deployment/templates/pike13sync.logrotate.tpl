# Log rotation configuration for Pike13Sync
${app_path}/logs/pike13sync.log {
    daily
    rotate 7
    compress
    missingok
    notifempty
    create 0644 root root
    postrotate
        # Nothing to restart since the app runs via cron/docker
    endscript
}

${app_path}/logs/cron_run.log {
    daily
    rotate 7
    compress
    missingok
    notifempty
    create 0644 root root
}
