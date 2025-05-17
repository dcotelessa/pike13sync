version: '3'
services:
  pike13sync:
    image: dcotelessa/pike13sync:latest
    container_name: pike13sync
    restart: "no"
    env_file:
      - .env
    environment:
      - DOCKER_ENV=true
    volumes:
      - ${app_path}/config:/app/config
      - ${app_path}/credentials:/app/credentials
      - ${app_path}/logs:/app/logs
