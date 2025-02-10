package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/radu2020/planet/config"
	"github.com/radu2020/planet/internal/service"
	"github.com/radu2020/planet/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type application struct {
	dataService *service.DataService
	server      *http.Server
}

func main() {
	// Config
	cfg := config.LoadConfig()

	// Database
	db, err := sql.Open("postgres", cfg.Database.ConnectionInfo())
	if err != nil {
		log.Fatalf("Failed to open postgres connection: %v", err)
	}
	defer db.Close()

	// Storage
	storage := storage.NewSqlStorage(db)

	// Service
	dataService := service.NewDataService(storage)

	// App
	app := &application{dataService: dataService}

	// Handlers
	http.HandleFunc("GET /files/collection", app.getCollectionHandler)
	http.HandleFunc("GET /organizations/ids", app.getOrgIDsHandler)

	// Create server
	app.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: nil,
	}

	// Start server in a goroutine
	go func() {
		fmt.Printf("Starting the server on :%d...\n", cfg.Port)
		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// Wait for SIGINT or SIGTERM to gracefully shut down
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a shutdown signal
	<-quit
	log.Println("Shutting down server...")

	// Set a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shut down the server
	if err := app.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown failed: %v", err)
	}

	log.Println("Server gracefully stopped")
}

func (app *application) getCollectionHandler(w http.ResponseWriter, r *http.Request) {
	// Get data
	collection, err := app.dataService.GetCollection()
	if err != nil {
		log.Printf("Failed to fetch collection: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal
	payload, err := collection.MarshalJSON()
	if err != nil {
		log.Printf("GeoJSON Collection marshal failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send back data
	w.Header().Set("Content-Disposition", "attachment; filename=test.geojson")
	w.Header().Set("Content-Type", "application/text")
	http.ServeContent(w, r, "test.geojson", time.Now(), bytes.NewReader(payload))
}

func (app *application) getOrgIDsHandler(w http.ResponseWriter, r *http.Request) {
	// Get data
	payload, err := app.dataService.GetOrgIDs()
	if err != nil {
		log.Printf("Fetching OrgIDs failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send back data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Response failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
