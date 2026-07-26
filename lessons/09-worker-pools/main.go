package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID int
}

type Result struct {
	JobID  int
	Value  int
	Worker int
}

func worker(ctx context.Context, id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("worker %d: cancelled (%v)\n", id, ctx.Err())
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			time.Sleep(time.Duration(50+job.ID*5) * time.Millisecond)
			select {
			case results <- Result{JobID: job.ID, Value: job.ID * job.ID, Worker: id}:
			case <-ctx.Done():
				return
			}
		}
	}
}

func produce(ctx context.Context, jobs chan<- Job, n int) {
	defer close(jobs)
	for i := 1; i <= n; i++ {
		select {
		case <-ctx.Done():
			return
		case jobs <- Job{ID: i}:
		}
	}
}

func main() {
	const (
		numWorkers = 4
		numJobs    = 12
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jobs := make(chan Job)
	results := make(chan Result)

	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(ctx, w, jobs, results, &wg)
	}

	go produce(ctx, jobs, numJobs)

	// Close results only after every worker has returned, so the
	// range loop in main terminates cleanly.
	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		fmt.Printf("worker %d: job %2d -> %3d\n", r.Worker, r.JobID, r.Value)
	}

	fmt.Println("all done")
}
