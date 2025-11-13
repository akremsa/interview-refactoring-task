package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Database struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
}

// NewDatabase creates a new database connection
func NewDatabase(host string, port int, username, password, dbName string) *Database {
	fmt.Printf("Establishing database connection to %s:%d...\n", host, port)

	return &Database{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		DBName:   dbName,
	}
}

// Save saves data to the database
func (db *Database) Save(data []byte) error {
	// Mock database save operation
	fmt.Printf("Saving data to database %s: %s\n", db.DBName, string(data))
	return nil
}

// DataStorage handles data persistence
type DataStorage struct {
	StorageType string
}

// SaveData saves data based on storage type
func (ds *DataStorage) SaveData(data []byte) error {
	// Hard-coded storage types - should use factory pattern
	if ds.StorageType == "file" {
		// File storage logic
		file, err := os.Create("data.txt")
		if err != nil {
			return err
		}
		defer file.Close()
		file.Write(data)
		fmt.Println("Data saved to file")
		return nil
	} else if ds.StorageType == "database" {
		db := NewDatabase("localhost", 5432, "admin", "password123", "app_database")
		db.Save(data)
		return nil
	} else {
		return fmt.Errorf("unsupported storage type: %s", ds.StorageType)
	}
}

// Request represents the incoming request
type Request struct {
	Data        []byte `json:"data"`
	StorageType string `json:"storage_type"`
}

// DataHandler handles HTTP requests
type DataHandler struct {
}

// HandleSaveData processes save data requests
func (h *DataHandler) HandleSaveData(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	var req Request
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	storage := &DataStorage{StorageType: req.StorageType}

	err = storage.SaveData(req.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data saved successfully"))
}

func main() {
	handler := &DataHandler{}
	http.HandleFunc("/save-data", handler.HandleSaveData)
	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
