package security

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"pqcd/crypto"
)

// RequestFeatures represents extracted features from a request
type RequestFeatures struct {
	// Request metadata
	Timestamp    int64  `json:"timestamp"`
	ClientIP     string `json:"client_ip"`
	UserAgent    string `json:"user_agent"`
	Algorithm    string `json:"algorithm"`
	Operation    string `json:"operation"`
	
	// Statistical features
	InputEntropy  float64 `json:"input_entropy"`
	InputSize     int     `json:"input_size"`
	
	// Temporal features
	InterRequestTime   float64 `json:"inter_request_time"`
	RequestsPerMinute  int     `json:"requests_per_minute"`
	
	// Result features
	OperationLatency   float64 `json:"operation_latency"`
	Success            bool    `json:"success"`
	
	// Sequence features (derived from history)
	LastOperation      string  `json:"last_operation"`
	SequenceHash       string  `json:"sequence_hash"`
}

// FeatureExtractor extracts features from HTTP requests
type FeatureExtractor struct {
	mu                 sync.Mutex
	// Keep track of last request time for each client IP
	lastRequestTime map[string]time.Time
	// Keep track of request count in the last minute for each client IP
	requestCounts map[string]int
	// Keep track of the last operation for each client IP
	lastOperations map[string]string
	// Keep track of operation sequence (used for anomaly detection)
	operationSequences map[string][]string
}

// NewFeatureExtractor creates a new feature extractor
func NewFeatureExtractor() *FeatureExtractor {
	return &FeatureExtractor{
		lastRequestTime:   make(map[string]time.Time),
		requestCounts:     make(map[string]int),
		lastOperations:    make(map[string]string),
		operationSequences: make(map[string][]string),
	}
}

// ExtractFeatures processes a request and extracts features
func (e *FeatureExtractor) ExtractFeatures(r *http.Request, algorithm crypto.Algorithm, operation string, inputData []byte, success bool, latencyMs float64) RequestFeatures {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	// Extract client IP
	clientIP := getClientIP(r)
	
	// Get current time
	now := time.Now()
	
	// Calculate inter-request time
	var interRequestTime float64
	lastRequest, exists := e.lastRequestTime[clientIP]
	if exists {
		interRequestTime = now.Sub(lastRequest).Seconds()
	} else {
		interRequestTime = -1 // First request from this IP
	}
	e.lastRequestTime[clientIP] = now
	
	// Calculate requests per minute
	// This is simplified, in a production system we would use a sliding window
	e.requestCounts[clientIP]++
	requestsPerMinute := e.requestCounts[clientIP]
	
	// Get last operation for sequence analysis
	lastOperation := e.lastOperations[clientIP]
	e.lastOperations[clientIP] = operation
	
	// Update operation sequence
	sequence := e.operationSequences[clientIP]
	if len(sequence) > 9 {
		// Keep last 10 operations
		sequence = sequence[1:]
	}
	sequence = append(sequence, operation)
	e.operationSequences[clientIP] = sequence
	
	// Calculate sequence hash for pattern detection
	sequenceHash := hashSequence(sequence)
	
	// Calculate entropy of input data
	entropy := calculateEntropy(inputData)
	
	return RequestFeatures{
		Timestamp:       now.Unix(),
		ClientIP:        clientIP,
		UserAgent:       r.UserAgent(),
		Algorithm:       string(algorithm),
		Operation:       operation,
		InputEntropy:    entropy,
		InputSize:       len(inputData),
		InterRequestTime: interRequestTime,
		RequestsPerMinute: requestsPerMinute,
		OperationLatency: latencyMs,
		Success:          success,
		LastOperation:    lastOperation,
		SequenceHash:     sequenceHash,
	}
}

// Helper functions

// calculateEntropy calculates Shannon entropy of a byte slice
func calculateEntropy(data []byte) float64 {
	if len(data) == 0 {
		return 0
	}
	
	// Count occurrences of each byte value
	counts := make(map[byte]int)
	for _, b := range data {
		counts[b]++
	}
	
	// Calculate entropy
	entropy := 0.0
	for _, count := range counts {
		p := float64(count) / float64(len(data))
		entropy -= p * math.Log2(p)
	}
	
	return entropy
}

// getClientIP extracts the client IP from a request
func getClientIP(r *http.Request) string {
	// Check for X-Forwarded-For header first (for clients behind proxies)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// The X-Forwarded-For header can contain multiple IPs
		// The leftmost is the original client IP
		return strings.Split(forwarded, ",")[0]
	}
	
	// Otherwise, use RemoteAddr
	return strings.Split(r.RemoteAddr, ":")[0]
}

// hashSequence creates a hash of the operation sequence for pattern matching
func hashSequence(sequence []string) string {
	if len(sequence) == 0 {
		return ""
	}
	
	// Join the sequence elements with a delimiter
	joined := strings.Join(sequence, ":")
	
	// Hash the joined string
	hash := sha256.Sum256([]byte(joined))
	
	return hex.EncodeToString(hash[:])
} 