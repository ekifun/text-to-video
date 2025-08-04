package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"text-to-video-api/models"

	"github.com/redis/go-redis/v9"
)

var (
	ctx         = context.Background()
	redisClient *redis.Client
)

// InitRedis initializes the Redis client (call this once in main.go)
func InitRedis(addr string) {
	redisClient = redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})
}

// SaveJob stores a job in Redis under key: job:<id>, field: data
func SaveJob(job models.Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		fmt.Printf("❌ JSON marshal error: %v\n", err)
		return err
	}

	err = redisClient.HSet(ctx, "job:"+job.ID, "data", data).Err()
	if err != nil {
		fmt.Printf("❌ Redis HSET error: %v\n", err)
	}
	return err
}

// GetJob retrieves a job by ID
func GetJob(id string) (models.Job, error) {
	val, err := redisClient.HGet(ctx, "job:"+id, "data").Result()
	if err != nil {
		fmt.Printf("❌ Redis HGET error for job %s: %v\n", id, err)
		return models.Job{}, err
	}

	var job models.Job
	if err := json.Unmarshal([]byte(val), &job); err != nil {
		fmt.Printf("❌ JSON unmarshal error: %v\n", err)
		return models.Job{}, err
	}
	return job, nil
}

// ListJobs retrieves all stored jobs
func ListJobs() ([]models.Job, error) {
	keys, err := redisClient.Keys(ctx, "job:*").Result()
	if err != nil {
		fmt.Printf("❌ Redis KEYS error: %v\n", err)
		return nil, err
	}

	var jobs []models.Job
	for _, key := range keys {
		val, err := redisClient.HGet(ctx, key, "data").Result()
		if err != nil {
			fmt.Printf("⚠️ Redis HGET error for key %s: %v\n", key, err)
			continue
		}
		var job models.Job
		if err := json.Unmarshal([]byte(val), &job); err != nil {
			fmt.Printf("⚠️ JSON unmarshal error for key %s: %v\n", key, err)
			continue
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}
