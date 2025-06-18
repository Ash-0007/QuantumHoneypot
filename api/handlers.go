package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"pqcd/crypto"
	"pqcd/benchmark"
)

// CryptoHandler handles crypto API requests
type CryptoHandler struct {
	registry *crypto.Registry
	metrics  *benchmark.MetricsCollector
}

// NewCryptoHandler creates a new handler for crypto operations
func NewCryptoHandler(registry *crypto.Registry, metrics *benchmark.MetricsCollector) *CryptoHandler {
	return &CryptoHandler{
		registry: registry,
		metrics:  metrics,
	}
}

// KeyGenRequest is empty for now since key generation doesn't need input
type KeyGenRequest struct{}

// KeyGenResponse is the response for key generation
type KeyGenResponse struct {
	PublicKey   string    `json:"publicKey"`
	PrivateKey  string    `json:"privateKey"`
	Algorithm   string    `json:"algorithm"`
	Fingerprint string    `json:"fingerprint"`
	Decoys      []string  `json:"decoys"`
	GeneratedAt time.Time `json:"generatedAt"`
}

// EncapsulateRequest is the request for encapsulation
type EncapsulateRequest struct {
	PublicKey string `json:"publicKey"`
	Algorithm string `json:"algorithm"`
}

// EncapsulateResponse is the response for encapsulation
type EncapsulateResponse struct {
	Ciphertext   string `json:"ciphertext"`
	SharedSecret string `json:"sharedSecret"`
}

// DecapsulateRequest is the request for decapsulation
type DecapsulateRequest struct {
	PrivateKey string `json:"privateKey"`
	Ciphertext string `json:"ciphertext"`
	Algorithm  string `json:"algorithm"`
}

// DecapsulateResponse is the response for decapsulation
type DecapsulateResponse struct {
	SharedSecret string `json:"sharedSecret"`
}

// SignRequest is the request for signing
type SignRequest struct {
	PrivateKey string `json:"privateKey"`
	Message    string `json:"message"`
}

// SignResponse is the response for signing
type SignResponse struct {
	Signature string `json:"signature"`
}

// VerifyRequest is the request for verification
type VerifyRequest struct {
	PublicKey string `json:"publicKey"`
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

// VerifyResponse is the response for verification
type VerifyResponse struct {
	Valid bool `json:"valid"`
}

// ErrorResponse represents an API error
type ErrorResponse struct {
	Error string `json:"error"`
}

// HealthCheckResponse represents the health check response
type HealthCheckResponse struct {
	Status string `json:"status"`
}

// DecoyGenerationRequest is the request for decoy generation
type DecoyGenerationRequest struct {
	Target     string `json:"target"`
	Complexity int    `json:"complexity"`
	Count      int    `json:"count"`
}

// DecoyGenerationResponse is the response for decoy generation
type DecoyGenerationResponse struct {
	Decoys []string `json:"decoys"`
}

// HandleHealthCheck handles health check requests
func (h *CryptoHandler) HandleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := HealthCheckResponse{Status: "ok"}
		respondWithJSON(w, http.StatusOK, response)
	}
}

// HandleDecoyGeneration handles decoy generation requests
func (h *CryptoHandler) HandleDecoyGeneration() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DecoyGenerationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		// Placeholder: Generate simple random strings as decoys
		decoys := make([]string, 0, req.Count)
		for i := 0; i < req.Count; i++ {
			// Generate a random-like string based on target and complexity
			hash := sha256.Sum256([]byte(fmt.Sprintf("%s-%d-%d", req.Target, req.Complexity, i)))
			decoys = append(decoys, hex.EncodeToString(hash[:10]))
		}

		response := DecoyGenerationResponse{
			Decoys: decoys,
		}

		respondWithJSON(w, http.StatusOK, response)
	}
}

// HandleKeyGen handles key generation requests
func (h *CryptoHandler) HandleKeyGen() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		algorithm := crypto.Algorithm(vars["alg"])
		logrus.WithField("algorithm", algorithm).Info("Handling key generation request")
		
		start := time.Now()
		var keyPair crypto.KeyPair
		var err error
		
		// Check if it's a KEM or signature algorithm
		if strings.HasPrefix(string(algorithm), "ml-kem") || algorithm == crypto.AlgECDH {
			provider, err := h.registry.GetKEMProvider(algorithm)
			if err != nil {
				logrus.WithError(err).Error("Failed to get KEM provider")
				respondWithError(w, http.StatusBadRequest, fmt.Sprintf("unsupported algorithm: %s", algorithm))
				return
			}
			keyPair, err = provider.KeyGen()
		} else {
			provider, err := h.registry.GetSignatureProvider(algorithm)
			if err != nil {
				logrus.WithError(err).Error("Failed to get signature provider")
				respondWithError(w, http.StatusBadRequest, fmt.Sprintf("unsupported algorithm: %s", algorithm))
				return
			}
			keyPair, err = provider.KeyGen()
		}
		
		if err != nil {
			logrus.WithError(err).Error("Key generation failed")
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("key generation failed: %v", err))
			return
		}
		
		duration := time.Since(start)
		h.metrics.RecordOperation(algorithm, "KeyGen", duration, len(keyPair.PublicKey), len(keyPair.PrivateKey), true)
		
		// Create a simple fingerprint (SHA-256 hash of the public key)
		hash := sha256.Sum256(keyPair.PublicKey)
		fingerprint := hex.EncodeToString(hash[:])

		// Encode keys as hex strings for JSON response
		response := KeyGenResponse{
			PublicKey:   hex.EncodeToString(keyPair.PublicKey),
			PrivateKey:  hex.EncodeToString(keyPair.PrivateKey),
			Algorithm:   string(keyPair.Algorithm),
			Fingerprint: fingerprint,
			Decoys:      []string{}, // Placeholder for now
			GeneratedAt: time.Now(),
		}
		
		respondWithJSON(w, http.StatusOK, response)
	}
}

// HandleEncapsulate handles encapsulation requests
func (h *CryptoHandler) HandleEncapsulate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req EncapsulateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		algorithm := crypto.Algorithm(req.Algorithm)
		
		// Decode public key from hex
		publicKey, err := hex.DecodeString(req.PublicKey)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid public key format")
			return
		}
		
		// Get the KEM provider
		provider, err := h.registry.GetKEMProvider(algorithm)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("unsupported algorithm: %s", algorithm))
			return
		}
		
		// Perform encapsulation
		start := time.Now()
		ciphertext, sharedSecret, err := provider.Encapsulate(publicKey)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("encapsulation failed: %v", err))
			return
		}
		duration := time.Since(start)
		
		h.metrics.RecordOperation(algorithm, "Encapsulate", duration, len(publicKey), len(ciphertext), true)
		
		// Prepare response
		response := EncapsulateResponse{
			Ciphertext:  hex.EncodeToString(ciphertext),
			SharedSecret: hex.EncodeToString(sharedSecret),
		}
		
		respondWithJSON(w, http.StatusOK, response)
	}
}

// HandleDecapsulate handles decapsulation requests
func (h *CryptoHandler) HandleDecapsulate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DecapsulateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		algorithm := crypto.Algorithm(req.Algorithm)
		
		// Decode private key and ciphertext from hex
		privateKey, err := hex.DecodeString(req.PrivateKey)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid private key format")
			return
		}
		
		ciphertext, err := hex.DecodeString(req.Ciphertext)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid ciphertext format")
			return
		}
		
		// Get the KEM provider
		provider, err := h.registry.GetKEMProvider(algorithm)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("unsupported algorithm: %s", algorithm))
			return
		}
		
		// Perform decapsulation
		start := time.Now()
		sharedSecret, err := provider.Decapsulate(privateKey, ciphertext)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("decapsulation failed: %v", err))
			return
		}
		duration := time.Since(start)
		
		h.metrics.RecordOperation(algorithm, "Decapsulate", duration, len(privateKey), len(ciphertext), true)
		
		// Prepare response
		response := DecapsulateResponse{
			SharedSecret: hex.EncodeToString(sharedSecret),
		}
		
		respondWithJSON(w, http.StatusOK, response)
	}
}

// HandleSign handles signature creation requests
func (h *CryptoHandler) HandleSign() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		algorithm := crypto.Algorithm(vars["alg"])
		
		var req SignRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		
		// Decode private key from hex
		privateKey, err := hex.DecodeString(req.PrivateKey)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid private key format")
			return
		}
		
		// Get the signature provider
		provider, err := h.registry.GetSignatureProvider(algorithm)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("unsupported algorithm: %s", algorithm))
			return
		}
		
		// Perform signing
		start := time.Now()
		signature, err := provider.Sign(privateKey, []byte(req.Message))
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("signing failed: %v", err))
			return
		}
		duration := time.Since(start)
		
		h.metrics.RecordOperation(algorithm, "Sign", duration, len(privateKey), len(signature), true)
		
		// Prepare response
		response := SignResponse{
			Signature: hex.EncodeToString(signature),
		}
		
		respondWithJSON(w, http.StatusOK, response)
	}
}

// HandleVerify handles signature verification requests
func (h *CryptoHandler) HandleVerify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		algorithm := crypto.Algorithm(vars["alg"])
		
		var req VerifyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		
		// Decode public key and signature from hex
		publicKey, err := hex.DecodeString(req.PublicKey)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid public key format")
			return
		}
		
		signature, err := hex.DecodeString(req.Signature)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid signature format")
			return
		}
		
		// Get the signature provider
		provider, err := h.registry.GetSignatureProvider(algorithm)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("unsupported algorithm: %s", algorithm))
			return
		}
		
		// Perform verification
		start := time.Now()
		valid, err := provider.Verify(publicKey, []byte(req.Message), signature)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("verification failed: %v", err))
			return
		}
		duration := time.Since(start)
		
		h.metrics.RecordOperation(algorithm, "Verify", duration, len(publicKey), len(signature), valid)
		
		// Prepare response
		response := VerifyResponse{
			Valid: valid,
		}
		
		respondWithJSON(w, http.StatusOK, response)
	}
}

// Helper functions for API responses

func respondWithError(w http.ResponseWriter, code int, message string) {
	logrus.WithFields(logrus.Fields{
		"status_code": code,
		"error":       message,
	}).Error("API error")
	
	respondWithJSON(w, code, ErrorResponse{Error: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		logrus.WithError(err).Error("Failed to encode JSON response")
	}
} 