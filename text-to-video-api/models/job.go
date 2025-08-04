package models

type Job struct {
    ID        string `json:"id"`
    Prompt    string `json:"prompt"`
    Status    string `json:"status"`     // pending, processing, completed, failed
    VideoPath string `json:"video_path"` // path to generated video
}
