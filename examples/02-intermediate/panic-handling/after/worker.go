package after

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type PanicHandler func(recovered interface{}, workerID int, job Job)

type WorkerPool struct {
	numWorkers   int
	jobs         chan Job
	panicHandler PanicHandler
	wg           sync.WaitGroup
}

type Job struct {
	ID   int
	Data string
}

type JobResult struct {
	Job       Job
	Err       error
	IsPanic   bool
	PanicInfo string
}

func NewWorkerPool(numWorkers int, panicHandler PanicHandler) *WorkerPool {
	if panicHandler == nil {
		panicHandler = defaultPanicHandler
	}

	return &WorkerPool{
		numWorkers:   numWorkers,
		jobs:         make(chan Job, 100),
		panicHandler: panicHandler,
	}
}

func defaultPanicHandler(recovered interface{}, workerID int, job Job) {
	log.Printf("PANIC RECOVERED in worker %d processing job %d: %v", workerID, job.ID, recovered)
}

func (wp *WorkerPool) Start(ctx context.Context, results chan<- JobResult) {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(ctx, i, results)
	}
}

func (wp *WorkerPool) worker(ctx context.Context, workerID int, results chan<- JobResult) {
	defer wp.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-wp.jobs:
			if !ok {
				return
			}
			result := wp.processJobWithRecovery(workerID, job)
			results <- result
		}
	}
}

func (wp *WorkerPool) processJobWithRecovery(workerID int, job Job) (result JobResult) {
	defer func() {
		if r := recover(); r != nil {
			wp.panicHandler(r, workerID, job)
			result = JobResult{
				Job:       job,
				IsPanic:   true,
				PanicInfo: fmt.Sprintf("%v", r),
				Err:       fmt.Errorf("panic recovered: %v", r),
			}
		}
	}()

	return wp.processJob(workerID, job)
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

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

func SafeGo(fn func(), panicHandler func(interface{})) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				if panicHandler != nil {
					panicHandler(r)
				} else {
					log.Printf("PANIC RECOVERED: %v", r)
				}
			}
		}()
		fn()
	}()
}
