package security

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type AISecurityMiddleware struct {
	// any dependencies, e.g., a logger
}

func NewAISecurityMiddleware() *AISecurityMiddleware {
	return &AISecurityMiddleware{}
}

func (m *AISecurityMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// --- 1. Extract Request Details ---
		// Attempt to get the real IP, falling back to RemoteAddr.
		ip := r.Header.Get("X-Real-IP")
		if ip == "" {
			ip = r.Header.Get("X-Forwarded-For")
		}
		if ip == "" {
			// This will likely be [::1]:port in local tests
			ip = strings.Split(r.RemoteAddr, ":")[0]
		}

		// Create the JSON payload for the analysis service.
		// NOTE: In a real system, you might buffer the request body to include payload hashes/entropy,
		// and wrap the ResponseWriter to get execution time and response size post-request.
		// For this integration, we are analyzing based on pre-request metadata.
		requestDetailsJSON := fmt.Sprintf(`{
			"timestamp": "%s",
			"ip": "%s",
			"user_agent": "%s",
			"endpoint": "%s",
			"method": "%s",
			"status_code": 200, "response_size": %d, "execution_time": 0.1, "key_size": 2000, "payload_size": %d
		}`, time.Now().UTC().Format(time.RFC3339), ip, r.UserAgent(), r.URL.Path, r.Method, r.ContentLength, r.ContentLength)

		// --- 2. Get AI Analysis ---
		analysis, err := AnalyzeRequest(requestDetailsJSON)
		if err != nil {
			logrus.WithError(err).Warn("AI analysis request failed. Passing request through.")
			next.ServeHTTP(w, r)
			return
		}

		logrus.WithFields(logrus.Fields{
			"ip":          ip,
			"threat_type": analysis.ThreatType,
			"confidence":  analysis.Confidence,
			"action":      analysis.Action,
		}).Info("AI analysis complete")

		// --- 3. Take Action ---
		switch analysis.Action {
		case "THROTTLE":
			logrus.WithField("ip", ip).Warn("Throttling request due to suspicious activity.")
			http.Error(w, "Request throttled", http.StatusTooManyRequests)
			return
		case "DECEIVE":
			logrus.WithField("ip", ip).Warn("Serving deceptive response.")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"error": "invalid_cryptographic_key_format"}`) // Send a fake, plausible error
			return
		case "REDIRECT":
			logrus.WithField("ip", ip).Warn("Redirecting suspicious request to honeypot.")
			http.Redirect(w, r, "http://localhost:8080/honeypot", http.StatusFound)
			return
		case "PASS":
			fallthrough // Explicitly fall through to the default case
		default:
			logrus.Debug("Passing request to next handler.")
			next.ServeHTTP(w, r)
		}
	})
} 