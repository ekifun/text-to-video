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
    data, _ := json.Marshal(job)
    return rdb.HSet(ctx, "job:"+job.ID, "data", data).Err()
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
