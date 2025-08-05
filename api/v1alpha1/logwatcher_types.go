package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LogWatcherSpec defines the desired state of LogWatcher
type LogWatcherSpec struct {
	// PodNamespace is the namespace where pods should be monitored
	PodNamespace string `json:"podNamespace"`

	// PodLabelSelector is the label selector to identify pods to monitor
	PodLabelSelector map[string]string `json:"podLabelSelector"`

	// MatchPattern is the regex pattern to match in logs
	MatchPattern string `json:"matchPattern"`

	// Actions define what to do when pattern is matched
	Actions LogWatcherActions `json:"actions"`

	// Alerting configuration
	Alerting *AlertingConfig `json:"alerting,omitempty"`

	// Scaling configuration
	Scaling *ScalingConfig `json:"scaling,omitempty"`

	// Annotations to add to pods when events occur
	Annotations map[string]string `json:"annotations,omitempty"`

	// Metrics configuration
	Metrics *MetricsConfig `json:"metrics,omitempty"`

	// ReconcileInterval is how often to check logs (in seconds)
	ReconcileInterval int32 `json:"reconcileInterval,omitempty"`

	// TailLines is the number of log lines to check
	TailLines int32 `json:"tailLines,omitempty"`
}

// LogWatcherActions defines what actions to take when pattern is matched
type LogWatcherActions struct {
	// RestartPod will restart the pod when pattern is matched
	RestartPod bool `json:"restartPod,omitempty"`

	// CreateJob will create a Kubernetes Job when pattern is matched
	CreateJob *JobConfig `json:"createJob,omitempty"`

	// RollingUpdate will trigger a rolling update of the deployment
	RollingUpdate *RollingUpdateConfig `json:"rollingUpdate,omitempty"`

	// Cleanup will clean up resources after a timeout
	Cleanup *CleanupConfig `json:"cleanup,omitempty"`
}

// JobConfig defines configuration for creating Kubernetes Jobs
type JobConfig struct {
	// JobTemplate is the job template to use
	JobTemplate string `json:"jobTemplate"`

	// Timeout is the timeout for the job (in seconds)
	Timeout int32 `json:"timeout,omitempty"`

	// Retries is the number of retries for the job
	Retries int32 `json:"retries,omitempty"`
}

// RollingUpdateConfig defines configuration for rolling updates
type RollingUpdateConfig struct {
	// DeploymentName is the name of the deployment to update
	DeploymentName string `json:"deploymentName"`

	// MaxUnavailable is the maximum number of unavailable pods
	MaxUnavailable int32 `json:"maxUnavailable,omitempty"`

	// MaxSurge is the maximum number of pods that can be created above desired replicas
	MaxSurge int32 `json:"maxSurge,omitempty"`
}

// CleanupConfig defines configuration for cleanup operations
type CleanupConfig struct {
	// Timeout is the timeout before cleanup (in seconds)
	Timeout int32 `json:"timeout"`

	// Resources to clean up
	Resources []string `json:"resources"`
}

// AlertingConfig defines alerting configuration
type AlertingConfig struct {
	// Slack configuration
	Slack *SlackConfig `json:"slack,omitempty"`

	// Email configuration
	Email *EmailConfig `json:"email,omitempty"`

	// Webhook configuration
	Webhook *WebhookConfig `json:"webhook,omitempty"`
}

// SlackConfig defines Slack alerting configuration
type SlackConfig struct {
	// WebhookURL is the Slack webhook URL
	WebhookURL string `json:"webhookURL"`

	// Channel is the Slack channel to send alerts to
	Channel string `json:"channel,omitempty"`

	// Username is the username for the Slack bot
	Username string `json:"username,omitempty"`
}

// EmailConfig defines email alerting configuration
type EmailConfig struct {
	// SMTPHost is the SMTP host
	SMTPHost string `json:"smtpHost"`

	// SMTPPort is the SMTP port
	SMTPPort int32 `json:"smtpPort"`

	// Username is the SMTP username
	Username string `json:"username"`

	// Password is the SMTP password
	Password string `json:"password"`

	// From is the sender email address
	From string `json:"from"`

	// To is the recipient email address
	To string `json:"to"`

	// Subject is the email subject
	Subject string `json:"subject,omitempty"`
}

// WebhookConfig defines webhook alerting configuration
type WebhookConfig struct {
	// URL is the webhook URL
	URL string `json:"url"`

	// Headers are the HTTP headers to send
	Headers map[string]string `json:"headers,omitempty"`

	// Method is the HTTP method to use
	Method string `json:"method,omitempty"`
}

// ScalingConfig defines scaling configuration
type ScalingConfig struct {
	// MinReplicas is the minimum number of replicas
	MinReplicas int32 `json:"minReplicas"`

	// MaxReplicas is the maximum number of replicas
	MaxReplicas int32 `json:"maxReplicas"`

	// ScaleUpThreshold is the log frequency threshold to scale up
	ScaleUpThreshold int32 `json:"scaleUpThreshold"`

	// ScaleDownThreshold is the log frequency threshold to scale down
	ScaleDownThreshold int32 `json:"scaleDownThreshold"`

	// DeploymentName is the name of the deployment to scale
	DeploymentName string `json:"deploymentName"`
}

// MetricsConfig defines metrics configuration
type MetricsConfig struct {
	// Enabled enables metrics export
	Enabled bool `json:"enabled"`

	// Port is the port to expose metrics on
	Port int32 `json:"port,omitempty"`

	// Path is the path to expose metrics on
	Path string `json:"path,omitempty"`
}

// LogWatcherStatus defines the observed state of LogWatcher
type LogWatcherStatus struct {
	// Conditions represent the latest available observations of a LogWatcher's current state
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// LastReconcileTime is the last time the controller reconciled
	LastReconcileTime metav1.Time `json:"lastReconcileTime,omitempty"`

	// MatchesCount is the number of pattern matches found
	MatchesCount int32 `json:"matchesCount,omitempty"`

	// PodsRestarted is the number of pods restarted
	PodsRestarted int32 `json:"podsRestarted,omitempty"`

	// AlertsSent is the number of alerts sent
	AlertsSent int32 `json:"alertsSent,omitempty"`

	// JobsCreated is the number of jobs created
	JobsCreated int32 `json:"jobsCreated,omitempty"`

	// ScalingEvents is the number of scaling events
	ScalingEvents int32 `json:"scalingEvents,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
//+kubebuilder:printcolumn:name="Matches",type="integer",JSONPath=".status.matchesCount"
//+kubebuilder:printcolumn:name="Restarts",type="integer",JSONPath=".status.podsRestarted"

// LogWatcher is the Schema for the logwatchers API
type LogWatcher struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LogWatcherSpec   `json:"spec,omitempty"`
	Status LogWatcherStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LogWatcherList contains a list of LogWatcher
type LogWatcherList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LogWatcher `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LogWatcher{}, &LogWatcherList{})
} 