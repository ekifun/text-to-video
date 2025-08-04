#!/bin/bash

set -e  # Exit on error

echo "🚀 Rebuilding Docker image..."
sudo nerdctl build -t text-to-video-api:latest .

echo "📦 Saving image to tar archive..."
sudo nerdctl save -o text-to-video-api_latest.tar text-to-video-api:latest

echo "📥 Importing image into containerd (Kubernetes)..."
sudo ctr -n k8s.io images import text-to-video-api_latest.tar

echo "🔁 Restarting text-to-video-api pod..."
kubectl delete pod -l app=text-to-video-api

echo "✅ Deployment complete. Use 'kubectl get pods -l app=text-to-video-api' to check status."
