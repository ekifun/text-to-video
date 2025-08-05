#!/usr/bin/env bash
set -e  # Exit on any error

echo "🚀 Rebuilding Docker image for mochi-worker..."

# Set correct paths
DOCKERFILE_PATH="mochi-worker/Dockerfile"
BUILD_CONTEXT="mochi-worker"

# Build Docker image using correct Dockerfile and build context
sudo nerdctl build -f "$DOCKERFILE_PATH" -t mochi-worker:latest "$BUILD_CONTEXT"

echo "📦 Saving image to tar archive..."
sudo nerdctl save -o mochi-worker/mochi-worker_latest.tar mochi-worker:latest

echo "📥 Importing image into containerd (Kubernetes)..."
sudo ctr -n k8s.io images import mochi-worker/mochi-worker_latest.tar

echo "🔁 Restarting mochi-worker pod..."
kubectl delete pod -l app=mochi-worker --ignore-not-found

echo "✅ Deployment complete. Run: kubectl get pods -l app=mochi-worker"
