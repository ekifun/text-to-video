# Base image with CUDA 12.1, PyTorch, and Python 3.10
FROM pytorch/pytorch:2.2.0-cuda12.1-cudnn8-runtime

# Set workdir
WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    git \
    ffmpeg \
    libgl1 \
    libglib2.0-0 \
    && rm -rf /var/lib/apt/lists/*

# Install Python dependencies
COPY requirements.txt .
RUN pip install --upgrade pip && pip install -r requirements.txt

# Add source code
COPY mochi_worker.py .
COPY --chmod=755 ./entrypoint.sh .

# HuggingFace model and utils will be mounted externally or downloaded on the fly

# Entrypoint
ENTRYPOINT ["./entrypoint.sh"]
