package build

import (
	"fmt"

	mcfgv1 "github.com/openshift/api/machineconfiguration/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
)

// Event types for OCL processes
const (
	// Build lifecycle events
	EventBuildStarted      = "BuildStarted"
	EventBuildPreparing    = "BuildPreparing"
	EventBuildBuilding     = "BuildBuilding"
	EventBuildCompleted    = "BuildCompleted"
	EventBuildFailed       = "BuildFailed"
	EventBuildInterrupted  = "BuildInterrupted"
	EventBuildRetrying     = "BuildRetrying"
	EventBuildDeleted      = "BuildDeleted"

	// Image operations events
	EventImagePullStarted   = "ImagePullStarted"
	EventImagePullCompleted = "ImagePullCompleted"
	EventImagePullFailed    = "ImagePullFailed"
	EventImagePushStarted   = "ImagePushStarted"
	EventImagePushCompleted = "ImagePushCompleted"
	EventImagePushFailed    = "ImagePushFailed"

	// Job events
	EventJobCreated   = "JobCreated"
	EventJobStarted   = "JobStarted"
	EventJobCompleted = "JobCompleted"
	EventJobFailed    = "JobFailed"
	EventJobDeleted   = "JobDeleted"

	// Config events
	EventConfigChanged      = "ConfigChanged"
	EventConfigReconciling  = "ConfigReconciling"
	EventConfigReconciled   = "ConfigReconciled"
	EventRebuildRequested   = "RebuildRequested"

	// Cleanup events
	EventCleanupStarted   = "CleanupStarted"
	EventCleanupCompleted = "CleanupCompleted"
	EventCleanupFailed    = "CleanupFailed"

	// Degradation events
	EventBuildDegraded    = "BuildDegraded"
	EventBuildRecovered   = "BuildRecovered"
)

// OCLEventRecorder wraps the Kubernetes event recorder with OCL-specific event helpers
type OCLEventRecorder struct {
	recorder record.EventRecorder
}

// NewOCLEventRecorder creates a new OCL event recorder
func NewOCLEventRecorder(recorder record.EventRecorder) *OCLEventRecorder {
	return &OCLEventRecorder{
		recorder: recorder,
	}
}

// Build lifecycle event recording

// RecordBuildStarted records when a build starts
func (r *OCLEventRecorder) RecordBuildStarted(mosb *mcfgv1.MachineOSBuild, mosc *mcfgv1.MachineOSConfig) {
	r.recorder.Event(mosb, corev1.EventTypeNormal, EventBuildStarted,
		fmt.Sprintf("Started build for pool %q with config %q",
			mosc.Spec.MachineConfigPool.Name, mosb.Spec.MachineConfig.Name))
}

// RecordBuildPreparing records when a build is in the preparing phase
func (r *OCLEventRecorder) RecordBuildPreparing(mosb *mcfgv1.MachineOSBuild, message string) {
	r.recorder.Event(mosb, corev1.EventTypeNormal, EventBuildPreparing,
		fmt.Sprintf("Preparing build: %s", message))
}

// RecordBuildBuilding records when a build transitions to building state
func (r *OCLEventRecorder) RecordBuildBuilding(mosb *mcfgv1.MachineOSBuild) {
	r.recorder.Event(mosb, corev1.EventTypeNormal, EventBuildBuilding,
		"Build is now in progress")
}

// RecordBuildCompleted records when a build completes successfully
func (r *OCLEventRecorder) RecordBuildCompleted(mosb *mcfgv1.MachineOSBuild, imagePullspec string) {
	r.recorder.Event(mosb, corev1.EventTypeNormal, EventBuildCompleted,
		fmt.Sprintf("Build completed successfully, image: %s", imagePullspec))
}

// RecordBuildFailed records when a build fails
func (r *OCLEventRecorder) RecordBuildFailed(mosb *mcfgv1.MachineOSBuild, reason string) {
	r.recorder.Event(mosb, corev1.EventTypeWarning, EventBuildFailed,
		fmt.Sprintf("Build failed: %s", reason))
}

// RecordBuildInterrupted records when a build is interrupted
func (r *OCLEventRecorder) RecordBuildInterrupted(mosb *mcfgv1.MachineOSBuild, reason string) {
	r.recorder.Event(mosb, corev1.EventTypeWarning, EventBuildInterrupted,
		fmt.Sprintf("Build interrupted: %s", reason))
}

// RecordBuildRetrying records when a build is retrying
func (r *OCLEventRecorder) RecordBuildRetrying(mosb *mcfgv1.MachineOSBuild, attempt int) {
	r.recorder.Event(mosb, corev1.EventTypeNormal, EventBuildRetrying,
		fmt.Sprintf("Retrying build (attempt %d)", attempt))
}

// RecordBuildDeleted records when a build is deleted
func (r *OCLEventRecorder) RecordBuildDeleted(mosb *mcfgv1.MachineOSBuild, reason string) {
	r.recorder.Event(mosb, corev1.EventTypeNormal, EventBuildDeleted,
		fmt.Sprintf("Build deleted: %s", reason))
}

// Image operation event recording

// RecordImagePullStarted records when an image pull operation starts
func (r *OCLEventRecorder) RecordImagePullStarted(mosb *mcfgv1.MachineOSBuild, imageRef string) {
	r.recorder.Event(mosb, corev1.EventTypeNormal, EventImagePullStarted,
		fmt.Sprintf("Started pulling base image: %s", imageRef))
}

// RecordImagePullCompleted records when an image pull completes
func (r *OCLEventRecorder) RecordImagePullCompleted(mosb *mcfgv1.MachineOSBuild, imageRef string) {
	r.recorder.Event(mosb, corev1.EventTypeNormal, EventImagePullCompleted,
		fmt.Sprintf("Successfully pulled base image: %s", imageRef))
}

// RecordImagePullFailed records when an image pull fails
func (r *OCLEventRecorder) RecordImagePullFailed(mosb *mcfgv1.MachineOSBuild, imageRef, reason string) {
	r.recorder.Event(mosb, corev1.EventTypeWarning, EventImagePullFailed,
		fmt.Sprintf("Failed to pull base image %s: %s", imageRef, reason))
}

// RecordImagePushStarted records when an image push operation starts
func (r *OCLEventRecorder) RecordImagePushStarted(mosb *mcfgv1.MachineOSBuild, imageRef string) {
	r.recorder.Event(mosb, corev1.EventTypeNormal, EventImagePushStarted,
		fmt.Sprintf("Started pushing built image: %s", imageRef))
}

// RecordImagePushCompleted records when an image push completes
func (r *OCLEventRecorder) RecordImagePushCompleted(mosb *mcfgv1.MachineOSBuild, digestedPullspec string) {
	r.recorder.Event(mosb, corev1.EventTypeNormal, EventImagePushCompleted,
		fmt.Sprintf("Successfully pushed built image: %s", digestedPullspec))
}

// RecordImagePushFailed records when an image push fails
func (r *OCLEventRecorder) RecordImagePushFailed(mosb *mcfgv1.MachineOSBuild, imageRef, reason string) {
	r.recorder.Event(mosb, corev1.EventTypeWarning, EventImagePushFailed,
		fmt.Sprintf("Failed to push built image %s: %s", imageRef, reason))
}

// Job event recording

// RecordJobCreated records when a build job is created
func (r *OCLEventRecorder) RecordJobCreated(mosb *mcfgv1.MachineOSBuild, job *batchv1.Job) {
	r.recorder.Event(mosb, corev1.EventTypeNormal, EventJobCreated,
		fmt.Sprintf("Created build job: %s", job.Name))
}

// RecordJobStarted records when a build job starts
func (r *OCLEventRecorder) RecordJobStarted(mosb *mcfgv1.MachineOSBuild, job *batchv1.Job) {
	r.recorder.Event(mosb, corev1.EventTypeNormal, EventJobStarted,
		fmt.Sprintf("Build job started: %s", job.Name))
}

// RecordJobCompleted records when a build job completes
func (r *OCLEventRecorder) RecordJobCompleted(mosb *mcfgv1.MachineOSBuild, job *batchv1.Job) {
	r.recorder.Event(mosb, corev1.EventTypeNormal, EventJobCompleted,
		fmt.Sprintf("Build job completed: %s", job.Name))
}

// RecordJobFailed records when a build job fails
func (r *OCLEventRecorder) RecordJobFailed(mosb *mcfgv1.MachineOSBuild, job *batchv1.Job, reason string) {
	r.recorder.Event(mosb, corev1.EventTypeWarning, EventJobFailed,
		fmt.Sprintf("Build job failed %s: %s", job.Name, reason))
}

// RecordJobDeleted records when a build job is deleted
func (r *OCLEventRecorder) RecordJobDeleted(mosb *mcfgv1.MachineOSBuild, jobName string) {
	r.recorder.Event(mosb, corev1.EventTypeNormal, EventJobDeleted,
		fmt.Sprintf("Build job deleted: %s", jobName))
}

// Config event recording

// RecordConfigChanged records when a config changes triggering a new build
func (r *OCLEventRecorder) RecordConfigChanged(mosc *mcfgv1.MachineOSConfig, oldConfig, newConfig string) {
	r.recorder.Event(mosc, corev1.EventTypeNormal, EventConfigChanged,
		fmt.Sprintf("Configuration changed from %s to %s, triggering new build", oldConfig, newConfig))
}

// RecordConfigReconciling records when a config is being reconciled
func (r *OCLEventRecorder) RecordConfigReconciling(mosc *mcfgv1.MachineOSConfig) {
	r.recorder.Event(mosc, corev1.EventTypeNormal, EventConfigReconciling,
		"Reconciling MachineOSConfig")
}

// RecordConfigReconciled records when a config reconciliation completes
func (r *OCLEventRecorder) RecordConfigReconciled(mosc *mcfgv1.MachineOSConfig) {
	r.recorder.Event(mosc, corev1.EventTypeNormal, EventConfigReconciled,
		"MachineOSConfig reconciled successfully")
}

// RecordRebuildRequested records when a rebuild is requested
func (r *OCLEventRecorder) RecordRebuildRequested(mosc *mcfgv1.MachineOSConfig, reason string) {
	r.recorder.Event(mosc, corev1.EventTypeNormal, EventRebuildRequested,
		fmt.Sprintf("Rebuild requested: %s", reason))
}

// Cleanup event recording

// RecordCleanupStarted records when cleanup of build artifacts starts
func (r *OCLEventRecorder) RecordCleanupStarted(mosb *mcfgv1.MachineOSBuild) {
	r.recorder.Event(mosb, corev1.EventTypeNormal, EventCleanupStarted,
		"Started cleanup of build artifacts")
}

// RecordCleanupCompleted records when cleanup completes
func (r *OCLEventRecorder) RecordCleanupCompleted(mosb *mcfgv1.MachineOSBuild) {
	r.recorder.Event(mosb, corev1.EventTypeNormal, EventCleanupCompleted,
		"Cleanup of build artifacts completed")
}

// RecordCleanupFailed records when cleanup fails
func (r *OCLEventRecorder) RecordCleanupFailed(mosb *mcfgv1.MachineOSBuild, reason string) {
	r.recorder.Event(mosb, corev1.EventTypeWarning, EventCleanupFailed,
		fmt.Sprintf("Cleanup of build artifacts failed: %s", reason))
}

// Degradation event recording

// RecordBuildDegraded records when a build enters degraded state
func (r *OCLEventRecorder) RecordBuildDegraded(mosc *mcfgv1.MachineOSConfig, mcp *mcfgv1.MachineConfigPool, reason string) {
	r.recorder.Event(mosc, corev1.EventTypeWarning, EventBuildDegraded,
		fmt.Sprintf("Build for pool %s degraded: %s", mcp.Name, reason))
}

// RecordBuildRecovered records when a build recovers from degraded state
func (r *OCLEventRecorder) RecordBuildRecovered(mosc *mcfgv1.MachineOSConfig, mcp *mcfgv1.MachineConfigPool) {
	r.recorder.Event(mosc, corev1.EventTypeNormal, EventBuildRecovered,
		fmt.Sprintf("Build for pool %s recovered from degraded state", mcp.Name))
}

// MachineConfigPool events

// RecordPoolConfigChanged records when a pool's config changes
func (r *OCLEventRecorder) RecordPoolConfigChanged(mcp *mcfgv1.MachineConfigPool, oldConfig, newConfig string) {
	r.recorder.Event(mcp, corev1.EventTypeNormal, EventConfigChanged,
		fmt.Sprintf("Rendered config changed from %s to %s", oldConfig, newConfig))
}
