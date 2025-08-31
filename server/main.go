package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Rissochek/kaspersky-container-security/controller"
	"github.com/Rissochek/kaspersky-container-security/model"
	"github.com/Rissochek/kaspersky-container-security/utils"
	"github.com/Rissochek/kaspersky-container-security/workerpool"
)

type Server struct {
	controller   controller.ControllerInterface
	workerPool   workerpool.WorkerPoolInterface
	shutdownChan chan struct{}
	workers      int
	wg           *sync.WaitGroup
}

func NewServer(controller controller.ControllerInterface, workerPool workerpool.WorkerPoolInterface, shutdownChan chan struct{}, workers int, wg *sync.WaitGroup) *Server {
	return &Server{controller: controller, workerPool: workerPool, shutdownChan: shutdownChan, workers: workers, wg: wg}
}

func InitDependencies() *Server {
	workers := utils.GetKeyFromEnv("WORKERS")
	queueSize := utils.GetKeyFromEnv("QUEUE_SIZE")

	taskQueue := make(chan model.Task, queueSize)
	shutdownChan := make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(workers)

	controller := controller.NewController(taskQueue, shutdownChan)
	pool := workerpool.NewWorkerPool(taskQueue, shutdownChan, wg)

	return NewServer(controller, pool, shutdownChan, workers, wg)
}

func main() {
	server := InitDependencies()
	http.HandleFunc("/enqueue", server.controller.HandleEnqueue)
	http.HandleFunc("/healthz", server.controller.HealthCheck)

	for workerId := range server.workers {
		go server.workerPool.HandleWorker(workerId)
	}
	go http.ListenAndServe(":8080", nil)

	server.HandleSignals()
}

func (server *Server) HandleSignals() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	signal := <-sigChan

	log.Printf("Received signal: %v", signal)
	close(server.shutdownChan)

	server.wg.Wait()
}
