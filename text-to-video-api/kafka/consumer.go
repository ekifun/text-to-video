package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"text-to-video-api/models"
	"text-to-video-api/redis"

	"github.com/segmentio/kafka-go"
)

func main() {
	// ✅ Initialize Redis client before usage
	redis.InitRedis("redis:6379")

	brokerAddress := "kafka:9092"
	topic := "video-jobs"

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
		GroupID: "video-worker-group",
	})

	defer reader.Close()
	fmt.Println("📥 Kafka consumer started, waiting for jobs...")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("❌ Failed to read message: %v", err)
			continue
		}

		var job models.Job
		err = json.Unmarshal(msg.Value, &job)
		if err != nil {
			log.Printf("❌ Failed to unmarshal job: %v", err)
			continue
		}

		log.Printf("🎬 Processing job: %s, prompt: %s", job.ID, job.Prompt)
		processJob(&job)
	}
}

func processJob(job *models.Job) {
	// Simulate generation time
	job.Status = "processing"
	err := redis.SaveJob(*job)
	if err != nil {
		log.Printf("❌ Failed to update job to 'processing': %v", err)
		return
	}

	// Simulate video generation (replace with actual Genmo call)
	time.Sleep(5 * time.Second)

	// Simulate video output path
	videoPath := fmt.Sprintf("http://storage.example.com/videos/%s.mp4", job.ID)

	job.Status = "completed"
	job.VideoPath = videoPath

	err = redis.SaveJob(*job)
	if err != nil {
		log.Printf("❌ Failed to update job after completion: %v", err)
	} else {
		log.Printf("✅ Job %s completed! Video path: %s", job.ID, videoPath)
	}
}
