#!/bin/bash
set -e

cd "$(dirname "$0")/.."  # Move to ~/text-to-video

echo "🚀 Rebuilding Docker image for mochi-worker..."
sudo nerdctl build -f mochi-worker/Dockerfile -t mochi-worker:latest .

echo "📦 Saving image to tar archive..."
sudo nerdctl save -o mochi-worker/mochi-worker_latest.tar mochi-worker:latest

echo "📥 Importing image into containerd (Kubernetes)..."
sudo ctr -n k8s.io images import mochi-worker/mochi-worker_latest.tar

echo "🔁 Restarting mochi-worker pod..."
kubectl delete pod -l app=mochi-worker

echo "✅ Deployment complete. Use 'kubectl get pods -l app=mochi-worker' to check status."
