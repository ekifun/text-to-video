package redis

import (
    "encoding/json"
    "text-to-video-api/models"
    "github.com/redis/go-redis/v9"
    "context"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

func SaveJob(job models.Job) error {
    ctx := context.Background()
    data, err := json.Marshal(job)
    if err != nil {
        fmt.Printf("❌ JSON marshal error: %v\n", err)
        return err
    }

    err = redisClient.Set(ctx, job.ID, data, 0).Err()
    if err != nil {
        fmt.Printf("❌ Redis SET error: %v\n", err)
    }
    return err
}

func GetJob(id string) (models.Job, error) {
    val, err := rdb.HGet(ctx, "job:"+id, "data").Result()
    if err != nil {
        return models.Job{}, err
    }

    var job models.Job
    json.Unmarshal([]byte(val), &job)
    return job, nil
}

func ListJobs() ([]models.Job, error) {
    keys, _ := rdb.Keys(ctx, "job:*").Result()
    var jobs []models.Job
    for _, key := range keys {
        val, _ := rdb.HGet(ctx, key, "data").Result()
        var job models.Job
        json.Unmarshal([]byte(val), &job)
        jobs = append(jobs, job)
    }
    return jobs, nil
}
