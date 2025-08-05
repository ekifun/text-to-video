#!/bin/bash

set -e

echo "🚀 Building the consumer image..."
sudo nerdctl build -t text-to-video-consumer:latest -f Dockerfile.consumer .

echo "📦 Saving image to tar archive..."
sudo nerdctl save -o text-to-video-consumer.tar text-to-video-consumer:latest

echo "📥 Importing image into containerd (Kubernetes)..."
sudo ctr -n k8s.io images import text-to-video-consumer.tar

echo "✅ Consumer image deployment complete. Now update the Deployment YAML with imagePullPolicy: Never and apply it:"
echo "👉 kubectl apply -f deployments/video-consumer.yaml"
