package security

import (
	"math"
	"sync"

	"github.com/sirupsen/logrus"
)

// AnomalyDetector implements a simple statistical anomaly detection system
// In a production system, this would be replaced with a trained machine learning model
type AnomalyDetector struct {
	// Statistical data about normal behavior
	mu                sync.RWMutex
	featureMeans      map[string]float64
	featureVariances  map[string]float64
	numSamples        int
	
	// Thresholds for anomaly detection
	entropyThreshold     float64
	interRequestTimeThreshold float64
	requestsPerMinuteThreshold int
	sequenceThreshold    float64
}

// NewAnomalyDetector creates a new anomaly detector
func NewAnomalyDetector() *AnomalyDetector {
	return &AnomalyDetector{
		featureMeans:     make(map[string]float64),
		featureVariances: make(map[string]float64),
		numSamples:       0,
		
		// Set initial thresholds
		entropyThreshold:     6.0,  // Lower entropy suggests non-random keys
		interRequestTimeThreshold: 0.1,  // Very fast requests are suspicious
		requestsPerMinuteThreshold: 20, // Too many requests per minute
		sequenceThreshold:    3.0,  // Mahalanobis distance threshold
	}
}

// Train updates the anomaly detector with normal traffic data
// This would be run periodically during a learning phase
func (d *AnomalyDetector) Train(features RequestFeatures) {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	d.numSamples++
	alpha := 1.0 / float64(d.numSamples)
	
	// Update means using incremental formula
	updateMean(d.featureMeans, "input_entropy", features.InputEntropy, alpha)
	updateMean(d.featureMeans, "inter_request_time", features.InterRequestTime, alpha)
	updateMean(d.featureMeans, "requests_per_minute", float64(features.RequestsPerMinute), alpha)
	updateMean(d.featureMeans, "operation_latency", features.OperationLatency, alpha)
	
	// Update variances using incremental formula (simplified)
	if d.numSamples > 1 {
		updateVariance(d.featureVariances, "input_entropy", features.InputEntropy, d.featureMeans["input_entropy"], alpha)
		updateVariance(d.featureVariances, "inter_request_time", features.InterRequestTime, d.featureMeans["inter_request_time"], alpha)
		updateVariance(d.featureVariances, "requests_per_minute", float64(features.RequestsPerMinute), d.featureMeans["requests_per_minute"], alpha)
		updateVariance(d.featureVariances, "operation_latency", features.OperationLatency, d.featureMeans["operation_latency"], alpha)
	}
	
	logrus.WithFields(logrus.Fields{
		"samples": d.numSamples,
		"means":   d.featureMeans,
	}).Debug("Anomaly detector trained with new sample")
}

// Detect checks if a request is anomalous
func (d *AnomalyDetector) Detect(features RequestFeatures) (bool, string, float64) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	// If we don't have enough samples for baseline, assume benign
	if d.numSamples < 10 {
		return false, "", 0.0
	}
	
	// Check statistical features
	if features.InputEntropy < d.entropyThreshold {
		return true, "LowEntropy", d.entropyThreshold - features.InputEntropy
	}
	
	// Check temporal features
	if features.InterRequestTime > 0 && features.InterRequestTime < d.interRequestTimeThreshold {
		return true, "RapidRequests", d.interRequestTimeThreshold - features.InterRequestTime
	}
	
	if features.RequestsPerMinute > d.requestsPerMinuteThreshold {
		return true, "HighFrequency", float64(features.RequestsPerMinute - d.requestsPerMinuteThreshold)
	}
	
	// Check for abnormal latency (potential side-channel attack)
	latencyScore := calculateZScore(features.OperationLatency, d.featureMeans["operation_latency"], d.featureVariances["operation_latency"])
	if math.Abs(latencyScore) > 3.0 {
		return true, "AbnormalLatency", math.Abs(latencyScore)
	}
	
	// In a real system, we would also check sequence patterns, calculate Mahalanobis distance, etc.
	
	// No anomaly detected
	return false, "", 0.0
}

// Helper functions

// updateMean updates the running mean for a feature
func updateMean(means map[string]float64, feature string, value float64, alpha float64) {
	oldMean, exists := means[feature]
	if !exists {
		means[feature] = value
		return
	}
	
	means[feature] = oldMean + alpha*(value-oldMean)
}

// updateVariance updates the running variance for a feature
func updateVariance(variances map[string]float64, feature string, value, mean, alpha float64) {
	oldVar, exists := variances[feature]
	if !exists {
		variances[feature] = 0.1 // Initialize with small non-zero value
		return
	}
	
	delta := value - mean
	variances[feature] = (1-alpha)*oldVar + alpha*delta*delta
}

// calculateZScore calculates the z-score of a value
func calculateZScore(value, mean, variance float64) float64 {
	if variance <= 0 {
		return 0
	}
	return (value - mean) / math.Sqrt(variance)
} 