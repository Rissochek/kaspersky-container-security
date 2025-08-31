package workerpool

import (
	"sync"

	"github.com/Rissochek/kaspersky-container-security/model"
)

type WorkerPool struct {
	taskQueue		chan model.Task
	shutdownChan 	chan struct{}
	wg				*sync.WaitGroup
}

func NewWorkerPool(taskQueue chan model.Task, shutdownChan chan struct{}, wg *sync.WaitGroup) *WorkerPool {
	return &WorkerPool{taskQueue: taskQueue, shutdownChan: shutdownChan, wg: wg}
}

type WorkerPoolInterface interface {
	HandleWorker(workerId int)
	HandleTask(task *model.Task)
}
