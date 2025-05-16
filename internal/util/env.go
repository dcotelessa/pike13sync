package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadEnvFile loads environment variables from a .env file
func LoadEnvFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening .env file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || len(strings.TrimSpace(line)) == 0 {
			continue
		}

		// Split on first equals sign
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Skip malformed lines
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		value = strings.Trim(value, `"'`)

		// Set environment variable
		os.Setenv(key, value)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .env file: %v", err)
	}

	return nil
}

// DisplayEnvironmentInfo shows information about the environment
func DisplayEnvironmentInfo(config interface{}) {
	fmt.Println("\n=== ENVIRONMENT INFORMATION ===")
	
	// Print config details using reflection
	fmt.Printf("GOOGLE_CREDENTIALS_FILE: %s\n", os.Getenv("GOOGLE_CREDENTIALS_FILE"))
	fmt.Printf("CALENDAR_ID: %s\n", os.Getenv("CALENDAR_ID"))
	fmt.Printf("PIKE13_CLIENT_ID: %s\n", os.Getenv("PIKE13_CLIENT_ID"))
	fmt.Printf("DOCKER_ENV: %s\n", os.Getenv("DOCKER_ENV"))
	fmt.Printf("TZ: %s\n", os.Getenv("TZ"))
	
	// You can add more details about the config if needed
	
	fmt.Println("=========================")
}
