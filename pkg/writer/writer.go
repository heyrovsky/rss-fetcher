package writer

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/heyrovsky/rsscurator/common/utils"
	"github.com/heyrovsky/rsscurator/pkg/content"
	"go.uber.org/zap"
)

// JsonWriter writes unique hashed news items to a JSON file, ensuring no duplicates across months.
func JsonWriter(data []content.NewsItemHashed, fileName, folderName, lastMonthfolderName string, logger zap.Logger) error {
	// Return early if there's no data to process
	if len(data) == 0 {
		return errors.New("no data provided to JsonWriter")
	}

	// Read and handle indices for the current and previous month
	currentMonthIndex, err := readIndexSafely(folderName)
	if err != nil {
		return fmt.Errorf("error reading current month index: %w", err)
	}

	previousMonthIndex, err := readIndexSafely(lastMonthfolderName)
	if err != nil {
		return fmt.Errorf("error reading previous month index: %w", err)
	}

	// Combine hashes from both months, ensuring uniqueness
	existingHashes, err := utils.ListUnique(currentMonthIndex, previousMonthIndex)
	if err != nil {
		return fmt.Errorf("error combining existing hashes: %w", err)
	}

	// Extract hashes from the new data and determine unique new items
	newHashes := content.ExtractHashes(data)
	uniqueNewHashes, err := utils.ListUniqueList2(existingHashes, newHashes)
	if err != nil {
		return fmt.Errorf("error finding unique hashes: %w", err)
	}

	// Filter data to only include unique items
	uniqueData := filterUniqueData(data, uniqueNewHashes)
	if len(uniqueData) == 0 {
		return errors.New("no unique data to write")
	}

	// Update and save the current month index with new unique hashes
	updatedIndex, err := utils.ListUnique(currentMonthIndex, uniqueNewHashes)
	if err != nil {
		return fmt.Errorf("error updating current month index: %w", err)
	}
	if err := IndexWriter(updatedIndex, folderName); err != nil {
		return fmt.Errorf("error writing updated index: %w", err)
	}

	// Write the unique data to the JSON file
	if err := writeJSONFile(fileName, uniqueData); err != nil {
		return fmt.Errorf("error writing JSON data: %w", err)
	}

	logger.Info(fmt.Sprintf("Successfully written %d unique items to %s\n", len(uniqueData), fileName))
	return nil
}

// readIndexSafely reads an index file, returning an empty slice if the file doesn't exist.
func readIndexSafely(folderName string) ([]string, error) {
	index, err := IndexReader(folderName)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return index, nil
}

// filterUniqueData filters the given data items to only include those with unique hashes.
func filterUniqueData(data []content.NewsItemHashed, uniqueHashes []string) []content.NewsItemHashed {
	uniqueSet := make(map[string]struct{}, len(uniqueHashes))
	for _, hash := range uniqueHashes {
		uniqueSet[hash] = struct{}{}
	}

	var uniqueData []content.NewsItemHashed
	for _, item := range data {
		if _, exists := uniqueSet[item.Hash]; exists {
			uniqueData = append(uniqueData, item)
		}
	}
	return uniqueData
}

// writeJSONFile appends or writes JSON data to a file based on its existence.
func writeJSONFile(fileName string, data []content.NewsItemHashed) error {
	var allData []content.NewsItemHashed

	// Check if the file already exists
	if _, err := os.Stat(fileName); err == nil {
		// File exists, so read the current content
		existingData, err := readJSONFile(fileName)
		if err != nil {
			return fmt.Errorf("error reading existing JSON file: %w", err)
		}
		// Append the new data to the existing data
		allData = append(existingData, data...)
	} else if errors.Is(err, os.ErrNotExist) {
		// File doesn't exist, so we'll create a new one with only the new data
		allData = data
	} else {
		// Some other error occurred
		return fmt.Errorf("error checking file existence: %w", err)
	}

	// Marshal the combined data and write it to the file
	jsonData, err := json.MarshalIndent(allData, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	if err := os.WriteFile(fileName, jsonData, 0644); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}

// readJSONFile reads and unmarshals JSON data from a file.
func readJSONFile(fileName string) ([]content.NewsItemHashed, error) {
	fileData, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %w", err)
	}

	var data []content.NewsItemHashed
	if err := json.Unmarshal(fileData, &data); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON data: %w", err)
	}
	return data, nil
}
