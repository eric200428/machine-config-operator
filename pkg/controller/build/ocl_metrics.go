package build

import (
	"fmt"
	"time"

	ctrlcommon "github.com/openshift/machine-config-operator/pkg/controller/common"
	"github.com/prometheus/client_golang/prometheus"
)

// OCL Build Metrics for tracking On-Cluster Layering processes
var (
	// oclBuildState tracks the current state of OCL builds per pool
	oclBuildState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ocl_build_state",
			Help: "Current state of OCL build for a pool (0=none, 1=pending, 2=building, 3=succeeded, 4=failed, 5=interrupted)",
		}, []string{"pool", "build_name", "state"})

	// oclBuildDuration tracks how long OCL builds take to complete
	oclBuildDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocl_build_duration_seconds",
			Help:    "Duration of OCL build processes in seconds",
			Buckets: []float64{60, 180, 300, 600, 900, 1200, 1800, 2400, 3000, 3600}, // 1m to 1h
		}, []string{"pool", "state"})

	// oclBuildStartTime tracks when a build started
	oclBuildStartTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ocl_build_start_timestamp_seconds",
			Help: "Timestamp when OCL build started",
		}, []string{"pool", "build_name"})

	// oclBuildEndTime tracks when a build completed/failed
	oclBuildEndTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ocl_build_end_timestamp_seconds",
			Help: "Timestamp when OCL build completed or failed",
		}, []string{"pool", "build_name", "state"})

	// oclBuildTotal counts total number of builds by state
	oclBuildTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ocl_build_total",
			Help: "Total number of OCL builds by final state",
		}, []string{"pool", "state"})

	// oclImagePullState tracks image pull operations
	oclImagePullState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ocl_image_pull_state",
			Help: "State of OCL image pull operations (0=none, 1=pulling, 2=succeeded, 3=failed)",
		}, []string{"pool", "build_name", "state"})

	// oclImagePullDuration tracks image pull duration
	oclImagePullDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ocl_image_pull_duration_seconds",
			Help:    "Duration of OCL image pull operations in seconds",
			Buckets: []float64{10, 30, 60, 120, 180, 300, 600}, // 10s to 10m
		}, []string{"pool", "state"})

	// oclBuildJobState tracks the state of build jobs
	oclBuildJobState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ocl_build_job_state",
			Help: "State of OCL build job (0=none, 1=active, 2=succeeded, 3=failed)",
		}, []string{"pool", "build_name", "job_name", "state"})

	// oclConfigChangeTotal counts config changes triggering builds
	oclConfigChangeTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ocl_config_change_total",
			Help: "Total number of config changes triggering OCL builds",
		}, []string{"pool"})

	// oclBuildRetries tracks build retry attempts
	oclBuildRetries = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ocl_build_retries_total",
			Help: "Total number of OCL build retry attempts",
		}, []string{"pool", "build_name"})

	// oclLayeredNodesCount tracks number of nodes using layered images
	oclLayeredNodesCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ocl_layered_nodes_count",
			Help: "Number of nodes currently using OCL layered images",
		}, []string{"pool"})

	// oclBuildArtifactSize tracks the size of build artifacts
	oclBuildArtifactSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ocl_build_artifact_size_bytes",
			Help: "Size of OCL build artifacts in bytes",
		}, []string{"pool", "build_name", "artifact_type"})
)

// Build state constants for consistent labeling
const (
	StateNone        = "none"
	StatePending     = "pending"
	StateBuilding    = "building"
	StateSucceeded   = "succeeded"
	StateFailed      = "failed"
	StateInterrupted = "interrupted"
	StatePulling     = "pulling"
)

// RegisterOCLMetrics registers all OCL-related Prometheus metrics
func RegisterOCLMetrics() error {
	err := ctrlcommon.RegisterMetrics([]prometheus.Collector{
		oclBuildState,
		oclBuildDuration,
		oclBuildStartTime,
		oclBuildEndTime,
		oclBuildTotal,
		oclImagePullState,
		oclImagePullDuration,
		oclBuildJobState,
		oclConfigChangeTotal,
		oclBuildRetries,
		oclLayeredNodesCount,
		oclBuildArtifactSize,
	})

	if err != nil {
		return fmt.Errorf("could not register OCL metrics: %w", err)
	}

	// Initialize GaugeVecs to ensure metrics are accessible even without values
	oclBuildState.WithLabelValues("init", "init", StateNone).Set(0)
	oclImagePullState.WithLabelValues("init", "init", StateNone).Set(0)
	oclBuildJobState.WithLabelValues("init", "init", "init", StateNone).Set(0)
	oclLayeredNodesCount.WithLabelValues("init").Set(0)

	return nil
}

// RecordBuildStarted records when a build starts
func RecordBuildStarted(pool, buildName string) {
	now := float64(time.Now().Unix())

	// Clear previous states for this build
	oclBuildState.DeletePartialMatch(prometheus.Labels{"pool": pool, "build_name": buildName})

	// Set new state
	oclBuildState.WithLabelValues(pool, buildName, StatePending).Set(1)
	oclBuildStartTime.WithLabelValues(pool, buildName).Set(now)
}

// RecordBuildBuilding records when a build transitions to building state
func RecordBuildBuilding(pool, buildName string) {
	// Clear pending state
	oclBuildState.DeletePartialMatch(prometheus.Labels{"pool": pool, "build_name": buildName})

	// Set building state
	oclBuildState.WithLabelValues(pool, buildName, StateBuilding).Set(1)
}

// RecordBuildCompleted records when a build completes successfully
func RecordBuildCompleted(pool, buildName string, startTime time.Time) {
	now := time.Now()
	duration := now.Sub(startTime).Seconds()

	// Clear previous states
	oclBuildState.DeletePartialMatch(prometheus.Labels{"pool": pool, "build_name": buildName})

	// Set succeeded state
	oclBuildState.WithLabelValues(pool, buildName, StateSucceeded).Set(1)
	oclBuildEndTime.WithLabelValues(pool, buildName, StateSucceeded).Set(float64(now.Unix()))
	oclBuildDuration.WithLabelValues(pool, StateSucceeded).Observe(duration)
	oclBuildTotal.WithLabelValues(pool, StateSucceeded).Inc()
}

// RecordBuildFailed records when a build fails
func RecordBuildFailed(pool, buildName string, startTime time.Time) {
	now := time.Now()
	duration := now.Sub(startTime).Seconds()

	// Clear previous states
	oclBuildState.DeletePartialMatch(prometheus.Labels{"pool": pool, "build_name": buildName})

	// Set failed state
	oclBuildState.WithLabelValues(pool, buildName, StateFailed).Set(1)
	oclBuildEndTime.WithLabelValues(pool, buildName, StateFailed).Set(float64(now.Unix()))
	oclBuildDuration.WithLabelValues(pool, StateFailed).Observe(duration)
	oclBuildTotal.WithLabelValues(pool, StateFailed).Inc()
}

// RecordBuildInterrupted records when a build is interrupted
func RecordBuildInterrupted(pool, buildName string) {
	// Clear previous states
	oclBuildState.DeletePartialMatch(prometheus.Labels{"pool": pool, "build_name": buildName})

	// Set interrupted state
	oclBuildState.WithLabelValues(pool, buildName, StateInterrupted).Set(1)
	oclBuildEndTime.WithLabelValues(pool, buildName, StateInterrupted).Set(float64(time.Now().Unix()))
	oclBuildTotal.WithLabelValues(pool, StateInterrupted).Inc()
}

// RecordImagePullStarted records when an image pull starts
func RecordImagePullStarted(pool, buildName string) {
	oclImagePullState.DeletePartialMatch(prometheus.Labels{"pool": pool, "build_name": buildName})
	oclImagePullState.WithLabelValues(pool, buildName, StatePulling).Set(1)
}

// RecordImagePullCompleted records when an image pull completes
func RecordImagePullCompleted(pool, buildName string, duration time.Duration) {
	oclImagePullState.DeletePartialMatch(prometheus.Labels{"pool": pool, "build_name": buildName})
	oclImagePullState.WithLabelValues(pool, buildName, StateSucceeded).Set(1)
	oclImagePullDuration.WithLabelValues(pool, StateSucceeded).Observe(duration.Seconds())
}

// RecordImagePullFailed records when an image pull fails
func RecordImagePullFailed(pool, buildName string, duration time.Duration) {
	oclImagePullState.DeletePartialMatch(prometheus.Labels{"pool": pool, "build_name": buildName})
	oclImagePullState.WithLabelValues(pool, buildName, StateFailed).Set(1)
	oclImagePullDuration.WithLabelValues(pool, StateFailed).Observe(duration.Seconds())
}

// RecordBuildJobState records the state of a build job
func RecordBuildJobState(pool, buildName, jobName, state string) {
	oclBuildJobState.DeletePartialMatch(prometheus.Labels{"pool": pool, "build_name": buildName, "job_name": jobName})

	stateValue := 0.0
	switch state {
	case "active":
		stateValue = 1.0
	case StateSucceeded:
		stateValue = 2.0
	case StateFailed:
		stateValue = 3.0
	}

	oclBuildJobState.WithLabelValues(pool, buildName, jobName, state).Set(stateValue)
}

// RecordConfigChange records when a config change triggers a build
func RecordConfigChange(pool string) {
	oclConfigChangeTotal.WithLabelValues(pool).Inc()
}

// RecordBuildRetry records a build retry attempt
func RecordBuildRetry(pool, buildName string) {
	oclBuildRetries.WithLabelValues(pool, buildName).Inc()
}

// UpdateLayeredNodesCount updates the count of nodes using layered images
func UpdateLayeredNodesCount(pool string, count int) {
	oclLayeredNodesCount.WithLabelValues(pool).Set(float64(count))
}

// RecordBuildArtifactSize records the size of a build artifact
func RecordBuildArtifactSize(pool, buildName, artifactType string, sizeBytes int64) {
	oclBuildArtifactSize.WithLabelValues(pool, buildName, artifactType).Set(float64(sizeBytes))
}
