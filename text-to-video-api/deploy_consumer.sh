#!/bin/bash

# Move to project root (if not already there)
cd "$(dirname "$0")/.."

# Build the consumer image
sudo nerdctl build -t text-to-video-consumer:latest -f Dockerfile.consumer .

# Save and load into containerd
sudo nerdctl save -o text-to-video-consumer.tar text-to-video-consumer:latest
sudo ctr -n k8s.io images import text-to-video-consumer.tar
