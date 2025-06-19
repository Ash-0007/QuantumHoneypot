package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	
	"pqcd/benchmark"
	"pqcd/crypto"
)

// RegisterRoutes sets up all API routes
func RegisterRoutes(r *mux.Router) {
	// Create the crypto registry
	registry := crypto.DefaultRegistry()
	
	// Create the metrics collector
	metrics := benchmark.NewMetricsCollector()
	
	// Create the handler
	handler := NewCryptoHandler(registry, metrics)
	
	// Set up the API subrouter with common path prefix
	api := r.PathPrefix("/api").Subrouter()
	
	// Register KEM endpoints
	registerKEMRoutes(api, handler)
	
	// Register signature endpoints
	registerSignatureRoutes(api, handler)
	
	// Register metrics endpoint
	api.HandleFunc("/metrics", metrics.HandleMetrics()).Methods("GET")

	// Register health check endpoint
	api.HandleFunc("/health", handler.HandleHealthCheck()).Methods("GET")
	
	// Register decoy generation endpoint
	api.HandleFunc("/decoys/generate", handler.HandleDecoyGeneration()).Methods("POST")

	// Register general encrypt/decrypt endpoints
	api.HandleFunc("/encrypt", handler.HandleEncapsulate()).Methods("POST")
	api.HandleFunc("/decrypt", handler.HandleDecapsulate()).Methods("POST")

	// This is a placeholder for your actual PQC API endpoints.
	// The AI middleware will wrap these routes.
	api.HandleFunc("/api/crypto/{algorithm}/{operation}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		algorithm := vars["algorithm"]
		operation := vars["operation"]

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := fmt.Sprintf(`{"status": "ok", "message": "successfully performed %s for %s"}`, operation, algorithm)
		fmt.Fprintln(w, response)
	}).Methods("POST", "GET")

	logrus.Info("API routes registered")
}

// registerKEMRoutes registers the Key Encapsulation Mechanism endpoints
func registerKEMRoutes(r *mux.Router, handler *CryptoHandler) {
	kemRoutes := r.PathPrefix("/{alg:(?:ml-kem-768|ecdh)}").Subrouter()
	kemRoutes.HandleFunc("/keygen", handler.HandleKeyGen()).Methods("POST")
	kemRoutes.HandleFunc("/encapsulate", handler.HandleEncapsulate()).Methods("POST")
	kemRoutes.HandleFunc("/decapsulate", handler.HandleDecapsulate()).Methods("POST")
}

// registerSignatureRoutes registers the Digital Signature endpoints
func registerSignatureRoutes(r *mux.Router, handler *CryptoHandler) {
	sigRoutes := r.PathPrefix("/{alg:(?:ml-dsa-65|ecdsa)}").Subrouter()
	sigRoutes.HandleFunc("/keygen", handler.HandleKeyGen()).Methods("POST")
	sigRoutes.HandleFunc("/sign", handler.HandleSign()).Methods("POST")
	sigRoutes.HandleFunc("/verify", handler.HandleVerify()).Methods("POST")
} 