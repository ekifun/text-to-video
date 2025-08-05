#!/bin/bash
set -e  # Exit immediately on error

echo "🚀 Rebuilding Docker image for mochi-worker..."

# Go to the project root directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$SCRIPT_DIR/.."
cd "$PROJECT_ROOT"

# Build Docker image using correct Dockerfile path
sudo nerdctl build -f mochi-worker/Dockerfile -t mochi-worker:latest mochi-worker/

echo "📦 Saving image to tar archive..."
sudo nerdctl save -o mochi-worker/mochi-worker_latest.tar mochi-worker:latest

echo "📥 Importing image into containerd (Kubernetes)..."
sudo ctr -n k8s.io images import mochi-worker/mochi-worker_latest.tar

echo "🔁 Restarting mochi-worker pod..."
kubectl delete pod -l app=mochi-worker --ignore-not-found

echo "✅ Deployment complete. Use 'kubectl get pods -l app=mochi-worker' to check status."
