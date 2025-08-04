package kafka

import (
    "context"
    "encoding/json"
    "text-to-video-api/models"

    "github.com/segmentio/kafka-go"
)

var writer = kafka.NewWriter(kafka.WriterConfig{
    Brokers:  []string{"localhost:9092"},
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

    return writer.WriteMessages(context.Background(), msg)
}
