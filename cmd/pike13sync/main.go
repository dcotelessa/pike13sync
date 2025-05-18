package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dcotelessa/pike13sync/internal/calendar"
	"github.com/dcotelessa/pike13sync/internal/config"
	"github.com/dcotelessa/pike13sync/internal/pike13"
	"github.com/dcotelessa/pike13sync/internal/sync"
	"github.com/dcotelessa/pike13sync/internal/util"
)

func main() {
	// Look for a .env file in the project root and load it
	envFile := ""
	
	// Allow specifying alternative .env file via ENV_FILE environment variable
	if customEnvFile := os.Getenv("ENV_FILE"); customEnvFile != "" {
		envFile = customEnvFile
	}
	
	// Load environment variables from .env file
	if err := util.LoadEnvFile(envFile); err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		fmt.Println("Continuing with existing environment variables")
	}

	// Setup logging
	logFile, err := util.SetupLogging()
	if err != nil {
		fmt.Printf("Error setting up logging: %v\n", err)
		fmt.Println("Continuing with console logging only")
	} else if logFile != nil {
		defer logFile.Close()
	}

	// Parse command-line flags
	dryRunFlag := flag.Bool("dry-run", false, "Dry run mode - don't actually modify Google Calendar")
	testFromDate := flag.String("from", "", "Test from date (format: 2025-01-01)")
	testToDate := flag.String("to", "", "Test to date (format: 2025-01-07)")
	debugMode := flag.Bool("debug", false, "Enable debug mode with extra logging")
	sampleOnly := flag.Bool("sample", false, "Only fetch and display sample events without syncing")
	configPath := flag.String("config", "", "Path to config file")
	showEnv := flag.Bool("show-env", false, "Show environment information and exit")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Printf("Error loading configuration: %v", err)
	}

	// Override with command-line flag if specified
	if *dryRunFlag {
		cfg.DryRun = true
	}

	if *showEnv {
		util.DisplayEnvironmentInfo(cfg)
		return
	}

	// Calculate date range
	fromDate, toDate := calculateDateRange(*testFromDate, *testToDate)
	log.Printf("Fetching events from %s to %s", fromDate, toDate)
	if cfg.DryRun || *debugMode {
		fmt.Printf("DRY RUN MODE: Fetching events from %s to %s\n", fromDate, toDate)
	}

	// Fetch Pike13 events
	pike13Client := pike13.NewClient(cfg)
	events, err := pike13Client.FetchEvents(fromDate, toDate)
	if err != nil {
		log.Printf("Error fetching Pike13 events: %v", err)
		if len(events.EventOccurrences) == 0 {
			log.Fatal("No events retrieved, exiting")
		}
	}

	eventCount := len(events.EventOccurrences)
	log.Printf("Retrieved %d events from Pike13", eventCount)
	if cfg.DryRun || *debugMode {
		fmt.Printf("Retrieved %d events from Pike13\n", eventCount)
	}

	// If sample-only mode, just display events and exit
	if *sampleOnly {
		pike13Client.DisplaySampleEvents(events)
		return
	}

	// Set up Google Calendar service
	calendarService, err := calendar.NewService(cfg)
	if err != nil {
		log.Fatalf("Error setting up Google Calendar: %v", err)
	}

	// Sync events
	syncService := sync.NewSyncService(calendarService, cfg)
	stats := syncService.SyncEvents(events.EventOccurrences)

	// Print summary
	printSummary(stats, cfg.DryRun)
}

func calculateDateRange(testFrom, testTo string) (string, string) {
	if testFrom != "" && testTo != "" {
		// Add time component if missing
		if len(testFrom) == 10 {
			testFrom += "T00:00:00Z"
		}
		if len(testTo) == 10 {
			testTo += "T00:00:00Z"
		}
		return testFrom, testTo
	}

	now := time.Now()
	
	// Special handling for Saturday
	if now.Weekday() == time.Saturday {
		// Start from today (Saturday)
		startDate := now
		
		// Calculate days until next Sunday (1 day from Saturday)
		daysUntilNextSunday := 1
		
		// Then add 7 more days to get to the following Sunday
		endDate := startDate.AddDate(0, 0, daysUntilNextSunday+7)
		
		return startDate.Format(time.RFC3339), endDate.Format(time.RFC3339)
	}
	
	// For all other days, get Sunday to Sunday
	// Calculate the most recent Sunday (start of week)
	daysToSubtract := int(now.Weekday())
	startOfWeek := now.AddDate(0, 0, -daysToSubtract)
	
	// End date is 7 days after start date (next Sunday)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	return startOfWeek.Format(time.RFC3339), endOfWeek.Format(time.RFC3339)
}

func printSummary(stats sync.SyncStats, dryRun bool) {
	log.Printf("Sync completed: %d created, %d updated, %d deleted, %d unchanged",
		stats.Created, stats.Updated, stats.Deleted, stats.Skipped)

	if dryRun {
		fmt.Printf("\n==== SYNC SUMMARY (DRY RUN) ====\n")
		fmt.Printf("Events that would be created: %d\n", stats.Created)
		fmt.Printf("Events that would be updated: %d\n", stats.Updated)
		fmt.Printf("Events that would be deleted: %d\n", stats.Deleted)
		fmt.Printf("Events that would be unchanged: %d\n", stats.Skipped)
		fmt.Printf("===============================\n")
		fmt.Println("No changes were made to Google Calendar (dry run mode)")
	} else {
		fmt.Printf("\n==== SYNC SUMMARY ====\n")
		fmt.Printf("Events created: %d\n", stats.Created)
		fmt.Printf("Events updated: %d\n", stats.Updated)
		fmt.Printf("Events deleted: %d\n", stats.Deleted)
		fmt.Printf("Events unchanged: %d\n", stats.Skipped)
		fmt.Printf("====================\n")
	}
}
