package before

import (
	"context"
	"fmt"
	"log"
	"time"
)

type WorkerPool struct {
	numWorkers int
	jobs       chan Job
}

type Job struct {
	ID   int
	Data string
}

type JobResult struct {
	Job Job
	Err error
}

func NewWorkerPool(numWorkers int) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		jobs:       make(chan Job, 100),
	}
}

func (wp *WorkerPool) Start(ctx context.Context, results chan<- JobResult) {
	for i := 0; i < wp.numWorkers; i++ {
		go func(workerID int) {
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-wp.jobs:
					if !ok {
						return
					}
					result := wp.processJob(workerID, job)
					results <- result
				}
			}
		}(i)
	}
}

func (wp *WorkerPool) processJob(workerID int, job Job) JobResult {
	log.Printf("Worker %d processing job %d\n", workerID, job.ID)
	time.Sleep(100 * time.Millisecond)

	if job.Data == "panic" {
		panic(fmt.Sprintf("Worker %d panicked on job %d", workerID, job.ID))
	}

	if job.Data == "error" {
		return JobResult{
			Job: job,
			Err: fmt.Errorf("worker %d failed to process job %d", workerID, job.ID),
		}
	}

	return JobResult{Job: job, Err: nil}
}

func (wp *WorkerPool) Submit(job Job) {
	wp.jobs <- job
}

func (wp *WorkerPool) Close() {
	close(wp.jobs)
}
