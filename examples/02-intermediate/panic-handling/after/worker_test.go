package after

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestWorkerPool_ProcessingNormalJobs(t *testing.T) {
	wp := NewWorkerPool(3, nil)
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
	wp := NewWorkerPool(3, nil)
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

func TestWorkerPool_PanicRecoveryInGoroutine(t *testing.T) {
	var panicHandlerCalled bool
	var mu sync.Mutex

	customPanicHandler := func(recovered interface{}, workerID int, job Job) {
		mu.Lock()
		panicHandlerCalled = true
		mu.Unlock()
		t.Logf("Custom panic handler called: recovered=%v, workerID=%d, jobID=%d", recovered, workerID, job.ID)
	}

	wp := NewWorkerPool(3, customPanicHandler)
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
	panicCount := 0
	successCount := 0

	for resultsReceived < len(jobs) {
		select {
		case result := <-results:
			resultsReceived++
			if result.IsPanic {
				panicCount++
				if !strings.Contains(result.PanicInfo, "panicked") {
					t.Errorf("Expected panic info to contain 'panicked', got: %s", result.PanicInfo)
				}
			} else if result.Err == nil {
				successCount++
			}
		case <-timeout:
			t.Fatalf("Only received %d out of %d results before timeout", resultsReceived, len(jobs))
		}
	}

	if resultsReceived != len(jobs) {
		t.Errorf("Expected to receive %d results, got %d", len(jobs), resultsReceived)
	}

	if panicCount != 1 {
		t.Errorf("Expected 1 panic result, got %d", panicCount)
	}

	if successCount != 2 {
		t.Errorf("Expected 2 successful results, got %d", successCount)
	}

	mu.Lock()
	if !panicHandlerCalled {
		t.Error("Expected custom panic handler to be called")
	}
	mu.Unlock()
}

func TestWorkerPool_MultiplePanics(t *testing.T) {
	panicCount := 0
	var mu sync.Mutex

	customPanicHandler := func(recovered interface{}, workerID int, job Job) {
		mu.Lock()
		panicCount++
		mu.Unlock()
	}

	wp := NewWorkerPool(3, customPanicHandler)
	results := make(chan JobResult, 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wp.Start(ctx, results)

	jobs := []Job{
		{ID: 1, Data: "panic"},
		{ID: 2, Data: "panic"},
		{ID: 3, Data: "panic"},
		{ID: 4, Data: "job4"},
	}

	for _, job := range jobs {
		wp.Submit(job)
	}

	wp.Close()

	timeout := time.After(3 * time.Second)
	resultsReceived := 0

	for resultsReceived < len(jobs) {
		select {
		case <-results:
			resultsReceived++
		case <-timeout:
			t.Fatalf("Only received %d out of %d results before timeout", resultsReceived, len(jobs))
		}
	}

	mu.Lock()
	if panicCount != 3 {
		t.Errorf("Expected panic handler to be called 3 times, got %d", panicCount)
	}
	mu.Unlock()
}

func TestSafeGo_WithoutPanic(t *testing.T) {
	done := make(chan bool)

	SafeGo(func() {
		done <- true
	}, nil)

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("SafeGo did not complete")
	}
}

func TestSafeGo_WithPanic(t *testing.T) {
	var panicHandlerCalled bool
	var mu sync.Mutex

	customPanicHandler := func(recovered interface{}) {
		mu.Lock()
		panicHandlerCalled = true
		mu.Unlock()
	}

	done := make(chan bool)

	SafeGo(func() {
		defer func() {
			done <- true
		}()
		panic("test panic")
	}, customPanicHandler)

	select {
	case <-done:
		mu.Lock()
		if !panicHandlerCalled {
			t.Error("Expected panic handler to be called")
		}
		mu.Unlock()
	case <-time.After(1 * time.Second):
		t.Fatal("SafeGo did not complete")
	}
}

func TestWorkerPool_Wait(t *testing.T) {
	wp := NewWorkerPool(2, nil)
	results := make(chan JobResult, 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wp.Start(ctx, results)

	jobs := []Job{
		{ID: 1, Data: "job1"},
		{ID: 2, Data: "job2"},
	}

	for _, job := range jobs {
		wp.Submit(job)
	}

	wp.Close()

	done := make(chan bool)
	go func() {
		wp.Wait()
		done <- true
	}()

	for i := 0; i < len(jobs); i++ {
		<-results
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Wait() did not complete after all jobs were processed")
	}
}
