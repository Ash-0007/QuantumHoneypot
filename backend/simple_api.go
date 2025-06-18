package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/rs/cors"
)

const (
	PORT                = 8082
	KyberCiphertextSize = 1088 // Kyber-768 ciphertext size
	KyberSharedKeySize  = 32   // Kyber-768 shared key size
)

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Create a new router
	mux := http.NewServeMux()

	// Add API routes
	mux.HandleFunc("/api/health", healthHandler)
	mux.HandleFunc("/api/status", statusHandler)
	mux.HandleFunc("/api/keys/generate", keyGenerateHandler)
	mux.HandleFunc("/api/encrypt", encryptHandler)
	mux.HandleFunc("/api/decrypt", decryptHandler)
	mux.HandleFunc("/api/decoys/generate", decoyGenerateHandler)

	// Add CORS middleware
	handler := cors.Default().Handler(mux)

	// Start server
	log.Printf("Starting simplified API server on port %d...", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), handler))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	setJSONResponse(w)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	setJSONResponse(w)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "operational",
		"version":     "0.1.0",
		"algorithms":  []string{"kyber", "saber", "ntru", "dilithium"},
		"connections": 1,
		"uptime":      "1h23m",
	})
}

func keyGenerateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		setJSONResponse(w)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Algorithm string `json:"algorithm"`
		Count     int    `json:"count"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	algorithm := req.Algorithm
	if algorithm == "" {
		algorithm = "kyber"
	}

	count := req.Count
	if count <= 0 {
		count = 0
	}

	// Generate mock key sizes based on algorithm
	var pubKeySize, privKeySize int
	switch algorithm {
	case "kyber":
		pubKeySize = 1184
		privKeySize = 2400
	case "saber":
		pubKeySize = 992
		privKeySize = 2304
	case "ntru":
		pubKeySize = 1138
		privKeySize = 1450
	case "dilithium":
		pubKeySize = 1312
		privKeySize = 2528
	default:
		pubKeySize = 1184
		privKeySize = 2400
	}

	// Generate random bytes for keys
	publicKey := randomBytes(pubKeySize)
	privateKey := randomBytes(privKeySize)

	// Generate fingerprint
	fingerprintBytes := sha256.Sum256(publicKey)
	fingerprint := hex.EncodeToString(fingerprintBytes[:16])

	// Generate decoys
	var decoys []string
	for i := 0; i < count; i++ {
		decoy := randomBytes(pubKeySize)
		decoys = append(decoys, base64.StdEncoding.EncodeToString(decoy))
	}

	// Create response
	resp := map[string]interface{}{
		"algorithm":    algorithm,
		"public_key":   base64.StdEncoding.EncodeToString(publicKey),
		"private_key":  base64.StdEncoding.EncodeToString(privateKey),
		"fingerprint":  fingerprint,
		"decoys":       decoys,
		"generated_at": time.Now(),
	}

	setJSONResponse(w)
	json.NewEncoder(w).Encode(resp)
}

func encryptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		setJSONResponse(w)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PublicKey string `json:"public_key"`
		Algorithm string `json:"algorithm"`
		Data      string `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Simulate encryption
	// In a real KEM, the public key would be used to derive a shared secret.
	// Here, we simulate this by creating a secret and hiding it in the ciphertext.
	sharedSecret := randomBytes(KyberSharedKeySize)
	encapsulation := make([]byte, KyberCiphertextSize)
	copy(encapsulation, sharedSecret) // "Hide" the key for decryption

	// Generate a nonce
	nonce := randomBytes(24)

	// Encrypt the data using a stream cipher from the shared secret and nonce
	keyStream := deriveKeyStream(sharedSecret, nonce, len(req.Data))
	encryptedData := make([]byte, len(req.Data))
	for i := range req.Data {
		encryptedData[i] = req.Data[i] ^ keyStream[i]
	}

	// The final ciphertext includes the "encapsulated" secret and the encrypted data
	finalCiphertext := append(encapsulation, encryptedData...)

	resp := map[string]interface{}{
		"algorithm":  req.Algorithm,
		"ciphertext": base64.StdEncoding.EncodeToString(finalCiphertext),
		"nonce":      base64.StdEncoding.EncodeToString(nonce),
	}

	setJSONResponse(w)
	json.NewEncoder(w).Encode(resp)
}

func decryptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		setJSONResponse(w)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PrivateKey string `json:"private_key"`
		Ciphertext string `json:"ciphertext"`
		Nonce      string `json:"nonce"`
		Algorithm  string `json:"algorithm"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Decode from base64
	ciphertext, err := base64.StdEncoding.DecodeString(req.Ciphertext)
	if err != nil {
		http.Error(w, "Invalid ciphertext format", http.StatusBadRequest)
		return
	}
	nonce, err := base64.StdEncoding.DecodeString(req.Nonce)
	if err != nil {
		http.Error(w, "Invalid nonce format", http.StatusBadRequest)
		return
	}

	// In a real system, the private key would be used to decapsulate the shared secret.
	// Here, we retrieve the "hidden" secret from the start of the ciphertext.
	if len(ciphertext) < KyberCiphertextSize {
		http.Error(w, "Invalid ciphertext: too short", http.StatusBadRequest)
		return
	}
	sharedSecret := ciphertext[:KyberSharedKeySize]
	actualCiphertext := ciphertext[KyberCiphertextSize:]

	// Re-derive the key stream to decrypt
	keyStream := deriveKeyStream(sharedSecret, nonce, len(actualCiphertext))
	plaintext := make([]byte, len(actualCiphertext))
	for i := range actualCiphertext {
		plaintext[i] = actualCiphertext[i] ^ keyStream[i]
	}

	resp := map[string]interface{}{
		"plaintext": string(plaintext),
	}

	setJSONResponse(w)
	json.NewEncoder(w).Encode(resp)
}

func decoyGenerateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		setJSONResponse(w)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Target     string `json:"target"`
		Complexity int    `json:"complexity"`
		Count      int    `json:"count"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

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

	// Generate decoys
	decoys := generateSimpleDecoys(req.Target, req.Complexity, req.Count)

	resp := map[string]interface{}{
		"decoys": decoys,
	}

	setJSONResponse(w)
	json.NewEncoder(w).Encode(resp)
}

func setJSONResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func randomBytes(length int) []byte {
	bytes := make([]byte, length)
	rand.Read(bytes) // Use crypto/rand for better random data
	return bytes
}

// deriveKeyStream creates a key stream from a key and nonce using SHA-256
func deriveKeyStream(key, nonce []byte, length int) []byte {
	keyStream := make([]byte, length)

	// Use SHA-256 to create a key stream seed
	h := sha256.New()
	h.Write(key)
	h.Write(nonce)
	seed := h.Sum(nil)

	// Expand the seed to the required length
	for i := 0; i < length; i += sha256.Size {
		h := sha256.New()
		h.Write(seed)
		// Add a counter to generate unique hashes for each chunk
		h.Write([]byte{byte(i / sha256.Size)})
		digest := h.Sum(nil)

		// Copy as much as needed
		copy(keyStream[i:], digest)
	}

	return keyStream
}

func generateSimpleDecoys(target string, complexity int, count int) []map[string]interface{} {
	decoys := make([]map[string]interface{}, count)
	
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
			sizes := []string{"512", "768", "1024", "1536"}
			size := sizes[rand.Intn(len(sizes))]
			if rand.Intn(3) == 0 {
				decoy = fmt.Sprintf("kyber_%s", size)
			} else {
				decoy = fmt.Sprintf("kyber%s", size)
			}
		case strings.Contains(target, "saber"):
			variants := []string{"lightsaber", "saber", "firesaber"}
			decoy = variants[rand.Intn(len(variants))]
		case strings.Contains(target, "ntru"):
			variants := []string{"hrss701", "hps2048509", "hps4096821"}
			decoy = fmt.Sprintf("ntru_%s", variants[rand.Intn(len(variants))])
		case strings.Contains(target, "dilithium"):
			mode := rand.Intn(5) + 1
			decoy = fmt.Sprintf("dilithium_mode%d", mode)
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