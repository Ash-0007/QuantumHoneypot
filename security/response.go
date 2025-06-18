package security

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

// ThreatType represents a classified threat
type ThreatType string

const (
	ThreatRecon        ThreatType = "Reconnaissance"
	ThreatSideChannel  ThreatType = "Side-Channel Probe"
	ThreatImplementation ThreatType = "Implementation Exploit"
	ThreatUnknown      ThreatType = "Unknown"
)

// ThreatLevel represents the severity of a threat
type ThreatLevel int

const (
	ThreatLevelLow     ThreatLevel = 1
	ThreatLevelMedium  ThreatLevel = 2
	ThreatLevelHigh    ThreatLevel = 3
	ThreatLevelCritical ThreatLevel = 4
)

// ActionType represents a response action
type ActionType string

const (
	ActionPass     ActionType = "Pass"
	ActionThrottle ActionType = "Throttle"
	ActionDeceive  ActionType = "Deceive"
	ActionRedirect ActionType = "Redirect"
	ActionBlock    ActionType = "Block"
)

// Threat contains details about a classified threat
type Threat struct {
	IP          string     `json:"ip"`
	Type        ThreatType `json:"type"`
	Level       ThreatLevel `json:"level"`
	Score       float64    `json:"score"`
	Description string     `json:"description"`
	Timestamp   time.Time  `json:"timestamp"`
	Features    RequestFeatures `json:"features"`
}

// ResponseEngine decides on and applies appropriate responses to threats
type ResponseEngine struct {
	// Map of IP to threat history
	mu              sync.RWMutex
	threatHistory   map[string][]Threat
	
	// Throttling limits for IPs
	throttlingMap   map[string]*rate.Limiter
	
	// Deception settings
	deceptionEnabled bool
	
	// Honeypot settings
	honeypotIP      string
	honeypotEnabled bool
}

// NewResponseEngine creates a new response engine
func NewResponseEngine() *ResponseEngine {
	return &ResponseEngine{
		threatHistory:    make(map[string][]Threat),
		throttlingMap:    make(map[string]*rate.Limiter),
		deceptionEnabled: true,
		honeypotEnabled:  true,
		honeypotIP:       "10.10.10.10", // In a real system, this would be a real honeypot server
	}
}

// ClassifyThreat determines the type and severity of a detected anomaly
func (r *ResponseEngine) ClassifyThreat(features RequestFeatures, anomalyType string, score float64) Threat {
	var threatType ThreatType
	var threatLevel ThreatLevel
	var description string
	
	// Classify based on anomaly type
	switch anomalyType {
	case "LowEntropy":
		threatType = ThreatImplementation
		description = "Suspiciously low entropy in cryptographic input"
		if score > 2.0 {
			threatLevel = ThreatLevelHigh
		} else {
			threatLevel = ThreatLevelMedium
		}
		
	case "RapidRequests":
		threatType = ThreatRecon
		description = "Abnormally rapid sequence of requests"
		if features.RequestsPerMinute > 1000 {
			threatLevel = ThreatLevelHigh
		} else {
			threatLevel = ThreatLevelMedium
		}
		
	case "HighFrequency":
		threatType = ThreatRecon
		description = "High frequency API usage"
		threatLevel = ThreatLevelMedium
		
	case "AbnormalLatency":
		threatType = ThreatSideChannel
		description = "Abnormal operation timing, possible side-channel probe"
		if score > 5.0 {
			threatLevel = ThreatLevelCritical
		} else {
			threatLevel = ThreatLevelHigh
		}
		
	default:
		threatType = ThreatUnknown
		description = "Unknown anomaly type"
		threatLevel = ThreatLevelLow
	}
	
	// Create threat record
	threat := Threat{
		IP:          features.ClientIP,
		Type:        threatType,
		Level:       threatLevel,
		Score:       score,
		Description: description,
		Timestamp:   time.Now(),
		Features:    features,
	}
	
	// Update threat history
	r.mu.Lock()
	threats := r.threatHistory[features.ClientIP]
	if len(threats) > 99 {
		// Keep last 100 threats
		threats = threats[1:]
	}
	threats = append(threats, threat)
	r.threatHistory[features.ClientIP] = threats
	r.mu.Unlock()
	
	return threat
}

// DecideAction determines the appropriate response to a threat
// This is a simplified version of what would normally be a reinforcement learning agent
func (r *ResponseEngine) DecideAction(threat Threat) ActionType {
	// Get threat history for this IP
	r.mu.RLock()
	history := r.threatHistory[threat.IP]
	historyLength := len(history)
	r.mu.RUnlock()
	
	// Simple rule-based decision making
	// In a real system, this would be a trained RL agent

	// For critical threats, always redirect to honeypot
	if threat.Level == ThreatLevelCritical && r.honeypotEnabled {
		return ActionRedirect
	}
	
	// For high-level threats or repeat offenders, use deception
	if threat.Level == ThreatLevelHigh || historyLength > 5 {
		if r.deceptionEnabled {
			return ActionDeceive
		}
		return ActionThrottle
	}
	
	// For medium-level threats, throttle
	if threat.Level == ThreatLevelMedium {
		return ActionThrottle
	}
	
	// For low-level threats, just pass
	return ActionPass
}

// ApplyAction executes the selected action
func (r *ResponseEngine) ApplyAction(action ActionType, clientIP string) {
	switch action {
	case ActionThrottle:
		// Create or update rate limiter for this IP
		r.mu.Lock()
		if _, exists := r.throttlingMap[clientIP]; !exists {
			// Initial limit: 5 requests per second with burst of 3
			r.throttlingMap[clientIP] = rate.NewLimiter(5, 3)
		} else {
			// Further restrict existing limiter
			current := r.throttlingMap[clientIP]
			// Cut rate in half each time
			newRate := current.Limit() / 2
			if newRate < 0.1 {
				newRate = 0.1 // Don't go below 1 request per 10 seconds
			}
			r.throttlingMap[clientIP] = rate.NewLimiter(newRate, 1)
		}
		r.mu.Unlock()
		
	case ActionRedirect:
		// In a real system, we would configure the proxy to redirect to honeypot
		logrus.WithFields(logrus.Fields{
			"action": "redirect",
			"ip":     clientIP,
			"target": r.honeypotIP,
		}).Info("Redirecting suspicious traffic to honeypot")
		
	case ActionDeceive:
		// No special setup needed, the handler will use this flag
		logrus.WithFields(logrus.Fields{
			"action": "deceive",
			"ip":     clientIP,
		}).Info("Activating deception for suspicious client")
		
	case ActionBlock:
		// In a real system, we would add to a block list
		logrus.WithFields(logrus.Fields{
			"action": "block",
			"ip":     clientIP,
		}).Info("Blocking malicious client")
	}
}

// ShouldThrottle checks if a request from this IP should be throttled
func (r *ResponseEngine) ShouldThrottle(clientIP string) bool {
	r.mu.RLock()
	limiter, exists := r.throttlingMap[clientIP]
	r.mu.RUnlock()
	
	if !exists {
		return false
	}
	
	return !limiter.Allow()
}

// ShouldDeceive checks if deception should be used for this IP
func (r *ResponseEngine) ShouldDeceive(clientIP string) bool {
	if !r.deceptionEnabled {
		return false
	}
	
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	history, exists := r.threatHistory[clientIP]
	if !exists || len(history) == 0 {
		return false
	}
	
	// Check if any recent action was to deceive
	for i := len(history) - 1; i >= 0 && i >= len(history)-5; i-- {
		// We don't store the action in the threat history, but in a real system we would
		// For now, assume high-level threats get deception
		if history[i].Level >= ThreatLevelHigh {
			return true
		}
	}
	
	return false
}

// ShouldRedirect checks if request should be redirected to honeypot
func (r *ResponseEngine) ShouldRedirect(clientIP string) bool {
	if !r.honeypotEnabled {
		return false
	}
	
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	history, exists := r.threatHistory[clientIP]
	if !exists || len(history) == 0 {
		return false
	}
	
	// Check if any recent threat was critical level
	for i := len(history) - 1; i >= 0 && i >= len(history)-3; i-- {
		if history[i].Level == ThreatLevelCritical {
			return true
		}
	}
	
	return false
} 