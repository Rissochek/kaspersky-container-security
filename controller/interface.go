package controller

import (
	"net/http"

	"github.com/Rissochek/kaspersky-container-security/model"
)

type Controller struct {
	taskQueue chan model.Task
	shutdownChan chan struct{}
}

func NewController(taskQueue chan model.Task, shutdownChan chan struct{}) *Controller {
	return &Controller{taskQueue: taskQueue, shutdownChan: shutdownChan}
}

type ControllerInterface interface {
	HandleEnqueue(w http.ResponseWriter, r *http.Request)
	HealthCheck(w http.ResponseWriter, r *http.Request)
}