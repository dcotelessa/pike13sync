# Pike13Sync scheduled task
${schedule} cd ${app_path} && /usr/local/bin/docker-compose up > ${app_path}/logs/cron_run.log 2>&1
