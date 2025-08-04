package main

import (
    "github.com/gin-gonic/gin"
    "text-to-video-api/handlers"
)

func main() {
    r := gin.Default()

    api := r.Group("/jobs")
    {
        api.POST("", handlers.SubmitJob)
        api.GET("/:id", handlers.GetJobStatus)
        api.GET("", handlers.ListJobs)
        api.GET("/:id/video", handlers.DownloadVideo)
    }

    r.Run(":8080")
}
