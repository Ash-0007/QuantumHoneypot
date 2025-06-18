package security

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"pqcd/crypto"
)

// AISecurityMiddleware implements the security layer as HTTP middleware
type AISecurityMiddleware struct {
	FeatureExtractor *FeatureExtractor
	AnomalyDetector  *AnomalyDetector
	ResponseEngine   *ResponseEngine
}

// NewAISecurityMiddleware creates a new AI security middleware
func NewAISecurityMiddleware() *AISecurityMiddleware {
	return &AISecurityMiddleware{
		FeatureExtractor: NewFeatureExtractor(),
		AnomalyDetector:  NewAnomalyDetector(),
		ResponseEngine:   NewResponseEngine(),
	}
}

// responseBodyWriter is a custom response writer that captures the entire response
type responseBodyWriter struct {
	http.ResponseWriter
	body   *bytes.Buffer
	header http.Header
	status int
}

func newResponseBodyWriter(w http.ResponseWriter) *responseBodyWriter {
	return &responseBodyWriter{
		ResponseWriter: w,
		body:           &bytes.Buffer{},
		header:         w.Header().Clone(),
	}
}

// Header returns the header map
func (w *responseBodyWriter) Header() http.Header {
	return w.header
}

// WriteHeader captures the status code but does not write to the original writer
func (w *responseBodyWriter) WriteHeader(statusCode int) {
	w.status = statusCode
}

// Write captures the response body but does not write to the original writer
func (w *responseBodyWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

// Middleware is the HTTP middleware function to process requests through the security layer
func (m *AISecurityMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip non-API paths or metrics endpoint
		if r.URL.Path == "/api/metrics" || !strings.HasPrefix(r.URL.Path, "/api") {
			next.ServeHTTP(w, r)
			return
		}
		
		// Extract path components to identify the algorithm and operation
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			next.ServeHTTP(w, r)
			return
		}
		
		algorithm := crypto.Algorithm(pathParts[2]) // e.g., "ml-kem-768"
		operation := pathParts[3]                  // e.g., "keygen"
		
		// Check if this client is already marked for redirection to honeypot
		clientIP := getClientIP(r)
		if m.ResponseEngine.ShouldRedirect(clientIP) {
			// In a real system, this would redirect to a honeypot
			// For now, we'll just log that we would have redirected
			logrus.WithFields(logrus.Fields{
				"client_ip": clientIP,
				"path":      r.URL.Path,
				"action":    "redirect",
			}).Info("Would redirect to honeypot")
			
			// Serve normal response to avoid tipping off the client
			next.ServeHTTP(w, r)
			return
		}
		
		// Check if this client should be throttled
		if m.ResponseEngine.ShouldThrottle(clientIP) {
			logrus.WithFields(logrus.Fields{
				"client_ip": clientIP,
				"path":      r.URL.Path,
				"action":    "throttle",
			}).Info("Throttling request")
			
			// Add artificial delay (25-75ms)
			time.Sleep(time.Duration(25+time.Now().UnixNano()%50) * time.Millisecond)
		}
		
		// Read the request body to extract features
		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		
		// Use our custom response writer to buffer the response
		responseWriter := newResponseBodyWriter(w)
		
		// Record start time to measure latency
		startTime := time.Now()
		
		// Process the request
		next.ServeHTTP(responseWriter, r)
		
		// Calculate latency
		latency := time.Since(startTime).Seconds() * 1000 // in milliseconds
		
		// Extract features
		features := m.FeatureExtractor.ExtractFeatures(r, algorithm, operation, bodyBytes, responseWriter.status >= 200 && responseWriter.status < 300, latency)
		
		// Determine if the client should receive deceptive responses
		shouldDeceive := m.ResponseEngine.ShouldDeceive(clientIP)
		if shouldDeceive && (operation == "encapsulate" || operation == "decapsulate" || operation == "verify") {
			logrus.WithFields(logrus.Fields{
				"client_ip": clientIP,
				"path":      r.URL.Path,
				"action":    "deceive",
			}).Info("Sending deceptive response")
			
			// Generate deceptive responses based on operation type
			// In a real system, these would be more sophisticated
			switch operation {
			case "encapsulate", "decapsulate":
				// Generate random "shared secret" that looks valid but isn't
				fakeSecret := make([]byte, 32)
				rand.Read(fakeSecret)
				// We would replace the response here, but that's complex in this simplified example
			case "verify":
				// Return "valid": false for signature verification
				// We would replace the response here
			}
		}
		
		// After normal processing, train the anomaly detector
		m.AnomalyDetector.Train(features)
		
		// Check for anomalies
		isAnomaly, anomalyType, anomalyScore := m.AnomalyDetector.Detect(features)
		if isAnomaly {
			// Classify the threat
			threat := m.ResponseEngine.ClassifyThreat(features, anomalyType, anomalyScore)
			
			logrus.WithFields(logrus.Fields{
				"client_ip":     clientIP,
				"threat_type":   threat.Type,
				"threat_level":  threat.Level,
				"threat_score":  threat.Score,
				"description":   threat.Description,
			}).Warn("Detected threat")
			
			// Decide and apply action
			action := m.ResponseEngine.DecideAction(threat)
			m.ResponseEngine.ApplyAction(action, clientIP)

			// Set headers to inform the client about the anomaly
			responseWriter.Header().Set("X-Anomaly-Detected", "true")
			responseWriter.Header().Set("X-Anomaly-Score", fmt.Sprintf("%f", threat.Score))
		}

		// Now, write the buffered response to the actual writer
		// Copy headers from our buffered writer to the real one
		for k, v := range responseWriter.header {
			w.Header()[k] = v
		}
		// Write the status code
		if responseWriter.status != 0 {
			w.WriteHeader(responseWriter.status)
		} else {
			w.WriteHeader(http.StatusOK) // Default if not set
		}
		// Write the body
		w.Write(responseWriter.body.Bytes())
	})
} 