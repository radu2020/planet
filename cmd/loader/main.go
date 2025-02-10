package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/radu2020/planet/config"
	"github.com/radu2020/planet/internal/data"
	"github.com/radu2020/planet/internal/storage"
	"log"
	"os"
)

func main() {
	// Config
	cfg := config.LoadConfig()

	// Database connection
	db, err := sql.Open("postgres", cfg.Database.ConnectionInfo())
	if err != nil {
		log.Fatalf("Failed to open postgres connection: %v", err)
	}
	defer db.Close()

	// Create table
	err = storage.CreateTable(db)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Open CSV file and process records
	file := openCSVFile(cfg.FilePath)
	defer file.Close()
	data.ProcessCSVRecords(file, db, cfg.BatchSize)
}

// openCSVFile opens the CSV file
func openCSVFile(filePath string) *os.File {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Error opening CSV file:", err)
	}
	return file
}
