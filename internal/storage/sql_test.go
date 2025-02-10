package storage

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestInsertBatch(t *testing.T) {
	// Mock DB setup
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	defer db.Close()

	// Mock the execution of the insert query
	timestampStr := "2025-02-09T15:04:05Z"
	parsedTime, _ := time.Parse(time.RFC3339, timestampStr)
	utcTime := parsedTime.UTC()

	mock.ExpectExec("INSERT INTO data").
		WithArgs("1", `{"type":"Feature"}`, utcTime).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Test data
	batch := [][]string{
		{"1", `{"type":"Feature"}`, "2025-02-09T15:04:05Z"},
	}

	// Call the insertBatch function
	err = InsertBatch(db, batch)
	assert.NoError(t, err)

	// Ensure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unmet expectations: %s", err)
	}
}

// Test InsertBatch with invalid timestamp
func TestInsertBatch_InvalidTimestamp(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	batch := [][]string{
		{"1", `{"type":"Feature"}`, "invalid-timestamp"},
	}

	// Since the function skips invalid timestamps, we do not expect any query execution
	err = InsertBatch(db, batch)
	assert.Error(t, err)
}

// Test GetCollection function
func TestGetCollection(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlStorage := NewSqlStorage(db)

	// Mock rows
	rows := sqlmock.NewRows([]string{"footprints_used"}).
		AddRow([]byte(`{"type":"Feature"}`))

	mock.ExpectQuery("SELECT footprints_used FROM data;").WillReturnRows(rows)

	fc, err := sqlStorage.GetCollection()

	assert.NoError(t, err)
	assert.Equal(t, "Feature", fc.Features[0].Type)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test GetCollection with query error
func TestGetCollection_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewSqlStorage(db)

	mock.ExpectQuery("SELECT footprints_used FROM data;").WillReturnError(sql.ErrNoRows)

	featureCollection, err := storage.GetCollection()
	assert.Error(t, err)
	assert.Nil(t, featureCollection)
}

// Test GetOrgIDs function
func TestGetOrgIDs(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewSqlStorage(db)

	rows := sqlmock.NewRows([]string{"org_id"}).AddRow(1).AddRow(2)
	mock.ExpectQuery("SELECT DISTINCT org_id FROM data").WillReturnRows(rows)

	orgIDs, err := storage.GetOrgIDs()
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2}, orgIDs)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test GetOrgIDs with error
func TestGetOrgIDs_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := NewSqlStorage(db)

	mock.ExpectQuery("SELECT DISTINCT org_id FROM data").WillReturnError(sql.ErrConnDone)

	orgIDs, err := storage.GetOrgIDs()
	assert.Error(t, err)
	assert.Nil(t, orgIDs)
}

// Test CreateTable function
func TestCreateTable(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS data").WillReturnResult(sqlmock.NewResult(1, 1))

	err = CreateTable(db)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test CreateTable with error
func TestCreateTable_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS data").WillReturnError(sql.ErrConnDone)

	err = CreateTable(db)
	assert.Error(t, err)
}
