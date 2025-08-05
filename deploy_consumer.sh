# deploy_consumer.sh
#!/bin/sh

set -e

echo "📦 Building consumer image..."
sudo nerdctl build --buildkit=false -t text-to-video-consumer:latest -f Dockerfile.consumer .

echo "💾 Saving image to tar archive..."
sudo nerdctl save -o text-to-video-consumer.tar text-to-video-consumer:latest

echo "📥 Importing image into containerd (Kubernetes)..."
sudo ctr -n k8s.io images import text-to-video-consumer.tar
