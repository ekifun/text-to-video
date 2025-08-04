#!/bin/bash
set -e

echo "🧠 Starting Mochi Worker..."
export PYTHONPATH=/app:/app/src:${PYTHONPATH}
exec python3 mochi_worker.py
