package workerpool

import (
	"sync"

	"github.com/Rissochek/kaspersky-container-security/model"
)

type WorkerPool struct {
	taskQueue		chan model.Task
	wg				*sync.WaitGroup
}

func NewWorkerPool(taskQueue chan model.Task, wg *sync.WaitGroup) *WorkerPool {
	return &WorkerPool{taskQueue: taskQueue, wg: wg}
}

type WorkerPoolInterface interface {
	HandleWorker(workerId int)
	HandleTask(task *model.Task)
}
