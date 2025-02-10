package data

import (
	"log"
	"strings"
	"time"
)

// isValidRecord validates if a record is well-formed
func isValidRecord(record []string) bool {
	if len(record) != 3 {
		log.Println("Invalid row. Each row must have exactly 3 columns")
		return false
	}

	for _, value := range record {
		if strings.TrimSpace(value) == "" {
			log.Println("Invalid row. Columns cannot be empty")
			return false
		}
	}

	if !isValidFootprintFormat(record[1]) {
		log.Println("Skipping invalid footprint format:", record[1])
		return false
	}

	if !isValidTimestamp(record[2]) {
		log.Println("Skipping invalid timestamp:", record[2])
		return false
	}

	return true
}

// isValidFootprintFormat checks if the footprint is in the correct JSON format
func isValidFootprintFormat(footprint string) bool {
	return strings.HasPrefix(footprint, `{"type":"Feature"`)
}

// isValidTimestamp checks if the timestamp is in the correct RFC3339 format
func isValidTimestamp(timestampStr string) bool {
	_, err := time.Parse(time.RFC3339, timestampStr)
	return err == nil
}
