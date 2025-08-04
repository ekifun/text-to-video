import os
import json
import redis
import torch
from huggingface_hub import snapshot_download
from kafka import KafkaConsumer
from genmo.text2video import Text2VideoPipeline
from genmo.utils.download import load_weights
from datetime import datetime

# ENV: these can be moved to env vars or argparse later
KAFKA_TOPIC = "video-jobs"
KAFKA_BROKERS = ["localhost:9092"]
REDIS_HOST = "localhost"
MODEL_DIR = "./model_decoder"
OUTPUT_DIR = "./videos"

# Create output directory
os.makedirs(OUTPUT_DIR, exist_ok=True)

# Setup Redis client
r = redis.Redis(host=REDIS_HOST, port=6379, decode_responses=True)

# Load the model once at startup
print("Loading Mochi model...")
load_weights(MODEL_DIR)
device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
pipe = Text2VideoPipeline.from_pretrained(MODEL_DIR).to(device).eval()
print("Model ready ✅")

# Kafka consumer
consumer = KafkaConsumer(
    KAFKA_TOPIC,
    bootstrap_servers=KAFKA_BROKERS,
    value_deserializer=lambda m: json.loads(m.decode("utf-8")),
    auto_offset_reset="earliest",
    group_id="mochi-workers"
)

# Worker loop
for msg in consumer:
    job = msg.value
    job_id = job.get("id")
    prompt = job.get("prompt")
    print(f"[{datetime.now()}] 🎬 Processing job {job_id} | prompt: {prompt}")

    try:
        # Update Redis status to 'processing'
        r.hset(f"job:{job_id}", "data", json.dumps({
            "id": job_id,
            "prompt": prompt,
            "status": "processing"
        }))

        # Run model inference
        result = pipe(prompt, num_frames=24)
        video_path = os.path.join(OUTPUT_DIR, f"{job_id}.mp4")
        result.save(video_path)

        # Update Redis status to 'completed'
        r.hset(f"job:{job_id}", "data", json.dumps({
            "id": job_id,
            "prompt": prompt,
            "status": "completed",
            "video_path": video_path
        }))
        print(f"✅ Completed job {job_id} → {video_path}")

    except Exception as e:
        # Update Redis status to 'failed'
        r.hset(f"job:{job_id}", "data", json.dumps({
            "id": job_id,
            "prompt": prompt,
            "status": "failed",
            "error": str(e)
        }))
        print(f"❌ Failed job {job_id}: {str(e)}")
