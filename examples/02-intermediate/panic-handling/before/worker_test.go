package before

import (
	"context"
	"testing"
	"time"
)

func TestWorkerPool_ProcessingNormalJobs(t *testing.T) {
	wp := NewWorkerPool(3)
	results := make(chan JobResult, 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wp.Start(ctx, results)

	jobs := []Job{
		{ID: 1, Data: "job1"},
		{ID: 2, Data: "job2"},
		{ID: 3, Data: "job3"},
	}

	for _, job := range jobs {
		wp.Submit(job)
	}

	wp.Close()

	successCount := 0
	timeout := time.After(2 * time.Second)

	for i := 0; i < len(jobs); i++ {
		select {
		case result := <-results:
			if result.Err == nil {
				successCount++
			}
		case <-timeout:
			t.Fatal("Test timed out waiting for results")
		}
	}

	if successCount != len(jobs) {
		t.Errorf("Expected %d successful jobs, got %d", len(jobs), successCount)
	}
}

func TestWorkerPool_ProcessingJobsWithErrors(t *testing.T) {
	wp := NewWorkerPool(3)
	results := make(chan JobResult, 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wp.Start(ctx, results)

	jobs := []Job{
		{ID: 1, Data: "job1"},
		{ID: 2, Data: "error"},
		{ID: 3, Data: "job3"},
	}

	for _, job := range jobs {
		wp.Submit(job)
	}

	wp.Close()

	successCount := 0
	errorCount := 0
	timeout := time.After(2 * time.Second)

	for i := 0; i < len(jobs); i++ {
		select {
		case result := <-results:
			if result.Err == nil {
				successCount++
			} else {
				errorCount++
			}
		case <-timeout:
			t.Fatal("Test timed out waiting for results")
		}
	}

	if successCount != 2 {
		t.Errorf("Expected 2 successful jobs, got %d", successCount)
	}

	if errorCount != 1 {
		t.Errorf("Expected 1 failed job, got %d", errorCount)
	}
}

func TestWorkerPool_PanicInGoroutine(t *testing.T) {
	wp := NewWorkerPool(3)
	results := make(chan JobResult, 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wp.Start(ctx, results)

	jobs := []Job{
		{ID: 1, Data: "job1"},
		{ID: 2, Data: "panic"},
		{ID: 3, Data: "job3"},
	}

	for _, job := range jobs {
		wp.Submit(job)
	}

	wp.Close()

	timeout := time.After(2 * time.Second)
	resultsReceived := 0

	for resultsReceived < len(jobs) {
		select {
		case <-results:
			resultsReceived++
		case <-timeout:
			t.Logf("Only received %d out of %d results before timeout", resultsReceived, len(jobs))
			return
		}
	}

	t.Errorf("Expected test to fail due to panic, but all results were received")
}
