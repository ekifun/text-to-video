package kafka

import (
    "context"
    "encoding/json"
    "text-to-video-api/models"

    "github.com/segmentio/kafka-go"
	"log"
)

var writer = kafka.NewWriter(kafka.WriterConfig{
    Brokers: []string{"kafka.default.svc.cluster.local:9092"},
    Topic:    "video-jobs",
    Balancer: &kafka.LeastBytes{},
})

func ProduceJob(job models.Job) error {
    data, err := json.Marshal(job)
    if err != nil {
        return err
    }

    msg := kafka.Message{
        Key:   []byte(job.ID),
        Value: data,
    }

    err = writer.WriteMessages(context.Background(), msg)
    if err != nil {
        log.Printf("❌ Kafka write error: %v\n", err)
    }
    return err
}
