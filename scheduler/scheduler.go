package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/vestverg/baymax/collections/queue"
)

type Scheduler interface {
	WithFixedDelay(job Run, delay time.Duration) error
	WithFixedRate(job Run, rate time.Duration, initialDelay time.Duration) error
	WithCronJob(job Run, cron string) error
	ShutDown()
}

type FailedJob struct {
	job Job
	err error
}
type scheduledExecutorService struct {
	sync.RWMutex
	cancel context.CancelFunc
	ctx    context.Context
	queue  queue.BlockingQueue[Job]
	errors []FailedJob
}

func NewScheduledExecutorService(ctx context.Context) Scheduler {
	ctx, cancel := context.WithCancel(ctx)

	s := &scheduledExecutorService{
		cancel: cancel,
		queue:  queue.NewDelayedQueue[Job](),
	}
	s.start(ctx)
	return s
}

func (s *scheduledExecutorService) WithFixedDelay(run Run, delay time.Duration) error {
	s.RLock()
	defer s.RUnlock()
	job, err := NewFixedDelayJob(run, delay)
	if err != nil {
		return fmt.Errorf("fialed to schedule task %w", err)
	}
	s.queue.Offer(job)
	return nil
}

func (s *scheduledExecutorService) WithFixedRate(run Run, rate time.Duration, initialDelay time.Duration) error {
	s.Lock()
	defer s.Unlock()
	job, err := NewFixedRateJob(run, rate, initialDelay)
	if err != nil {
		return fmt.Errorf("fialed to schedule task %w", err)
	}
	s.queue.Offer(job)
	return nil
}

func (s *scheduledExecutorService) WithCronJob(run Run, cron string) error {
	s.Lock()
	defer s.Unlock()
	job, err := NewCronJob(run, cron)
	if err != nil {
		return fmt.Errorf("fialed to schedule task %w", err)
	}
	s.queue.Offer(job)
	return nil

}

func (s *scheduledExecutorService) ShutDown() {
	s.cancel()
	s.queue.Interrupt()
}

func (s *scheduledExecutorService) start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				job := s.pickJob()
				if job == nil {
					continue
				}
				go s.runJob(ctx, *job)
			}
		}

	}()
}

func (s *scheduledExecutorService) pickJob() *Job {
	job := s.queue.TakeWithTimeout(5 * time.Second)
	return job
}

func (s *scheduledExecutorService) runJob(ctx context.Context, job Job) {
	err := job.Run(ctx)
	defer func() {
		if r := recover(); r != nil {
			s.errors = append(s.errors, FailedJob{
				job: job,
				err: fmt.Errorf("job failed %s", r),
			})
		}
	}()
	if err != nil {
		s.errors = append(s.errors, FailedJob{
			job: job,
			err: fmt.Errorf("job failed %w", err),
		})
		return
	}
	s.queue.Offer(job)
}
