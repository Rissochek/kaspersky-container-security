package controller

import "net/http"

//обработчик для /healthz . Всегда возвращает 200
func (controller *Controller) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
}
