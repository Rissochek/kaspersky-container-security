package model

const (
	Queued int = iota
	Running
	Done
	Failed
)

type Task struct {
	Id         string	`json:"id"`
	Payload    string	`json:"payload"`
	MaxRetries int		`json:"max_retries"`
	Status     int
}
