# Environment Configuration for Pike13Sync

Pike13Sync uses environment variables as the primary method for configuration. This document explains how to set up your environment for different deployment scenarios.

## Configuration Methods

Pike13Sync can be configured in several ways, in order of precedence:

1. Command-line flags (highest priority)
2. Environment variables
3. Configuration file (`config/config.json`)
4. Default values (lowest priority)

## Environment Variables

The following environment variables are supported:

| Variable | Description | Default Value |
|----------|-------------|---------------|
| `CALENDAR_ID` | Google Calendar ID | "primary" |
| `GOOGLE_CREDENTIALS_FILE` | Path to Google API credentials | "./credentials/credentials.json" |
| `PIKE13_CLIENT_ID` | Pike13 API client ID | (required, no default) |
| `PIKE13_URL` | Pike13 API endpoint URL | "https://herosjourneyfitness.pike13.com/api/v2/front/event_occurrences.json" |
| `TZ` | Time zone for calendar events | "America/Los_Angeles" |
| `LOG_PATH` | Path to log file | "./logs/pike13sync.log" |
| `DRY_RUN` | Whether to run without making changes | "false" |
| `DOCKER_ENV` | Set to "true" when running in Docker | (not set) |

## Using .env Files

For local development and simple deployments, you can use a `.env` file to set environment variables. The application automatically looks for a `.env` file in the project root.

1. Copy the template to create your own `.env` file:
   ```bash
   cp .env.example .env
   ```

2. Edit the `.env` file with your specific settings:
   ```bash
   nano .env
   ```

### Using Custom .env Files

You can specify a custom `.env` file using the `ENV_FILE` environment variable:

```bash
ENV_FILE=/path/to/custom.env ./run.sh
```

Or with the `--env-file` flag in the run script:

```bash
./run.sh --env-file /path/to/custom.env
```

## Docker Environment

When running in Docker, set the `DOCKER_ENV` environment variable to `true` in your Docker Compose file or Dockerfile:

```yaml
environment:
  - DOCKER_ENV=true
  - CALENDAR_ID=your_calendar_id
  - PIKE13_CLIENT_ID=your_pike13_client_id
```

## Configuration Files vs Environment Variables

While Pike13Sync still supports the `config/config.json` file for backward compatibility, the recommended approach is to use environment variables (either directly or via `.env` file) for all configuration.

Benefits of environment variables:
- Better security (credentials aren't stored in version control)
- Easier to change between environments
- Standard approach for containerized applications
- No need to modify files for different deployments

## Testing Configuration

To verify your environment configuration, run:

```bash
./run.sh --show-env
```

This will display all environment variables and the resulting configuration that will be used by the application.
