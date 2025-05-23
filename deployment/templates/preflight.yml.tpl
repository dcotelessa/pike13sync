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
        uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'  # Updated to a newer version for better compatibility
          
      - name: Create directory structure
        run: |
          mkdir -p config
          mkdir -p credentials
          mkdir -p logs

      - name: Set up Google credentials
        run: |
          echo '${{ secrets.GOOGLE_CREDENTIALS }}' > credentials/credentials.json
          chmod 600 credentials/credentials.json

      - name: Set up Pike13 credentials
        run: |
          echo '{"client_id": "${{ secrets.PIKE13_CLIENT_ID }}"}' > credentials/pike13_credentials.json

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
            
            // Read credentials from file
            credBytes, err := os.ReadFile("../credentials/credentials.json")
            if err != nil {
              fmt.Printf("Error reading credentials file: %v\n", err)
              os.Exit(1)
            }
            
            // Configure JWT
            config, err := google.JWTConfigFromJSON(credBytes, calendar.CalendarScope)
            if err != nil {
              fmt.Printf("Error creating JWT config: %v\n", err)
              os.Exit(1)
            }
            
            // Create client
            client := config.Client(ctx)
            
            // Create calendar service
            srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
            if err != nil {
              fmt.Printf("Error creating calendar service: %v\n", err)
              os.Exit(1)
            }
            
            // Get calendar ID from environment
            calendarID := os.Getenv("CALENDAR_ID")
            if calendarID == "" {
              calendarID = "primary"
            }
            
            // List next 10 events
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
            
            // Print success and event count
            fmt.Printf("Successfully connected to Google Calendar API\n")
            fmt.Printf("Found %d upcoming events\n", len(events.Items))
            
            // Print upcoming events
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
            "encoding/json"
            "fmt"
            "io"
            "net/http"
            "os"
            "time"
          )

          type Pike13Credentials struct {
            ClientID string \`json:"client_id"\`
          }

          type Pike13Response struct {
            EventOccurrences []struct {
              ID   int    \`json:"id"\`
              Name string \`json:"name"\`
            } \`json:"event_occurrences"\`
          }

          func main() {
            // Read Pike13 credentials
            credFile, err := os.ReadFile("../credentials/pike13_credentials.json")
            if err != nil {
              fmt.Printf("Error reading Pike13 credentials: %v\n", err)
              os.Exit(1)
            }
            
            var creds Pike13Credentials
            if err := json.Unmarshal(credFile, &creds); err != nil {
              fmt.Printf("Error parsing Pike13 credentials: %v\n", err)
              os.Exit(1)
            }
            
            if creds.ClientID == "" {
              fmt.Println("Pike13 client ID not found in credentials file")
              os.Exit(1)
            }
            
            // Calculate date range
            now := time.Now()
            fromDate := now.Format(time.RFC3339)
            toDate := now.AddDate(0, 0, 7).Format(time.RFC3339)
            
            // Build URL
            url := fmt.Sprintf("https://herosjourneyfitness.pike13.com/api/v2/front/event_occurrences.json?from=%s&to=%s&client_id=%s", 
                               fromDate, toDate, creds.ClientID)
            
            // Make request
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
            
            // Read and parse response
            body, err := io.ReadAll(resp.Body)
            if err != nil {
              fmt.Printf("Error reading response body: %v\n", err)
              os.Exit(1)
            }
            
            var response Pike13Response
            if err := json.Unmarshal(body, &response); err != nil {
              fmt.Printf("Error parsing JSON response: %v\n", err)
              os.Exit(1)
            }
            
            // Print success and event count
            fmt.Printf("Successfully connected to Pike13 API\n")
            fmt.Printf("Retrieved %d events from Pike13\n", len(response.EventOccurrences))
            
            // Print sample events
            maxEvents := 5
            if len(response.EventOccurrences) < maxEvents {
              maxEvents = len(response.EventOccurrences)
            }
            
            for i := 0; i < maxEvents; i++ {
              fmt.Printf("%d) Event ID: %d, Name: %s\n", i+1, 
                         response.EventOccurrences[i].ID, 
                         response.EventOccurrences[i].Name)
            }
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
          echo "Docker version: $(docker --version)"
          echo "Working directory: $(pwd)"
          echo "Directory listing:"
          ls -la
