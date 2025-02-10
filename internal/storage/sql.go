package storage

import (
	"database/sql"
	"fmt"
	"github.com/paulmach/orb/geojson"
	"log"
	"strings"
	"time"
)

type SqlStorage struct {
	db *sql.DB
}

func NewSqlStorage(db *sql.DB) *SqlStorage {
	return &SqlStorage{db: db}
}

// GetCollection queries the entities from the db and reads the entities
// directly in a geojson.Feature and returns a geojson.FeatureCollection
func (s *SqlStorage) GetCollection() (*geojson.FeatureCollection, error) {
	// Query db
	rows, err := s.db.Query("SELECT footprints_used FROM data;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create FeatureCollection
	fc := geojson.NewFeatureCollection()

	// Populate collection with features
	for rows.Next() {
		var payload []byte

		err := rows.Scan(&payload)
		if err != nil {
			log.Println(err)
			continue
		}
		f, err := geojson.UnmarshalFeature(payload)
		if err != nil {
			log.Println(err)
			continue
		}

		fc.Append(f)
	}

	return fc, nil
}

// GetOrgIDs fetches the Org IDs from the DB and returns a slice of int
func (s *SqlStorage) GetOrgIDs() ([]int, error) {
	rows, err := s.db.Query(`
		SELECT DISTINCT org_id 
		FROM data 
		WHERE footprints_used IS NOT NULL 
		AND jsonb_typeof(footprints_used) != 'null';
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgIDs []int
	for rows.Next() {
		var orgID int
		if err := rows.Scan(&orgID); err != nil {
			return nil, err
		}
		orgIDs = append(orgIDs, orgID)
	}

	return orgIDs, nil
}

// CreateTable creates the data table in the database if it doesn't exist
func CreateTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS data (
			org_id Int,
			footprints_used JSONB,
			source_event_timestamp timestamptz
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

// InsertBatch inserts a batch of records into the database
func InsertBatch(db *sql.DB, batch [][]string) error {
	query := "INSERT INTO data (org_id, footprints_used, source_event_timestamp) VALUES "
	values := []string{}
	args := []interface{}{}
	argCount := 1

	for _, record := range batch {
		timestampStr := record[2]
		parsedTime, err := time.Parse(time.RFC3339, timestampStr)
		if err != nil {
			log.Println("Skipping invalid timestamp:", timestampStr)
			continue
		}
		utcTime := parsedTime.UTC()

		values = append(values, fmt.Sprintf("($%d, $%d, $%d)", argCount, argCount+1, argCount+2))
		args = append(args, record[0], record[1], utcTime)
		argCount += 3
	}

	query += strings.Join(values, ",")
	_, err := db.Exec(query, args...)
	return err
}
