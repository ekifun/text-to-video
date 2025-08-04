package handlers

import (
    "net/http"
    "text-to-video-api/kafka"
    "text-to-video-api/models"
    "text-to-video-api/redis"
    "text-to-video-api/utils"

    "github.com/gin-gonic/gin"
)

func SubmitJob(c *gin.Context) {
    fmt.Println("➡️  Entered SubmitJob handler")

    var req struct {
        Prompt string `json:"prompt"`
    }
    if err := c.BindJSON(&req); err != nil {
        fmt.Printf("❌ Failed to parse request: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    fmt.Printf("📨 Received prompt: %s\n", req.Prompt)

    jobID := utils.GenerateUUID()
    job := models.Job{
        ID:     jobID,
        Prompt: req.Prompt,
        Status: "pending",
    }

    fmt.Println("💾 Saving job to Redis")
    if err := redis.SaveJob(job); err != nil {
        fmt.Printf("❌ Failed to save job: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save job"})
        return
    }

    fmt.Println("📤 Publishing job to Kafka")
    if err := kafka.ProduceJob(job); err != nil {
        fmt.Printf("❌ Failed to produce job to Kafka: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue job"})
        return
    }

    fmt.Printf("✅ Job %s submitted successfully\n", jobID)
    c.JSON(http.StatusAccepted, gin.H{"job_id": jobID})
}

func GetJobStatus(c *gin.Context) {
    id := c.Param("id")
    job, err := redis.GetJob(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
        return
    }
    c.JSON(http.StatusOK, job)
}

func ListJobs(c *gin.Context) {
    jobs, err := redis.ListJobs()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list jobs"})
        return
    }
    c.JSON(http.StatusOK, jobs)
}

func DownloadVideo(c *gin.Context) {
    id := c.Param("id")
    job, err := redis.GetJob(id)
    if err != nil || job.Status != "completed" {
        c.JSON(http.StatusNotFound, gin.H{"error": "Video not available"})
        return
    }

    c.File(job.VideoPath)
}
