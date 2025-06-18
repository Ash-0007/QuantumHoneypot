package benchmark

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"pqcd/crypto"
)

// OperationStats contains aggregated stats for a crypto operation
type OperationStats struct {
	Operation   string         `json:"operation"`
	Algorithm   crypto.Algorithm `json:"algorithm"`
	Count       int            `json:"count"`
	AvgLatency  float64        `json:"avg_latency_us"`
	MinLatency  float64        `json:"min_latency_us"`
	MaxLatency  float64        `json:"max_latency_us"`
	AvgInputSize int           `json:"avg_input_size_bytes"`
	AvgOutputSize int          `json:"avg_output_size_bytes"`
	SuccessRate  float64       `json:"success_rate"`
}

// MetricsCollector collects and reports performance metrics
type MetricsCollector struct {
	mutex  sync.RWMutex
	stats  map[string]*OperationStats // Key is "algorithm:operation"
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		stats: make(map[string]*OperationStats),
	}
}

// RecordOperation records a single operation
func (m *MetricsCollector) RecordOperation(algorithm crypto.Algorithm, operation string, duration time.Duration, inputSize, outputSize int, success bool) {
	key := string(algorithm) + ":" + operation
	
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	// Get or create stats for this operation
	stats, exists := m.stats[key]
	if !exists {
		stats = &OperationStats{
			Operation: operation,
			Algorithm: algorithm,
			MinLatency: float64(duration.Microseconds()),
			MaxLatency: float64(duration.Microseconds()),
		}
		m.stats[key] = stats
	}
	
	// Update stats
	stats.Count++
	
	latencyUs := float64(duration.Microseconds())
	
	// Update running average
	stats.AvgLatency = ((stats.AvgLatency * float64(stats.Count-1)) + latencyUs) / float64(stats.Count)
	stats.AvgInputSize = ((stats.AvgInputSize * (stats.Count-1)) + inputSize) / stats.Count
	stats.AvgOutputSize = ((stats.AvgOutputSize * (stats.Count-1)) + outputSize) / stats.Count
	
	// Update min/max
	if latencyUs < stats.MinLatency {
		stats.MinLatency = latencyUs
	}
	if latencyUs > stats.MaxLatency {
		stats.MaxLatency = latencyUs
	}
	
	// Update success rate
	if success {
		stats.SuccessRate = ((stats.SuccessRate * float64(stats.Count-1)) + 1) / float64(stats.Count)
	} else {
		stats.SuccessRate = (stats.SuccessRate * float64(stats.Count-1)) / float64(stats.Count)
	}
	
	// Log detailed metrics for this operation
	logrus.WithFields(logrus.Fields{
		"algorithm":    algorithm,
		"operation":    operation,
		"latency_us":   latencyUs,
		"input_bytes":  inputSize,
		"output_bytes": outputSize,
		"success":      success,
	}).Debug("Operation recorded")
}

// GetAllStats returns all collected stats
func (m *MetricsCollector) GetAllStats() []OperationStats {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	stats := make([]OperationStats, 0, len(m.stats))
	for _, stat := range m.stats {
		stats = append(stats, *stat)
	}
	
	return stats
}

// HandleMetrics returns an HTTP handler for metrics endpoint
func (m *MetricsCollector) HandleMetrics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats := m.GetAllStats()
		
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(stats); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.WithError(err).Error("Failed to encode metrics")
		}
	}
} 