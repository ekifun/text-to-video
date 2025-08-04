package main

import (
	"github.com/gin-gonic/gin"
	"text-to-video-api/handlers"
	"text-to-video-api/redis"
)

func main() {
	// Initialize Redis client (host should be "redis" in Kubernetes)
	redis.InitRedis("redis:6379")

	// Initialize Gin router
	r := gin.Default()

	// Optional: basic health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	// Job API group
	api := r.Group("/jobs")
	{
		api.POST("", handlers.SubmitJob)
		api.GET("/:id", handlers.GetJobStatus)
		api.GET("", handlers.ListJobs)
		api.GET("/:id/video", handlers.DownloadVideo)
	}

	// Start HTTP server on port 8080
	r.Run(":8080")
}
