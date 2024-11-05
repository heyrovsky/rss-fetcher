package writer

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/heyrovsky/rsscurator/common/constants"
)

// IndexWriter writes the given string slice to a JSON file at the specified file location.
// It creates the necessary directories and writes the data atomically.
func IndexWriter(data []string, filelocation string) error {
	if len(data) == 0 {
		return errors.New("no data provided for writing")
	}
	if strings.TrimSpace(filelocation) == "" {
		return errors.New("file location cannot be empty")
	}

	// Clean and validate data entries
	cleanData := cleanAndTrimData(data)

	// Create directory if it does not exist
	if err := ensureDirectoryExists(filelocation); err != nil {
		return fmt.Errorf("could not create directory: %w", err)
	}

	// Serialize data to JSON format
	jsonData, err := json.MarshalIndent(cleanData, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON data: %w", err)
	}

	// Write to a temporary file first, then rename to ensure atomic write
	fullPath := filepath.Join(filelocation, constants.INDEX_FILE)
	if err := writeAtomic(fullPath, jsonData); err != nil {
		return err
	}

	return nil
}

// IndexReader reads and unmarshals JSON data from the specified file location.
// Returns an empty slice if the file does not exist.
func IndexReader(filelocation string) ([]string, error) {
	if strings.TrimSpace(filelocation) == "" {
		return nil, errors.New("file location cannot be empty")
	}

	fullPath := filepath.Join(filelocation, constants.INDEX_FILE)
	if !isValidPath(fullPath) {
		return nil, fmt.Errorf("invalid file path: %s", fullPath)
	}

	// Return empty slice if file does not exist
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return []string{}, nil
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Unmarshal data into a slice
	var result []string
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON data: %w", err)
	}

	return cleanAndTrimData(result), nil
}

// ensureDirectoryExists creates the directory if it does not already exist.
func ensureDirectoryExists(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}
	return nil
}

// cleanAndTrimData removes empty and whitespace-only entries from a slice.
func cleanAndTrimData(data []string) []string {
	cleaned := make([]string, 0, len(data))
	for _, item := range data {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return cleaned
}

// writeAtomic writes data to a temporary file and renames it to ensure atomic write.
func writeAtomic(filePath string, data []byte) error {
	tempFile := filePath + ".tmp"

	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	if err := os.Rename(tempFile, filePath); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to rename temporary file: %w", err)
	}

	return nil
}

// isValidPath validates that the given path is absolute and does not contain directory traversal segments.
func isValidPath(path string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil || strings.Contains(filepath.Clean(absPath), "..") {
		return false
	}
	return true
}
