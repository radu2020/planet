package data

import (
	"database/sql"
	"encoding/csv"
	"github.com/radu2020/planet/internal/storage"
	"io"
	"log"
	"os"
)

// ProcessCSVRecords processes CSV records and inserts them into the database
func ProcessCSVRecords(file *os.File, db *sql.DB, batchSize int) {
	reader := csv.NewReader(file)

	// Skip header
	_, err := reader.Read()
	if err != nil {
		log.Fatalf("Error reading header: %v", err)
		return
	}

	var batch [][]string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading line: %v", err)
			continue
		}

		if isValidRecord(record) {
			batch = append(batch, record)
			if len(batch) >= batchSize {
				if err := storage.InsertBatch(db, batch); err != nil {
					log.Println("Batch insert error:", err)
				}
				batch = nil
			}
		}
	}

	if len(batch) > 0 {
		if err := storage.InsertBatch(db, batch); err != nil {
			log.Println("Final batch insert error:", err)
		}
	}

	log.Println("CSV data successfully loaded into the database!")
}
