package util

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"reflect"
)

// LoadEnvFile loads environment variables from a .env file
func LoadEnvFile(filePath string) error {
	// Check if filePath is provided
	if filePath == "" {
		// Try to find .env in the project root
		rootDir, err := findProjectRoot()
		if err != nil {
			return fmt.Errorf("error finding project root: %v", err)
		}
		
		filePath = filepath.Join(rootDir, ".env")
	}
	
	// Check if .env file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("Notice: .env file does not exist at %s", filePath)
		return nil // Not an error, just continue without .env
	}
	
	// Open .env file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening .env file: %v", err)
	}
	defer file.Close()

	log.Printf("Loading environment variables from: %s", filePath)
	
	// Read line by line
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		
		// Skip comments and empty lines
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split on first equals sign
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Printf("Warning: Skipping malformed line %d in .env file: %s", lineNum, line)
			continue // Skip malformed lines
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		value = strings.Trim(value, `"'`)

		// Set environment variable if not already set
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
			log.Printf("Set environment variable: %s", key)
		} else {
			log.Printf("Environment variable already set, not overriding: %s", key)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .env file: %v", err)
	}

	return nil
}

// findProjectRoot attempts to locate the root directory of the project
func findProjectRoot() (string, error) {
	// If TEST_BASE_DIR is set, use it
	if testDir := os.Getenv("TEST_BASE_DIR"); testDir != "" {
		return testDir, nil
	}
	
	// If DOCKER_ENV is set, use /app
	if os.Getenv("DOCKER_ENV") == "true" {
		return "/app", nil
	}
	
	// Try to determine the project root from the current directory
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %v", err)
	}
	
	// Look for indicators of the project root (go.mod, .git, etc.)
	for dir != "/" && dir != "." && dir != "" {
		if fileExists(filepath.Join(dir, "go.mod")) || 
		   fileExists(filepath.Join(dir, ".git")) {
			return dir, nil
		}
		
		// Move up one directory
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break // Avoid infinite loop
		}
		dir = parentDir
	}
	
	// Fallback to current directory if project root not found
	cwd, err := os.Getwd()
	if err != nil {
		return ".", nil
	}
	return cwd, nil
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// DisplayEnvironmentInfo shows information about the environment
func DisplayEnvironmentInfo(config interface{}) {
	fmt.Println("\n=== ENVIRONMENT INFORMATION ===")
	
	// Print environment variables
	vars := []string{
		"GOOGLE_CREDENTIALS_FILE",
		"CALENDAR_ID",
		"PIKE13_CLIENT_ID",
		"PIKE13_URL",
		"DOCKER_ENV",
		"TZ",
		"LOG_PATH",
		"DRY_RUN",
	}
	
	fmt.Println("Environment Variables:")
	for _, v := range vars {
		value := os.Getenv(v)
		if value == "" {
			value = "(not set)"
		}
		fmt.Printf("  %s: %s\n", v, value)
	}
	
	// Print the config details if needed
	fmt.Printf("\n=== CONFIG INFORMATION ===\n")
	
	// Use reflection to print config details
	val := reflect.ValueOf(config)
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem() // Dereference pointer
		
		// If config is a struct, print its fields
		if val.Kind() == reflect.Struct {
			fmt.Println("Configuration Settings:")
			
			// Print each field name and value
			for i := 0; i < val.NumField(); i++ {
				field := val.Type().Field(i)
				value := val.Field(i)
				
				// Skip unexported fields
				if field.PkgPath != "" {
					continue
				}
				
				// Get JSON tag or use field name
				fieldName := field.Name
				if tag := field.Tag.Get("json"); tag != "" && tag != "-" {
					// Split tag to get the name part (before any comma)
					parts := strings.Split(tag, ",")
					if parts[0] != "" {
						fieldName = parts[0]
					}
				}
				
				// Format value based on kind
				var valueStr string
				switch value.Kind() {
				case reflect.String:
					valueStr = value.String()
				case reflect.Bool:
					valueStr = fmt.Sprintf("%v", value.Bool())
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					valueStr = fmt.Sprintf("%d", value.Int())
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					valueStr = fmt.Sprintf("%d", value.Uint())
				case reflect.Float32, reflect.Float64:
					valueStr = fmt.Sprintf("%g", value.Float())
				default:
					valueStr = fmt.Sprintf("%v", value.Interface())
				}
				
				fmt.Printf("  %s: %s\n", fieldName, valueStr)
			}
		} else {
			fmt.Printf("Config type: %T\n", config)
		}
	} else {
		fmt.Printf("Config type: %T\n", config)
	}
	
	fmt.Println("=========================")
}
