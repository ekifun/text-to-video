#!/bin/bash

echo "✅ Checking for BuildKit..."

# Check if buildkitd is running
if ! pgrep -x "buildkitd" > /dev/null; then
    echo "🔧 BuildKit not running. Installing..."

    # Create directory for BuildKit binaries
    mkdir -p ~/buildkit-install && cd ~/buildkit-install

    # Download and extract BuildKit
    curl -LO https://github.com/moby/buildkit/releases/download/v0.12.5/buildkit-v0.12.5.linux-amd64.tar.gz
    tar -xvf buildkit-v0.12.5.linux-amd64.tar.gz

    # Install binaries
    sudo cp bin/* /usr/local/bin/

    # Start buildkitd in the background
    echo "🚀 Starting buildkitd..."
    sudo nohup buildkitd > /var/log/buildkitd.log 2>&1 &
    sleep 3
    echo "✅ buildkitd started"
else
    echo "✅ BuildKit is already running"
fi

echo "📦 Building consumer image..."
sudo nerdctl build -t text-to-video-consumer:latest -f Dockerfile.consumer .

echo "💾 Saving image to tar archive..."
sudo nerdctl save -o text-to-video-consumer.tar text-to-video-consumer:latest

echo "📥 Importing image into containerd (Kubernetes)..."
sudo ctr -n k8s.io images import text-to-video-consumer.tar

echo "✅ Done."
