package util

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindProjectRoot attempts to locate the root directory of the project
func FindProjectRoot() (string, error) {
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
