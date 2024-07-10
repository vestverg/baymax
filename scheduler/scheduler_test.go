package scheduler

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewScheduledExecutorService(t *testing.T) {

	ctx := context.Background()

	scheduler := NewScheduledExecutorService(ctx)

	var wg sync.WaitGroup

	var completion int64
	start := time.Now().UnixNano()
	wg.Add(2)
	_ = scheduler.WithFixedRate(func(ctx context.Context) error {

		atomic.StoreInt64(&completion, time.Now().UnixNano())
		defer wg.Done()

		return nil
	}, 100*time.Millisecond, 0)
	wg.Wait()
	fmt.Println((completion - start) / int64(time.Millisecond))

}
func TestWithFixedDelay(t *testing.T) {
	ctx := context.Background()
	scheduler := NewScheduledExecutorService(ctx)

	var wg sync.WaitGroup
	var completion int64

	wg.Add(1)
	err := scheduler.WithFixedDelay(func(ctx context.Context) error {
		atomic.StoreInt64(&completion, time.Now().UnixNano())
		defer wg.Done()
		return nil
	}, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("failed to schedule task: %v", err)
	}

	wg.Wait()
	if completion == 0 {
		t.Fatalf("task did not run")
	}
}

func TestWithFixedRate(t *testing.T) {
	ctx := context.Background()
	scheduler := NewScheduledExecutorService(ctx)

	var wg sync.WaitGroup
	var completion int64

	wg.Add(1)
	err := scheduler.WithFixedRate(func(ctx context.Context) error {
		atomic.StoreInt64(&completion, time.Now().UnixNano())
		defer wg.Done()
		return nil
	}, 100*time.Millisecond, 0)
	if err != nil {
		t.Fatalf("failed to schedule task: %v", err)
	}

	wg.Wait()
	if completion == 0 {
		t.Fatalf("task did not run")
	}
}

func TestWithCronJob(t *testing.T) {
	ctx := context.Background()
	scheduler := NewScheduledExecutorService(ctx)

	var wg sync.WaitGroup
	var completion int64

	wg.Add(1)
	err := scheduler.WithCronJob(func(ctx context.Context) error {
		atomic.StoreInt64(&completion, time.Now().UnixNano())
		defer wg.Done()
		return nil
	}, "* * * * * *")
	if err != nil {
		t.Fatalf("failed to schedule task: %v", err)
	}

	wg.Wait()
	if completion == 0 {
		t.Fatalf("task did not run")
	}
}

func TestShutDown(t *testing.T) {
	ctx := context.Background()
	scheduler := NewScheduledExecutorService(ctx)

	scheduler.ShutDown()

	if scheduler.(*scheduledExecutorService).ctx.Err() != context.Canceled {
		t.Fatalf("scheduler did not shut down properly")
	}
}
