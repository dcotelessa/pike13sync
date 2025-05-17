# Pike13Sync Test Suite

This document outlines the comprehensive test suite for Pike13Sync, explaining the testing strategy and how to run the tests.

## Test Structure

The Pike13Sync test suite is organized into several levels:

1. **Unit Tests**: Test individual functions and components in isolation
2. **Integration Tests**: Test multiple components working together
3. **Connection Tests**: Verify connectivity to Pike13 and Google Calendar APIs
4. **End-to-End Tests**: Test the entire application workflow

## Running Tests

### Using the Test Script

The easiest way to run all tests is to use the provided test script:

```bash
# Make the script executable if needed
chmod +x run_tests.sh

# Run all tests
./run_tests.sh
```

This script will:
- Run unit tests for all components
- Run integration tests
- Run connection tests
- Generate a test coverage report

### Running Individual Tests

You can also run tests for specific packages:

```bash
# Run tests for a specific package
go test -v ./internal/util

# Run tests with coverage
go test -v -coverprofile=coverage.out ./internal/util
go tool cover -html=coverage.out -o coverage.html
```

### Running Integration Tests

Integration tests verify that multiple components work together correctly:

```bash
# Run integration tests
go test -v ./tests -run TestPike13ToGoogleCalendarIntegration
```

### Running Connection Tests

Connection tests verify connectivity to the Pike13 and Google Calendar APIs:

```bash
# Run connection tests
go test -v ./tests -run TestConnectionOnly
```

## Test Descriptions

### Unit Tests

| Test Package | Description |
|--------------|-------------|
| `internal/util` | Tests utility functions for environment loading, logging, and time formatting |
| `internal/config` | Tests configuration loading and environment variable handling |
| `internal/pike13` | Tests Pike13 API client functionality |
| `internal/calendar` | Tests Google Calendar service interactions |
| `internal/sync` | Tests synchronization logic between Pike13 and Google Calendar |
| `cmd/pike13sync` | Tests main application functions |

### Integration Tests

| Test | Description |
|------|-------------|
| `TestPike13ToGoogleCalendarIntegration` | Tests the full integration between Pike13 and Google Calendar using a mock Pike13 API server |

### Connection Tests 

| Test | Description |
|------|-------------|
| `TestConnectionOnly` | Tests connectivity to both Pike13 and Google Calendar APIs |
| `run_script_test.go` | Tests that the run.sh script executes correctly |

## Mock Servers

For integration testing, we create a mock Pike13 API server that simulates the Pike13 API responses. This allows us to test integration scenarios without requiring actual API access.

## Coverage Reports

After running the test script, a coverage report will be generated at `./coverage/coverage.html`. This report shows which lines of code are covered by tests and which are not.

## Adding New Tests

When adding new functionality, make sure to add corresponding tests:

1. Add unit tests for new functions
2. Update integration tests if necessary
3. Run the full test suite to ensure everything still works correctly

## Testing Guidelines

1. Aim for at least 80% code coverage
2. Test both success and error conditions
3. Use mocks for external dependencies when appropriate
4. Keep tests focused and fast

## Troubleshooting

If tests fail, check the following:

- Ensure all required environment variables are set
- Check that credential files exist in the correct locations
- Look for error messages in the test output
- If using Docker, ensure the container has the necessary permissions
