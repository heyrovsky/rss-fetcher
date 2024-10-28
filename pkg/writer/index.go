package writer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/heyrovsky/rsscurator/common/constants"
)

func IndexWriter(data []string, filelocation string) error {
	if len(data) == 0 {
		return fmt.Errorf("empty data provided for writing")
	}
	if strings.TrimSpace(filelocation) == "" {
		return fmt.Errorf("invalid file location provided")
	}

	cleanData := make([]string, 0, len(data))
	for _, item := range data {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			cleanData = append(cleanData, trimmed)
		}
	}

	jsonData, err := json.MarshalIndent(cleanData, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON data: %w", err)
	}

	if err := os.MkdirAll(filelocation, 0755); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	fullPath := filepath.Join(filelocation, constants.INDEX_FILE)

	tempFile := fullPath + ".tmp"
	if err := os.WriteFile(tempFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	if err := os.Rename(tempFile, fullPath); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to finalize file write: %w", err)
	}

	return nil
}

func IndexReader(filelocation string) ([]string, error) {
	if strings.TrimSpace(filelocation) == "" {
		return nil, fmt.Errorf("invalid file location provided")
	}

	fullPath := filepath.Join(filelocation, constants.INDEX_FILE)
	if !isValidPath(fullPath) {
		return nil, fmt.Errorf("invalid file path: %s", fullPath)
	}

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return []string{}, nil
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	if len(data) == 0 {
		return []string{}, nil
	}

	var result []string

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON data: %w", err)
	}

	cleanResult := make([]string, 0, len(result))
	for _, item := range result {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			cleanResult = append(cleanResult, trimmed)
		}
	}

	return cleanResult, nil
}

func isValidPath(path string) bool {
	if !filepath.IsAbs(path) {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return false
		}
		path = absPath
	}

	cleanPath := filepath.Clean(path)

	return !strings.Contains(cleanPath, "..")
}
