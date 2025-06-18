package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pqcd/backend/crypto"
	"github.com/rs/cors"
)

const PORT = 8082

// Request structures
type KeyGenRequest struct {
	Algorithm string `json:"algorithm"`
	Count     int    `json:"count"`
}

type EncryptRequest struct {
	PublicKey string `json:"public_key"`
	Algorithm string `json:"algorithm"`
	Data      string `json:"data"`
}

type DecryptRequest struct {
	PrivateKey string `json:"private_key"`
	Ciphertext string `json:"ciphertext"`
	Nonce      string `json:"nonce"`
	Algorithm  string `json:"algorithm"`
}

type DecoyGenRequest struct {
	Target     string `json:"target"`
	Complexity int    `json:"complexity"`
	Count      int    `json:"count"`
}

// Response structures
type KeyGenResponse struct {
	Algorithm   string    `json:"algorithm"`
	PublicKey   string    `json:"public_key"`
	PrivateKey  string    `json:"private_key"`
	Fingerprint string    `json:"fingerprint"`
	Decoys      []string  `json:"decoys,omitempty"`
	GeneratedAt time.Time `json:"generated_at"`
}

type EncryptResponse struct {
	Algorithm  string `json:"algorithm"`
	Ciphertext string `json:"ciphertext"`
	Nonce      string `json:"nonce"`
}

type DecryptResponse struct {
	Plaintext string `json:"plaintext"`
}

type DecoyGenResponse struct {
	Decoys []map[string]interface{} `json:"decoys"`
}

func main() {
	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/api/health", healthHandler)
	mux.HandleFunc("/api/status", statusHandler)
	mux.HandleFunc("/api/keys/generate", keyGenerateHandler)
	mux.HandleFunc("/api/encrypt", encryptHandler)
	mux.HandleFunc("/api/decrypt", decryptHandler)
	mux.HandleFunc("/api/decoys/generate", decoyGenerateHandler)

	// Setup CORS
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(mux)

	// Start server
	log.Printf("Starting API server on port %d...", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), handler))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "operational",
		"version":     "0.1.0",
		"algorithms":  []string{crypto.AlgoKyber, crypto.AlgoSaber, crypto.AlgoNTRU, crypto.AlgoDilithium},
		"connections": 1,
		"uptime":      "1h23m",
	})
}

func keyGenerateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request
	var req KeyGenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate algorithm
	algorithm := req.Algorithm
	if algorithm == "" {
		algorithm = crypto.AlgoKyber
	}

	// Generate key pair
	keyPair, err := crypto.GenerateKeyPair(algorithm)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating key pair: %s", err), http.StatusInternalServerError)
		return
	}

	// Convert binary keys to base64 for JSON
	publicKeyB64 := base64.StdEncoding.EncodeToString(keyPair.PublicKey)
	privateKeyB64 := base64.StdEncoding.EncodeToString(keyPair.PrivateKey)
	
	// Generate fingerprint
	fingerprint := crypto.FingerPrint(keyPair.PublicKey)

	// Generate decoys if requested
	var decoyKeys []string
	if req.Count > 0 {
		decoys, err := crypto.GenerateCognitiveDecoyKeys(keyPair, req.Count)
		if err == nil {
			for _, decoy := range decoys {
				decoyKeyB64 := base64.StdEncoding.EncodeToString(decoy.PublicKey)
				decoyKeys = append(decoyKeys, decoyKeyB64)
			}
		}
	}

	// Create response
	resp := KeyGenResponse{
		Algorithm:   algorithm,
		PublicKey:   publicKeyB64,
		PrivateKey:  privateKeyB64,
		Fingerprint: fingerprint,
		Decoys:      decoyKeys,
		GeneratedAt: time.Now(),
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func encryptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request
	var req EncryptRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Decode base64 public key
	publicKey, err := base64.StdEncoding.DecodeString(req.PublicKey)
	if err != nil {
		http.Error(w, "Invalid public key", http.StatusBadRequest)
		return
	}

	// Encrypt data
	algorithm := req.Algorithm
	if algorithm == "" {
		algorithm = crypto.AlgoKyber
	}

	encryptedData, err := crypto.Encrypt([]byte(req.Data), publicKey, algorithm)
	if err != nil {
		http.Error(w, fmt.Sprintf("Encryption error: %s", err), http.StatusInternalServerError)
		return
	}

	// Create response
	resp := EncryptResponse{
		Algorithm:  algorithm,
		Ciphertext: base64.StdEncoding.EncodeToString(encryptedData.Ciphertext),
		Nonce:      base64.StdEncoding.EncodeToString(encryptedData.Nonce),
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func decryptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request
	var req DecryptRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Decode base64 private key
	privateKey, err := base64.StdEncoding.DecodeString(req.PrivateKey)
	if err != nil {
		http.Error(w, "Invalid private key", http.StatusBadRequest)
		return
	}

	// Decode base64 ciphertext
	ciphertext, err := base64.StdEncoding.DecodeString(req.Ciphertext)
	if err != nil {
		http.Error(w, "Invalid ciphertext", http.StatusBadRequest)
		return
	}

	// Decode base64 nonce
	nonce, err := base64.StdEncoding.DecodeString(req.Nonce)
	if err != nil {
		http.Error(w, "Invalid nonce", http.StatusBadRequest)
		return
	}

	// Create encrypted data structure
	encryptedData := &crypto.EncryptedData{
		Ciphertext: ciphertext,
		Algorithm:  req.Algorithm,
		Nonce:      nonce,
	}

	// Decrypt data
	decryptedData, err := crypto.Decrypt(encryptedData, privateKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Decryption error: %s", err), http.StatusInternalServerError)
		return
	}

	// Create response
	resp := DecryptResponse{
		Plaintext: string(decryptedData),
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func decoyGenerateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request
	var req DecoyGenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Target == "" {
		http.Error(w, "Target text is required", http.StatusBadRequest)
		return
	}

	if req.Count <= 0 {
		req.Count = 5
	}

	if req.Complexity <= 0 {
		req.Complexity = 5
	}

	// Generate decoys - in a real implementation this would use machine learning
	// Here we'll just use a simple algorithm to generate text variations
	decoys := generateSimpleDecoys(req.Target, req.Complexity, req.Count)

	// Create response
	resp := DecoyGenResponse{
		Decoys: decoys,
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Simple algorithm to generate decoy text
func generateSimpleDecoys(target string, complexity int, count int) []map[string]interface{} {
	decoys := make([]map[string]interface{}, count)
	
	// Generate random bytes for variation
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	
	// Create decoys with varying similarity
	for i := 0; i < count; i++ {
		// Calculate similarity - higher complexity = less similar
		similarity := 1.0 - float64(complexity)/10.0 - float64(i)/float64(count*2)
		if similarity < 0.1 {
			similarity = 0.1
		}
		
		// Generate decoy value based on target
		var decoy string
		
		switch {
		case strings.Contains(target, "kyber"):
			suffix := hex.EncodeToString(randBytes)[0:4]
			decoy = fmt.Sprintf("kyber%d", 512+i*128)
			if i%3 == 0 {
				decoy = fmt.Sprintf("kyber_%s", suffix)
			}
		case strings.Contains(target, "saber"):
			suffix := hex.EncodeToString(randBytes)[0:4]
			decoy = fmt.Sprintf("saber%d", 1+i)
			if i%3 == 0 {
				decoy = fmt.Sprintf("lightsaber_%s", suffix)
			}
		case strings.Contains(target, "ntru"):
			decoy = fmt.Sprintf("ntru_%d", 509+i*100)
		case strings.Contains(target, "dilithium"):
			decoy = fmt.Sprintf("dilithium_mode%d", 1+i)
		default:
			// Add some variations to the target
			chars := []rune(target)
			if len(chars) > 3 {
				pos := 1 + (i % (len(chars) - 2))
				chars[pos] = chars[pos] + 1
			}
			decoy = string(chars)
		}
		
		decoys[i] = map[string]interface{}{
			"value":      decoy,
			"similarity": similarity,
		}
	}
	
	return decoys
} 