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
