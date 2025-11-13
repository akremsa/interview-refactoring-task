package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// SOLUTION: Proper design patterns implementation

// Storage interface - IMPLEMENTS Polymorphism and Strategy Pattern
type StorageInterface interface {
	Save(data []byte) error
}

// FileStorage implements Storage interface
type FileStorage struct {
	filename string
}

func (fs *FileStorage) Save(data []byte) error {
	file, err := os.Create(fs.filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	fmt.Println("Data saved to file")
	return nil
}

// DatabaseStorage implements Storage interface
type DatabaseStorage struct {
	db *DatabaseConnection
}

func (ds *DatabaseStorage) Save(data []byte) error {
	return ds.db.Save(data)
}

// DatabaseConnection - properly structured with dependency injection
type DatabaseConnection struct {
	Host      string
	Port      int
	Username  string
	Password  string
	DBName    string
	connected bool
}

// NewDatabaseConnection creates a new database connection - PROPER INITIALIZATION
func NewDatabaseConnection(host string, port int, username, password, dbName string) (*DatabaseConnection, error) {
	fmt.Printf("Establishing database connection to %s:%d...\n", host, port)

	db := &DatabaseConnection{
		Host:      host,
		Port:      port,
		Username:  username,
		Password:  password,
		DBName:    dbName,
		connected: true,
	}

	fmt.Printf("Successfully connected to database: %s\n", dbName)
	return db, nil
}

func (db *DatabaseConnection) Save(data []byte) error {
	if !db.connected {
		return fmt.Errorf("database connection not established")
	}
	fmt.Printf("Saving data to database %s: %s\n", db.DBName, string(data))
	return nil
}

func (db *DatabaseConnection) Close() error {
	fmt.Printf("Closing database connection to %s\n", db.DBName)
	db.connected = false
	return nil
}

// StorageFactory - IMPLEMENTS Factory Pattern
type StorageFactory interface {
	CreateStorage(storageType string) (StorageInterface, error)
}

// ConcreteStorageFactory implements StorageFactory
type ConcreteStorageFactory struct {
	database *DatabaseConnection
}

func NewStorageFactory(database *DatabaseConnection) *ConcreteStorageFactory {
	return &ConcreteStorageFactory{
		database: database,
	}
}

func (f *ConcreteStorageFactory) CreateStorage(storageType string) (StorageInterface, error) {
	switch storageType {
	case "file":
		return &FileStorage{filename: "data.txt"}, nil
	case "database":
		if f.database == nil {
			return nil, fmt.Errorf("database connection not available")
		}
		return &DatabaseStorage{db: f.database}, nil
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", storageType)
	}
}

// SaveRequest - FOLLOWS Single Responsibility
type SaveRequest struct {
	Data        []byte `json:"data"`
	StorageType string `json:"storage_type"`
}

// Validator - IMPLEMENTS Single Responsibility
type RequestValidator struct{}

func NewRequestValidator() *RequestValidator {
	return &RequestValidator{}
}

func (v *RequestValidator) ValidateRequest(req *SaveRequest) error {
	if len(req.Data) == 0 {
		return fmt.Errorf("data cannot be empty")
	}
	if req.StorageType == "" {
		return fmt.Errorf("storage type cannot be empty")
	}
	validTypes := []string{"file", "database"}
	for _, validType := range validTypes {
		if req.StorageType == validType {
			return nil
		}
	}
	return fmt.Errorf("invalid storage type: %s", req.StorageType)
}

// DataService - IMPLEMENTS Single Responsibility and Dependency Injection
type DataService struct {
	factory   StorageFactory
	validator *RequestValidator
}

func NewDataService(factory StorageFactory, validator *RequestValidator) *DataService {
	return &DataService{
		factory:   factory,
		validator: validator,
	}
}

func (ds *DataService) SaveData(req *SaveRequest) error {
	// Validate request
	if err := ds.validator.ValidateRequest(req); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Use factory to create storage
	storage, err := ds.factory.CreateStorage(req.StorageType)
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}

	// Save data
	if err := storage.Save(req.Data); err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}

	return nil
}

// HTTPHandler - IMPLEMENTS Single Responsibility and Dependency Injection
type HTTPHandler struct {
	dataService *DataService
}

func NewHTTPHandler(dataService *DataService) *HTTPHandler {
	return &HTTPHandler{dataService: dataService}
}

func (h *HTTPHandler) HandleSaveData(w http.ResponseWriter, r *http.Request) {
	// Validate HTTP method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read and parse request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req SaveRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Process request
	err = h.dataService.SaveData(&req)
	if err != nil {
		// Determine appropriate HTTP status code based on error type
		statusCode := http.StatusInternalServerError
		if err.Error() == "validation failed" {
			statusCode = http.StatusBadRequest
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	// Send structured JSON response
	response := map[string]string{
		"message": "Data saved successfully",
		"status":  "success",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Configuration - IMPLEMENTS Configuration Management
type Configuration struct {
	Port         string
	DatabaseHost string
	DatabasePort int
	DatabaseUser string
	DatabasePass string
	DatabaseName string
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Port:         "8080",
		DatabaseHost: "localhost",
		DatabasePort: 5432,
		DatabaseUser: "admin",
		DatabasePass: "password123",
		DatabaseName: "app_database",
	}
}

// APIServer - IMPLEMENTS Proper Server Structure
type APIServer struct {
	config   *Configuration
	handler  *HTTPHandler
	database *DatabaseConnection
}

func NewAPIServer(config *Configuration) (*APIServer, error) {
	// Initialize database connection ONCE at startup
	database, err := NewDatabaseConnection(
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseUser,
		config.DatabasePass,
		config.DatabaseName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Create dependencies using dependency injection
	validator := NewRequestValidator()
	factory := NewStorageFactory(database)
	dataService := NewDataService(factory, validator)
	handler := NewHTTPHandler(dataService)

	return &APIServer{
		config:   config,
		handler:  handler,
		database: database,
	}, nil
}

func (s *APIServer) Start() error {
	http.HandleFunc("/save-data", s.handler.HandleSaveData)

	// Add health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{"status": "healthy"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	fmt.Printf("Server starting on :%s\n", s.config.Port)
	return http.ListenAndServe(":"+s.config.Port, nil)
}

func (s *APIServer) Shutdown() error {
	fmt.Println("Shutting down server...")
	return s.database.Close()
}

// Properly structured main function with dependency injection
func main() {
	// Load configuration
	config := NewConfiguration()

	// Initialize server with all dependencies
	server, err := NewAPIServer(config)
	if err != nil {
		log.Fatal("Failed to initialize server:", err)
	}

	// Graceful shutdown would be implemented here in production
	defer func() {
		if err := server.Shutdown(); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
	}()

	// Start server with proper error handling
	if err := server.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
