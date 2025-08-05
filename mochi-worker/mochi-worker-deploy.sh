#!/bin/bash
set -e  # Exit immediately on error

echo "🚀 Rebuilding Docker image for mochi-worker..."
sudo nerdctl build -t mochi-worker:latest .

echo "📦 Saving image to tar archive..."
sudo nerdctl save -o mochi-worker_latest.tar mochi-worker:latest

echo "📥 Importing image into containerd (Kubernetes)..."
sudo ctr -n k8s.io images import mochi-worker_latest.tar

echo "🔁 Restarting mochi-worker pod..."
kubectl delete pod -l app=mochi-worker

echo "✅ Deployment complete. Use 'kubectl get pods -l app=mochi-worker' to check status."
