package writer

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/heyrovsky/rsscurator/common/utils"
	"github.com/heyrovsky/rsscurator/pkg/content"
)

// JsonWriter writes unique hashed news items to a JSON file, ensuring no duplicates across months.
func JsonWriter(data []content.NewsItemHashed, fileName, folderName, lastMonthfolderName string) error {
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

	fmt.Printf("Successfully written %d unique items to %s\n", len(uniqueData), fileName)
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

// writeJSONFile marshals and writes JSON data to a file.
func writeJSONFile(fileName string, data []content.NewsItemHashed) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	if err := os.WriteFile(fileName, jsonData, 0644); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}
