package workerpool

import (
	"log"
	"math/rand"
	"time"

	"github.com/Rissochek/kaspersky-container-security/model"
)

func (pool *WorkerPool) HandleWorker(workerId int) {
	log.Printf("worker %v is starting", workerId)
	for task := range pool.taskQueue{
		pool.HandleTask(&task)
	}

	log.Printf("task queue closed, worker %v finished", workerId)
	pool.wg.Done()
}
//Здесь находится логика имитации работы с задачей
//jitter - пропорциональный к слип тайму, а также симметричный
//baseSleepTime растет экспоненциально благодаря битовому сдвигу (каждый раз увеличивается степень двойки)
func (pool *WorkerPool) HandleTask(task *model.Task) {
	task.Status = model.Running
	log.Printf("task %v running", task.Id)
	for fail := 0; fail <= task.MaxRetries; fail++ {
		executingTime := time.Duration((rand.Intn(401) + 100)) * time.Millisecond
		log.Printf("task %v executing for %v", task.Id, executingTime)
		time.Sleep(executingTime)

		failNum := rand.Intn(5)
		//failNum это что-то из [0, 1, 2, 3, 4] => равенство одному из чисел имеет шанс 1/5 * 100% = 20%
		if failNum == 0 {
			if fail != task.MaxRetries {
				baseSleepTime := time.Duration(1<<fail) * time.Second
				jitter := time.Duration(float64(baseSleepTime) * 0.15 * (rand.Float64()*2 - 1))
				sleepTime := baseSleepTime + jitter

				log.Printf("task %v sleeping for %v", task.Id, sleepTime)
				time.Sleep(sleepTime)
			}
			continue
		}

		task.Status = model.Done
		log.Printf("task %v done", task.Id)
		return
	}

	task.Status = model.Failed
	log.Printf("task %v failed", task.Id)
}
