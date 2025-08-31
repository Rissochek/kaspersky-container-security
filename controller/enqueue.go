package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Rissochek/kaspersky-container-security/model"
)

//контроллер для обработки входящих запросов по /enqueue 
func (controller *Controller) HandleEnqueue(w http.ResponseWriter, r *http.Request) {
	select {
	case <- controller.shutdownChan:
		log.Printf("task after shutting down is rejected")
		http.Error(w, "server is shutting down, new tasks is not accepting", http.StatusServiceUnavailable)
		return
	default:
	}
	
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	task := model.Task{}
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("failed to parse json with error: %v", err)
		http.Error(w, "wrong json format", http.StatusBadRequest)
		return
	}

	if err := ValidateData(&task); err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task.Status = model.Queued
	controller.taskQueue <- task
}

//функция для валидации полей id, payload. Поле max_retries валидируется само при биндинге в структуру (если оно пустое, то max_retires = 0)
func ValidateData(task *model.Task) error {
	if task.Id == "" {
		return errors.New("id field is empty")
	}

	if task.Payload == "" {
		return errors.New("payload field is empty")
	}

	return nil
}
