#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    print_error "kubectl is not installed. Please install kubectl first."
    exit 1
fi

# Check if we can connect to the cluster
if ! kubectl cluster-info &> /dev/null; then
    print_error "Cannot connect to Kubernetes cluster. Please check your kubeconfig."
    exit 1
fi

print_status "Starting deployment of LogWatcher Controller..."

# Create namespace if it doesn't exist
print_status "Creating namespace..."
kubectl create namespace logwatcher-system --dry-run=client -o yaml | kubectl apply -f -

# Build the controller
print_status "Building controller..."
make build

# Build Docker image
print_status "Building Docker image..."
make docker-build

# Install CRDs
print_status "Installing CRDs..."
make install

# Deploy the controller
print_status "Deploying controller..."
make deploy

# Wait for deployment to be ready
print_status "Waiting for deployment to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/logwatcher-controller-manager -n logwatcher-system

# Deploy examples
print_status "Deploying example LogWatcher resources..."
kubectl apply -f examples/

print_status "Deployment completed successfully!"

# Show status
print_status "Checking deployment status..."
kubectl get pods -n logwatcher-system
kubectl get logwatchers -A

print_status "To access metrics, run:"
echo "kubectl port-forward -n logwatcher-system deployment/logwatcher-controller-manager 8080:8080"
echo "curl http://localhost:8080/metrics"

print_status "To view logs, run:"
echo "kubectl logs -n logwatcher-system deployment/logwatcher-controller-manager -f"

print_status "To delete the deployment, run:"
echo "make undeploy"
echo "make uninstall" 