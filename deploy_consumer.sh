#!/bin/bash

# Go to the script's directory
cd "$(dirname "$0")"

# Move to project root
cd ..

# Build the consumer image
sudo nerdctl build -t text-to-video-consumer:latest -f Dockerfile.consumer .

# Save and import into containerd for Kubernetes
sudo nerdctl save -o text-to-video-consumer.tar text-to-video-consumer:latest
sudo ctr -n k8s.io images import text-to-video-consumer.tar
