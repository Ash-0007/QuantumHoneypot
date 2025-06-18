package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	// Define a simple multiplexer
	mux := http.NewServeMux()

	// Add routes
	mux.HandleFunc("/api/status", statusHandler)
	mux.HandleFunc("/api/keys/generate", keyGenerationHandler)
	mux.HandleFunc("/api/health", healthHandler)
	mux.HandleFunc("/api/decoys/generate", decoyGenerationHandler)
	mux.HandleFunc("/api/encrypt", encryptHandler)
	mux.HandleFunc("/api/decrypt", decryptHandler)

	// Create a server with CORS middleware
	handler := corsMiddleware(mux)

	// Start the server
	port := "8083"
	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for all responses
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
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

	// Mock response for POST requests
	response := map[string]interface{}{
		"public_key":   "mock_public_key_data_12345",
		"private_key":  "mock_private_key_data_67890",
		"algorithm":    "kyber",
		"fingerprint":  "abcdef123456",
		"generated_at": time.Now(),
		"decoy_count":  5,
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

	// Mock response for POST requests
	mockDecoys := []string{
		"kyber512",
		"dilithium2",
		"falcon512",
		"sphincs-haraka-128f",
		"saber",
		"ntru-hps-2048-509",
		"frodokem-640-aes",
		"picnic-L1-FS",
		"rainbow-Ia-classic",
		"mceliece348864",
	}

	response := map[string]interface{}{
		"decoys":      mockDecoys,
		"target":      "kyber768",
		"complexity":  5,
		"count":       len(mockDecoys),
		"generated_at": time.Now(),
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

	// Mock response for POST requests
	response := map[string]interface{}{
		"ciphertext":   "encrypted_data_mock_123456789",
		"nonce":        "random_nonce_value_987654321",
		"algorithm":    "kyber",
		"encrypted_at": time.Now(),
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

	// Mock response for POST requests
	response := map[string]interface{}{
		"plaintext":    "Hello, post-quantum world!",
		"algorithm":    "kyber",
		"decrypted_at": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
} 