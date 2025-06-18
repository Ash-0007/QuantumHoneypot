package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rs/cors"
	_ "github.com/mattn/go-sqlite3"
	
	"github.com/pqcd/backend/crypto"
)

// Global database connection
var db *sql.DB

// Configuration
type Config struct {
	Port         string
	DatabasePath string
	AIServiceURL string
	LogLevel     string
}

// Response structures
type KeyResponse struct {
	PublicKey   string    `json:"public_key"`
	Algorithm   string    `json:"algorithm"`
	Fingerprint string    `json:"fingerprint"`
	GeneratedAt time.Time `json:"generated_at"`
	DecoyCount  int       `json:"decoy_count"`
}

type EncryptResponse struct {
	Ciphertext string    `json:"ciphertext"`
	Nonce      string    `json:"nonce"`
	Algorithm  string    `json:"algorithm"`
	EncryptedAt time.Time `json:"encrypted_at"`
}

type DecryptResponse struct {
	Plaintext   string    `json:"plaintext"`
	Algorithm   string    `json:"algorithm"`
	DecryptedAt time.Time `json:"decrypted_at"`
}

type DecoyResponse struct {
	Decoys      []string  `json:"decoys"`
	Target      string    `json:"target"`
	Complexity  int       `json:"complexity"`
	Count       int       `json:"count"`
	GeneratedAt time.Time `json:"generated_at"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}

// Request structures
type KeyRequest struct {
	Algorithm string `json:"algorithm"`
	Count     int    `json:"count"`
}

type EncryptRequest struct {
	Plaintext string `json:"plaintext"`
	PublicKey string `json:"public_key"`
	Algorithm string `json:"algorithm"`
}

type DecryptRequest struct {
	Ciphertext string `json:"ciphertext"`
	PrivateKey string `json:"private_key"`
	Nonce      string `json:"nonce"`
	Algorithm  string `json:"algorithm"`
}

type DecoyRequest struct {
	Target     string `json:"target"`
	Complexity int    `json:"complexity"`
	Count      int    `json:"count"`
}

func main() {
	// Load configuration
	config := loadConfig()

	// Initialize database
	initDB(config.DatabasePath)
	defer db.Close()

	// Create router
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/api/status", statusHandler)
	mux.HandleFunc("/api/health", healthHandler)
	mux.HandleFunc("/api/keys/generate", keyGenerationHandler)
	mux.HandleFunc("/api/decoys/generate", decoyGenerationHandler)
	mux.HandleFunc("/api/encrypt", encryptHandler)
	mux.HandleFunc("/api/decrypt", decryptHandler)

	// Add CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            config.LogLevel == "debug",
	})
	handler := c.Handler(mux)

	// Start server
	log.Printf("Server starting on port %s...", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, handler))
}

// Load configuration from environment variables
func loadConfig() *Config {
	return &Config{
		Port:         getEnv("PORT", "8083"),
		DatabasePath: getEnv("DB_PATH", "./pqcd.db"),
		AIServiceURL: getEnv("AI_SERVICE_URL", "http://localhost:5000"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
	}
}

// Helper to get environment variable with default
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Initialize database connection
func initDB(dbPath string) {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Printf("Connected to database at %s", dbPath)
	
	// Create tables if they don't exist
	createTables()
}

// Create necessary tables if they don't exist
func createTables() {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS key_pairs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			public_key BLOB NOT NULL,
			private_key BLOB NOT NULL,
			fingerprint TEXT NOT NULL,
			algorithm TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			is_real BOOLEAN DEFAULT 1,
			tags TEXT
		)`,
		`CREATE INDEX IF NOT EXISTS idx_key_pairs_fingerprint ON key_pairs(fingerprint)`,
		`CREATE INDEX IF NOT EXISTS idx_key_pairs_algorithm ON key_pairs(algorithm)`,
		`CREATE TABLE IF NOT EXISTS decoys (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			decoy_text TEXT NOT NULL,
			target_text TEXT NOT NULL,
			complexity INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			effectiveness_score REAL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_decoys_target ON decoys(target_text)`,
		`CREATE TABLE IF NOT EXISTS event_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			event_type TEXT NOT NULL,
			description TEXT,
			source_ip TEXT,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			severity TEXT CHECK (severity IN ('INFO', 'WARNING', 'ERROR', 'CRITICAL')),
			related_item_id INTEGER,
			related_item_type TEXT
		)`,
		`CREATE INDEX IF NOT EXISTS idx_event_logs_type_time ON event_logs(event_type, timestamp)`,
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			last_login TIMESTAMP,
			role TEXT CHECK (role IN ('admin', 'user', 'readonly')) DEFAULT 'user'
		)`,
	}

	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			log.Printf("Error creating table: %v", err)
		}
	}

	// Insert default admin user if it doesn't exist
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'admin'").Scan(&count)
	if err != nil {
		log.Printf("Error checking for admin user: %v", err)
		return
	}

	if count == 0 {
		// Insert default admin user (password: admin123)
		_, err = db.Exec("INSERT INTO users (username, password_hash, role) VALUES ('admin', '$2b$12$1B7tQMt1IjuOXZvGdyz9A.J1DWnbwOgqYJMBpWY2Rcx.UmNMy9.cG', 'admin')")
		if err != nil {
			log.Printf("Error creating admin user: %v", err)
		} else {
			log.Printf("Created default admin user")
		}
	}
}

// Status handler
func statusHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"service":   "pqcd-backend",
		"status":    "operational",
		"version":   "0.1.0",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Health handler
func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Check database connection
	err := db.Ping()
	if err != nil {
		sendErrorResponse(w, "Database connection error", http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]string{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Key generation handler
func keyGenerationHandler(w http.ResponseWriter, r *http.Request) {
	// Only handle POST and OPTIONS methods
	if r.Method != "POST" && r.Method != "OPTIONS" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// For OPTIONS requests, the headers are already set by the middleware
	if r.Method == "OPTIONS" {
		return
	}

	// Parse request
	var req KeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Invalid request format", http.StatusBadRequest, err.Error())
		return
	}

	// Validate request
	if req.Algorithm == "" {
		req.Algorithm = crypto.AlgoKyber // Default algorithm
	}
	if req.Count <= 0 {
		req.Count = 5 // Default decoy count
	}

	// Generate key pair
	keyPair, err := crypto.GenerateKeyPair(req.Algorithm)
	if err != nil {
		sendErrorResponse(w, "Key generation failed", http.StatusInternalServerError, err.Error())
		return
	}

	// Generate fingerprint
	fingerprint := crypto.FingerPrint(keyPair.PublicKey)

	// Store in database
	_, err = db.Exec(
		"INSERT INTO key_pairs (public_key, private_key, fingerprint, algorithm, is_real) VALUES (?, ?, ?, ?, ?)",
		keyPair.PublicKey, keyPair.PrivateKey, fingerprint, keyPair.Algorithm, true,
	)
	if err != nil {
		sendErrorResponse(w, "Failed to store key pair", http.StatusInternalServerError, err.Error())
		return
	}

	// Generate decoys
	decoys, err := crypto.GenerateCognitiveDecoyKeys(keyPair, req.Count)
	if err != nil {
		sendErrorResponse(w, "Failed to generate decoys", http.StatusInternalServerError, err.Error())
		return
	}

	// Store decoys in database
	for _, decoy := range decoys {
		decoyFingerprint := crypto.FingerPrint(decoy.PublicKey)
		_, err = db.Exec(
			"INSERT INTO key_pairs (public_key, private_key, fingerprint, algorithm, is_real) VALUES (?, ?, ?, ?, ?)",
			decoy.PublicKey, decoy.PrivateKey, decoyFingerprint, decoy.Algorithm, false,
		)
		if err != nil {
			log.Printf("Failed to store decoy: %v", err)
		}
	}

	// Prepare response
	response := KeyResponse{
		PublicKey:   base64.StdEncoding.EncodeToString(keyPair.PublicKey),
		Algorithm:   keyPair.Algorithm,
		Fingerprint: fingerprint,
		GeneratedAt: time.Now(),
		DecoyCount:  req.Count,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Decoy generation handler
func decoyGenerationHandler(w http.ResponseWriter, r *http.Request) {
	// Only handle POST and OPTIONS methods
	if r.Method != "POST" && r.Method != "OPTIONS" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// For OPTIONS requests, the headers are already set by the middleware
	if r.Method == "OPTIONS" {
		return
	}

	// Parse request
	var req DecoyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Invalid request format", http.StatusBadRequest, err.Error())
		return
	}

	// Validate request
	if req.Target == "" {
		sendErrorResponse(w, "Target is required", http.StatusBadRequest, "")
		return
	}
	if req.Complexity <= 0 {
		req.Complexity = 5 // Default complexity
	}
	if req.Count <= 0 {
		req.Count = 10 // Default count
	}

	// Call AI service to generate decoys
	aiServiceURL := loadConfig().AIServiceURL + "/generate"
	aiReq, err := json.Marshal(map[string]interface{}{
		"target":     req.Target,
		"complexity": req.Complexity,
		"count":      req.Count,
	})
	if err != nil {
		sendErrorResponse(w, "Failed to prepare AI request", http.StatusInternalServerError, err.Error())
		return
	}

	// Make HTTP request to AI service
	aiResp, err := http.Post(aiServiceURL, "application/json", ioutil.NopCloser(bytes.NewReader(aiReq)))
	if err != nil {
		// Fallback to local decoy generation if AI service is unavailable
		log.Printf("AI service unavailable: %v. Using fallback decoys.", err)
		
		// Generate some fallback decoys
		mockDecoys := []string{
			req.Target + "_v1",
			req.Target + "_prime",
			req.Target + "-light",
			req.Target + "-extended",
			"light" + req.Target,
			"extended" + req.Target,
			req.Target + "512",
			req.Target + "1024",
			req.Target + "-hps",
			req.Target + "-hrss",
		}

		// Limit to requested count
		if len(mockDecoys) > req.Count {
			mockDecoys = mockDecoys[:req.Count]
		}

		response := DecoyResponse{
			Decoys:      mockDecoys,
			Target:      req.Target,
			Complexity:  req.Complexity,
			Count:       len(mockDecoys),
			GeneratedAt: time.Now(),
		}

		// Store decoys in database
		for _, decoy := range mockDecoys {
			_, err = db.Exec(
				"INSERT INTO decoys (decoy_text, target_text, complexity) VALUES (?, ?, ?)",
				decoy, req.Target, req.Complexity,
			)
			if err != nil {
				log.Printf("Failed to store decoy: %v", err)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	defer aiResp.Body.Close()

	// Parse AI service response
	var aiResponse map[string]interface{}
	if err := json.NewDecoder(aiResp.Body).Decode(&aiResponse); err != nil {
		sendErrorResponse(w, "Failed to parse AI response", http.StatusInternalServerError, err.Error())
		return
	}

	// Extract decoys from AI response
	decoys, ok := aiResponse["decoys"].([]interface{})
	if !ok {
		sendErrorResponse(w, "Invalid AI response format", http.StatusInternalServerError, "")
		return
	}

	// Convert to string array
	decoyStrings := make([]string, 0, len(decoys))
	for _, d := range decoys {
		if str, ok := d.(string); ok {
			decoyStrings = append(decoyStrings, str)
			
			// Store in database
			_, err = db.Exec(
				"INSERT INTO decoys (decoy_text, target_text, complexity) VALUES (?, ?, ?)",
				str, req.Target, req.Complexity,
			)
			if err != nil {
				log.Printf("Failed to store decoy: %v", err)
			}
		}
	}

	// Prepare response
	response := DecoyResponse{
		Decoys:      decoyStrings,
		Target:      req.Target,
		Complexity:  req.Complexity,
		Count:       len(decoyStrings),
		GeneratedAt: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Encrypt handler
func encryptHandler(w http.ResponseWriter, r *http.Request) {
	// Only handle POST and OPTIONS methods
	if r.Method != "POST" && r.Method != "OPTIONS" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// For OPTIONS requests, the headers are already set by the middleware
	if r.Method == "OPTIONS" {
		return
	}

	// Parse request
	var req EncryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Invalid request format", http.StatusBadRequest, err.Error())
		return
	}

	// Validate request
	if req.Plaintext == "" {
		sendErrorResponse(w, "Plaintext is required", http.StatusBadRequest, "")
		return
	}
	if req.PublicKey == "" {
		sendErrorResponse(w, "Public key is required", http.StatusBadRequest, "")
		return
	}
	if req.Algorithm == "" {
		req.Algorithm = crypto.AlgoKyber // Default algorithm
	}

	// Decode public key
	publicKey, err := base64.StdEncoding.DecodeString(req.PublicKey)
	if err != nil {
		sendErrorResponse(w, "Invalid public key format", http.StatusBadRequest, err.Error())
		return
	}

	// Encrypt data
	encryptedData, err := crypto.Encrypt([]byte(req.Plaintext), publicKey, req.Algorithm)
	if err != nil {
		sendErrorResponse(w, "Encryption failed", http.StatusInternalServerError, err.Error())
		return
	}

	// Log encryption event
	_, err = db.Exec(
		"INSERT INTO event_logs (event_type, description, severity) VALUES (?, ?, ?)",
		"encryption", fmt.Sprintf("Encrypted message using %s algorithm", req.Algorithm), "INFO",
	)
	if err != nil {
		log.Printf("Failed to log encryption event: %v", err)
	}

	// Prepare response
	response := EncryptResponse{
		Ciphertext: base64.StdEncoding.EncodeToString(encryptedData.Ciphertext),
		Nonce:      base64.StdEncoding.EncodeToString(encryptedData.Nonce),
		Algorithm:  encryptedData.Algorithm,
		EncryptedAt: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Decrypt handler
func decryptHandler(w http.ResponseWriter, r *http.Request) {
	// Only handle POST and OPTIONS methods
	if r.Method != "POST" && r.Method != "OPTIONS" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// For OPTIONS requests, the headers are already set by the middleware
	if r.Method == "OPTIONS" {
		return
	}

	// Parse request
	var req DecryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Invalid request format", http.StatusBadRequest, err.Error())
		return
	}

	// Validate request
	if req.Ciphertext == "" {
		sendErrorResponse(w, "Ciphertext is required", http.StatusBadRequest, "")
		return
	}
	if req.PrivateKey == "" {
		sendErrorResponse(w, "Private key is required", http.StatusBadRequest, "")
		return
	}
	if req.Nonce == "" {
		sendErrorResponse(w, "Nonce is required", http.StatusBadRequest, "")
		return
	}
	if req.Algorithm == "" {
		req.Algorithm = crypto.AlgoKyber // Default algorithm
	}

	// Decode request data
	ciphertext, err := base64.StdEncoding.DecodeString(req.Ciphertext)
	if err != nil {
		sendErrorResponse(w, "Invalid ciphertext format", http.StatusBadRequest, err.Error())
		return
	}

	privateKey, err := base64.StdEncoding.DecodeString(req.PrivateKey)
	if err != nil {
		sendErrorResponse(w, "Invalid private key format", http.StatusBadRequest, err.Error())
		return
	}

	nonce, err := base64.StdEncoding.DecodeString(req.Nonce)
	if err != nil {
		sendErrorResponse(w, "Invalid nonce format", http.StatusBadRequest, err.Error())
		return
	}

	// Create encrypted data structure
	encryptedData := &crypto.EncryptedData{
		Ciphertext: ciphertext,
		Algorithm:  req.Algorithm,
		Nonce:      nonce,
	}

	// Decrypt data
	plaintext, err := crypto.Decrypt(encryptedData, privateKey)
	if err != nil {
		sendErrorResponse(w, "Decryption failed", http.StatusInternalServerError, err.Error())
		return
	}

	// Log decryption event
	_, err = db.Exec(
		"INSERT INTO event_logs (event_type, description, severity) VALUES (?, ?, ?)",
		"decryption", fmt.Sprintf("Decrypted message using %s algorithm", req.Algorithm), "INFO",
	)
	if err != nil {
		log.Printf("Failed to log decryption event: %v", err)
	}

	// Prepare response
	response := DecryptResponse{
		Plaintext:   string(plaintext),
		Algorithm:   req.Algorithm,
		DecryptedAt: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper function to send error responses
func sendErrorResponse(w http.ResponseWriter, message string, code int, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	
	response := ErrorResponse{
		Error:   message,
		Code:    code,
		Details: details,
	}
	
	json.NewEncoder(w).Encode(response)
	
	// Log error
	log.Printf("Error: %s (%d): %s", message, code, details)
} 