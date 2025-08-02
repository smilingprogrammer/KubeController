# Kubernetes LogWatcher Controller

A comprehensive Kubernetes controller that monitors pod logs and performs automated actions based on log patterns. This controller provides advanced features for pod monitoring, alerting, scaling, job automation, and observability.

## Features

### 🚀 Pod/Container Monitoring
- **Log Pattern Detection**: Monitor pods for specific log patterns using regex
- **Auto-restart Pods**: Automatically restart pods when log patterns are detected
- **Pod Annotations**: Add custom annotations to pods when events occur
- **Configurable Monitoring**: Set custom reconcile intervals and log tail lines

### 📢 Alerting System
- **Slack Integration**: Send alerts to Slack channels with custom messages
- **Email Alerts**: Configure SMTP settings for email notifications
- **Webhook Support**: Send HTTP requests to external systems
- **Customizable Messages**: Include pod details, logs, and timestamps

### 📈 Auto-scaling
- **Log-based Scaling**: Scale deployments based on log frequency
- **Configurable Thresholds**: Set scale-up and scale-down thresholds
- **Min/Max Replicas**: Define scaling boundaries
- **Deployment Targeting**: Scale specific deployments

### 🔧 Job/Lifecycle Automation
- **Kubernetes Job Creation**: Create jobs when log patterns are detected
- **Job Templates**: Use custom job templates with retry logic
- **Timeout Management**: Set job timeouts and retry limits
- **Automatic Cleanup**: Clean up old jobs and pods after timeouts

### 🔄 Rolling Updates
- **Conditional Updates**: Trigger rolling updates on specific conditions
- **Deployment Targeting**: Update specific deployments
- **Configurable Strategy**: Set max unavailable and surge values

### 📊 Observability
- **Prometheus Metrics**: Export comprehensive metrics for monitoring
- **Custom Metrics**: Track matches, restarts, alerts, jobs, and scaling events
- **Health Checks**: Built-in health and readiness probes
- **Structured Logging**: Detailed logging with structured fields

## Installation

### Prerequisites
- Kubernetes cluster (1.19+)
- kubectl configured
- Go 1.21+ (for building from source)

### Quick Start

1. **Clone the repository**:
```bash
git clone <repository-url>
cd KubeController
```

2. **Install dependencies**:
```bash
go mod tidy
```

3. **Build the controller**:
```bash
go build -o bin/manager main.go
```

4. **Deploy to cluster**:
```bash
# Create namespace
kubectl create namespace logwatcher-system

# Apply CRDs
kubectl apply -f config/crd/bases/

# Deploy the controller
kubectl apply -f config/default/
```

## Usage Examples

### Basic Log Monitoring

```yaml
apiVersion: monitoring.kubecontroller.com/v1alpha1
kind: LogWatcher
metadata:
  name: error-monitor
  namespace: default
spec:
  podNamespace: default
  podLabelSelector:
    app: myapp
  matchPattern: "ERROR|FATAL|Exception"
  actions:
    restartPod: true
  reconcileInterval: 30
  tailLines: 100
```

### Advanced Monitoring with Alerts

```yaml
apiVersion: monitoring.kubecontroller.com/v1alpha1
kind: LogWatcher
metadata:
  name: production-monitor
  namespace: default
spec:
  podNamespace: production
  podLabelSelector:
    app: production-app
  matchPattern: "OutOfMemoryError|ConnectionTimeout"
  actions:
    restartPod: true
    createJob:
      jobTemplate: "echo 'Critical error detected' && kubectl get pods -n production"
      timeout: 300
      retries: 3
    rollingUpdate:
      deploymentName: production-app
      maxUnavailable: 1
      maxSurge: 1
  alerting:
    slack:
      webhookURL: "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"
      channel: "#alerts"
      username: "LogWatcher"
    email:
      smtpHost: "smtp.gmail.com"
      smtpPort: 587
      username: "your-email@gmail.com"
      password: "your-password"
      from: "alerts@company.com"
      to: "ops@company.com"
      subject: "Critical Error Detected"
  scaling:
    minReplicas: 2
    maxReplicas: 10
    scaleUpThreshold: 5
    scaleDownThreshold: 1
    deploymentName: production-app
  annotations:
    logwatcher.io/error-detected: "true"
    logwatcher.io/last-error-time: "{{.Timestamp}}"
  metrics:
    enabled: true
    port: 8080
    path: "/metrics"
  reconcileInterval: 60
  tailLines: 200
```

### Job Automation Example

```yaml
apiVersion: monitoring.kubecontroller.com/v1alpha1
kind: LogWatcher
metadata:
  name: backup-trigger
  namespace: default
spec:
  podNamespace: database
  podLabelSelector:
    app: database
  matchPattern: "backup.*completed|backup.*failed"
  actions:
    createJob:
      jobTemplate: |
        #!/bin/bash
        echo "Database backup status detected"
        kubectl logs -n database -l app=database --tail=50
        # Additional backup verification logic
      timeout: 600
      retries: 2
    cleanup:
      timeout: 3600
      resources:
        - "jobs"
        - "pods"
  alerting:
    webhook:
      url: "https://api.company.com/backup-webhook"
      method: "POST"
      headers:
        Authorization: "Bearer your-token"
        Content-Type: "application/json"
```

### Scaling Based on Log Frequency

```yaml
apiVersion: monitoring.kubecontroller.com/v1alpha1
kind: LogWatcher
metadata:
  name: load-based-scaling
  namespace: default
spec:
  podNamespace: webapp
  podLabelSelector:
    app: webapp
  matchPattern: "request.*processed|error.*rate.*high"
  scaling:
    minReplicas: 3
    maxReplicas: 20
    scaleUpThreshold: 10
    scaleDownThreshold: 2
    deploymentName: webapp-deployment
  metrics:
    enabled: true
```

## Configuration Reference

### LogWatcherSpec

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| `podNamespace` | string | Namespace to monitor pods in | Yes |
| `podLabelSelector` | map[string]string | Label selector for pods | Yes |
| `matchPattern` | string | Regex pattern to match in logs | Yes |
| `actions` | LogWatcherActions | Actions to take when pattern matches | Yes |
| `alerting` | AlertingConfig | Alerting configuration | No |
| `scaling` | ScalingConfig | Auto-scaling configuration | No |
| `annotations` | map[string]string | Annotations to add to pods | No |
| `metrics` | MetricsConfig | Metrics configuration | No |
| `reconcileInterval` | int32 | Reconcile interval in seconds | No |
| `tailLines` | int32 | Number of log lines to check | No |

### LogWatcherActions

| Field | Type | Description |
|-------|------|-------------|
| `restartPod` | bool | Restart pod when pattern matches |
| `createJob` | JobConfig | Create Kubernetes Job |
| `rollingUpdate` | RollingUpdateConfig | Trigger rolling update |
| `cleanup` | CleanupConfig | Cleanup configuration |

### AlertingConfig

| Field | Type | Description |
|-------|------|-------------|
| `slack` | SlackConfig | Slack alerting configuration |
| `email` | EmailConfig | Email alerting configuration |
| `webhook` | WebhookConfig | Webhook alerting configuration |

### ScalingConfig

| Field | Type | Description |
|-------|------|-------------|
| `minReplicas` | int32 | Minimum number of replicas |
| `maxReplicas` | int32 | Maximum number of replicas |
| `scaleUpThreshold` | int32 | Log frequency threshold to scale up |
| `scaleDownThreshold` | int32 | Log frequency threshold to scale down |
| `deploymentName` | string | Deployment to scale |

## Metrics

The controller exports the following Prometheus metrics:

- `logwatcher_matches_total`: Total number of log pattern matches
- `logwatcher_pods_restarted_total`: Total number of pods restarted
- `logwatcher_alerts_sent_total`: Total number of alerts sent (by type)
- `logwatcher_jobs_created_total`: Total number of jobs created
- `logwatcher_scaling_events_total`: Total number of scaling events (by direction)

### Accessing Metrics

```bash
# Port forward to access metrics
kubectl port-forward -n logwatcher-system deployment/logwatcher-controller-manager 8080:8080

# Access metrics endpoint
curl http://localhost:8080/metrics
```

## Development

### Building from Source

```bash
# Build the controller
make build

# Run tests
make test

# Generate manifests
make manifests
```

### Local Development

```bash
# Run locally
go run main.go

# Run with specific flags
go run main.go --metrics-bind-address=:8080 --health-probe-bind-address=:8081
```

## Troubleshooting

### Common Issues

1. **Controller not starting**: Check RBAC permissions and CRD installation
2. **No logs being monitored**: Verify pod label selectors and namespace
3. **Alerts not sending**: Check webhook URLs, SMTP settings, and network connectivity
4. **Scaling not working**: Verify deployment names and replica limits

### Debugging

```bash
# Check controller logs
kubectl logs -n logwatcher-system deployment/logwatcher-controller-manager

# Check LogWatcher status
kubectl get logwatchers -o yaml

# Check metrics
kubectl port-forward -n logwatcher-system deployment/logwatcher-controller-manager 8080:8080
curl http://localhost:8080/metrics
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 