name: Pike13Sync Preflight Test

on:
  # Allow manual triggering
  workflow_dispatch:
  # Also run on a schedule for testing (adjust as needed)
  schedule:
    # Run daily at 2 AM UTC
    - cron: '0 2 * * *'

jobs:
  preflight:
    runs-on: ubuntu-latest
    steps:
      - name: Check out repository
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '${go_version}'  # Updated to a newer version for better compatibility
          
      - name: Create directory structure
        run: |
          mkdir -p config
          mkdir -p credentials
          mkdir -p logs

      - name: Set up Google credentials
        run: |
          echo '${{ secrets.GOOGLE_CREDENTIALS }}' > credentials/credentials.json
          chmod 600 credentials/credentials.json

      - name: Create .env file
        run: |
          cat > .env << EOF
          CALENDAR_ID=${{ secrets.CALENDAR_ID }}
          PIKE13_CLIENT_ID=${{ secrets.PIKE13_CLIENT_ID }}
          GOOGLE_CREDENTIALS_FILE=./credentials/credentials.json
          TZ=America/Los_Angeles
          ${{ secrets.PIKE13_URL != '' && format('PIKE13_URL={0}', secrets.PIKE13_URL) || '' }}
          EOF

      - name: Test Google Calendar API Connection with compatible versions
        run: |
          # Create a separate directory for the Google Calendar test
          mkdir -p test_calendar
          cd test_calendar
          
          # Create go.mod with version constraints
          cat > go.mod << EOF
          module test_calendar

          go 1.21

          require (
            golang.org/x/oauth2 v0.13.0
            google.golang.org/api v0.149.0
          )
          EOF
          
          # Create a simple test script to check Google Calendar API connectivity
          cat > test_calendar.go << EOF
          package main

          import (
            "context"
            "fmt"
            "os"
            "time"

            "golang.org/x/oauth2/google"
            "google.golang.org/api/calendar/v3"
            "google.golang.org/api/option"
          )

          func main() {
            ctx := context.Background()
            
            # Read credentials from file
            credBytes, err := os.ReadFile("../credentials/credentials.json")
            if err != nil {
              fmt.Printf("Error reading credentials file: %v\n", err)
              os.Exit(1)
            }
            
            # Configure JWT
            config, err := google.JWTConfigFromJSON(credBytes, calendar.CalendarScope)
            if err != nil {
              fmt.Printf("Error creating JWT config: %v\n", err)
              os.Exit(1)
            }
            
            # Create client
            client := config.Client(ctx)
            
            # Create calendar service
            srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
            if err != nil {
              fmt.Printf("Error creating calendar service: %v\n", err)
              os.Exit(1)
            }
            
            # Get calendar ID from environment
            calendarID := os.Getenv("CALENDAR_ID")
            if calendarID == "" {
              calendarID = "primary"
            }
            
            # List next 10 events
            t := time.Now().Format(time.RFC3339)
            events, err := srv.Events.List(calendarID).
              MaxResults(10).
              TimeMin(t).
              OrderBy("startTime").
              SingleEvents(true).
              Do()
            if err != nil {
              fmt.Printf("Error retrieving events: %v\n", err)
              os.Exit(1)
            }
            
            # Print success and event count
            fmt.Printf("Successfully connected to Google Calendar API\n")
            fmt.Printf("Found %d upcoming events\n", len(events.Items))
            
            # Print upcoming events
            for i, item := range events.Items {
              date := item.Start.DateTime
              if date == "" {
                date = item.Start.Date
              }
              fmt.Printf("%d) %s (%s)\n", i+1, item.Summary, date)
            }
          }
          EOF
          
          # Run go mod tidy to ensure dependencies are in sync
          go mod tidy
          
          # Run the test script
          CALENDAR_ID="${{ secrets.CALENDAR_ID }}" go run test_calendar.go

      - name: Test Pike13 API Connection
        run: |
          # Create a separate directory for the Pike13 test
          mkdir -p test_pike13
          cd test_pike13
          
          # Create go.mod file with minimal dependencies
          cat > go.mod << EOF
          module test_pike13

          go 1.21
          EOF
          
          # Create a simple test script to check Pike13 API connectivity
          cat > test_pike13.go << EOF
          package main

          import (
            "fmt"
            "io"
            "net/http"
            "os"
            "time"
          )

          func main() {
            # Get Pike13 client ID from environment
            clientID := os.Getenv("PIKE13_CLIENT_ID")
            if clientID == "" {
              fmt.Println("PIKE13_CLIENT_ID not found in environment")
              os.Exit(1)
            }
            
            # Get Pike13 URL from environment or use default
            pike13URL := os.Getenv("PIKE13_URL")
            if pike13URL == "" {
              pike13URL = "https://herosjourneyfitness.pike13.com/api/v2/front/event_occurrences.json"
            }
            
            # Calculate date range
            now := time.Now()
            fromDate := now.Format(time.RFC3339)
            toDate := now.AddDate(0, 0, 7).Format(time.RFC3339)
            
            # Build URL
            url := fmt.Sprintf("%s?from=%s&to=%s&client_id=%s", 
                               pike13URL, fromDate, toDate, clientID)
            
            fmt.Printf("Testing Pike13 API connectivity with URL: %s\n", url)
            
            # Make request
            resp, err := http.Get(url)
            if err != nil {
              fmt.Printf("Error making Pike13 API request: %v\n", err)
              os.Exit(1)
            }
            defer resp.Body.Close()
            
            if resp.StatusCode != http.StatusOK {
              body, _ := io.ReadAll(resp.Body)
              fmt.Printf("Pike13 API returned non-OK status: %d - %s\n", resp.StatusCode, string(body))
              os.Exit(1)
            }
            
            # Read the response body
            body, err := io.ReadAll(resp.Body)
            if err != nil {
              fmt.Printf("Error reading response body: %v\n", err)
              os.Exit(1)
            }
            
            # Print success and first 200 characters of response
            fmt.Printf("Successfully connected to Pike13 API\n")
            responseSummary := string(body)
            if len(responseSummary) > 200 {
                responseSummary = responseSummary[:200] + "..."
            }
            fmt.Printf("Response summary: %s\n", responseSummary)
          }
          EOF
          
          # Run go mod tidy to ensure dependencies are in sync
          go mod tidy
          
          # Run the test script
          go run test_pike13.go

      - name: Show environment information
        run: |
          echo "GitHub Actions environment:"
          echo "Ubuntu version: $(lsb_release -d)"
          echo "Go version: $(go version)"
          echo "Working directory: $(pwd)"
          echo "Directory listing:"
          ls -la
          
      - name: Upload test logs
        uses: actions/upload-artifact@v4
        with:
          name: pike13sync-preflight-logs
          path: |
            logs/
          retention-days: 7
