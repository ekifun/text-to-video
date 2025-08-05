import os
import json
import redis
import torch
from kafka import KafkaConsumer
from datetime import datetime

from genmo.mochi_preview.pipelines import (
    DecoderModelFactory,
    DitModelFactory,
    MochiSingleGPUPipeline,
    T5ModelFactory,
    linear_quadratic_schedule,
)

# ──────── 🔧 Config ──────────
KAFKA_TOPIC = "video-jobs"
KAFKA_BROKERS = ["kafka.default.svc.cluster.local:9092"]
REDIS_HOST = "redis"
MODEL_DIR = os.getenv("MODEL_DIR", "/models/mochi")  # Mount this in your pod
OUTPUT_DIR = "./videos"
os.makedirs(OUTPUT_DIR, exist_ok=True)

# ──────── 🧠 Load Model ────────
print("🧠 Loading Mochi model...")
pipeline = MochiSingleGPUPipeline(
    text_encoder_factory=T5ModelFactory(),
    dit_factory=DitModelFactory(
        model_path=f"{MODEL_DIR}/dit.safetensors", model_dtype="bf16"
    ),
    decoder_factory=DecoderModelFactory(
        model_path=f"{MODEL_DIR}/vae.safetensors",
    ),
    cpu_offload=True,
    decode_type="tiled_full",
)
print("✅ Model ready.")

# ──────── 🔌 Setup Redis ────────
r = redis.Redis(host=REDIS_HOST, port=6379, decode_responses=True)

# ──────── 🔁 Kafka Worker Loop ────────
consumer = KafkaConsumer(
    KAFKA_TOPIC,
    bootstrap_servers=KAFKA_BROKERS,
    value_deserializer=lambda m: json.loads(m.decode("utf-8")),
    auto_offset_reset="earliest",
    group_id="mochi-workers"
)

for msg in consumer:
    job = msg.value
    job_id = job.get("id")
    prompt = job.get("prompt")
    print(f"[{datetime.now()}] 🎬 Processing job {job_id} | prompt: {prompt}")

    try:
        # Update Redis: job status = processing
        r.hset(f"job:{job_id}", mapping={
            "id": job_id,
            "prompt": prompt,
            "status": "processing"
        })

        # Run inference
        video = pipeline(
            height=480,
            width=848,
            num_frames=31,
            num_inference_steps=64,
            sigma_schedule=linear_quadratic_schedule(64, 0.025),
            cfg_schedule=[4.5] * 64,
            batch_cfg=False,
            prompt=prompt,
            negative_prompt="",
            seed=12345,
        )

        video_path = os.path.join(OUTPUT_DIR, f"{job_id}.mp4")
        video.save(video_path)

        # Update Redis: job status = completed
        r.hset(f"job:{job_id}", mapping={
            "id": job_id,
            "prompt": prompt,
            "status": "completed",
            "video_path": video_path
        })

        print(f"✅ Completed job {job_id} → {video_path}")

    except Exception as e:
        r.hset(f"job:{job_id}", mapping={
            "id": job_id,
            "prompt": prompt,
            "status": "failed",
            "error": str(e)
        })
        print(f"❌ Failed job {job_id}: {str(e)}")
