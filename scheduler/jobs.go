package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/vestverg/baymax/collections/queue"

	"github.com/vestverg/baymax/cron"
)

type Run func(ctx context.Context) error

type Trigger interface {
	queue.Delayed
	GetNextExecution() time.Time
}

type Job interface {
	Trigger
	Run(ctx context.Context) error
}

type FixedRateJob struct {
	Trigger
	run           Run
	rate          time.Duration
	lastExecution *time.Time
}

func NewFixedRateJob(run Run, rate time.Duration, initialDelay time.Duration) (Job, error) {

	if run == nil {
		return nil, fmt.Errorf("invalid argument, run is nil")
	}

	if rate == 0 {
		return nil, fmt.Errorf("invalid argument, rate can't be 0")
	}

	initialExecution := time.Now().Add(-rate)
	initialExecution = initialExecution.Add(initialDelay)

	return &FixedRateJob{
		run:           run,
		rate:          rate,
		lastExecution: &initialExecution,
	}, nil
}

func (f *FixedRateJob) Run(ctx context.Context) error {
	now := time.Now()
	f.lastExecution = &now
	return f.run(ctx)
}

func (f *FixedRateJob) GetNextExecution() time.Time {
	now := time.Now()
	if f.lastExecution == nil {
		f.lastExecution = &now
		return now
	}
	return f.lastExecution.Add(f.rate)
}

func (f *FixedRateJob) GetDelay() int64 {
	if f.lastExecution == nil {
		return 0
	}
	return f.GetNextExecution().UnixNano() - time.Now().UnixNano()
}

type FixedDelayJob struct {
	Trigger
	run            Run
	delay          time.Duration
	lastCompletion *time.Time
}

func NewFixedDelayJob(run Run, delay time.Duration) (Job, error) {
	if run == nil {
		return nil, fmt.Errorf("invalid argument, run is nil")
	}
	if delay == 0 {
		return nil, fmt.Errorf("invalid argument, delay can't be 0")
	}
	initial := time.Now()
	return &FixedDelayJob{
		run:            run,
		delay:          delay,
		lastCompletion: &initial,
	}, nil
}

func (f *FixedDelayJob) Run(ctx context.Context) error {
	err := f.run(ctx)
	now := time.Now()
	f.lastCompletion = &now
	return err
}

func (f *FixedDelayJob) GetNextExecution() time.Time {
	from := f.lastCompletion
	now := time.Now()
	if from == nil {
		from = &now
	}
	from.Add(f.delay)
	return *from
}

func (f *FixedDelayJob) GetDelay() int64 {
	if f.lastCompletion == nil {
		return f.delay.Nanoseconds()
	}
	return f.GetNextExecution().UnixNano() - time.Now().UnixNano()
}

type CronJob struct {
	Job
	run            Run
	expression     *cron.CronExpression
	lastCompletion *time.Time
}

func NewCronJob(run Run, cronExpression string) (*CronJob, error) {
	expression, err := cron.Parse(cronExpression)
	if err != nil {
		return nil, fmt.Errorf("can't create CronJob: %w", err)
	}
	return &CronJob{
		run:        run,
		expression: expression,
	}, nil
}

func (cr *CronJob) Run(ctx context.Context) error {
	err := cr.run(ctx)
	now := time.Now()
	cr.lastCompletion = &now
	return err
}

func (cr *CronJob) GetNextExecution() time.Time {
	now := time.Now()
	if cr.lastCompletion == nil {
		cr.expression.Next(now)
	}
	return cr.expression.Next(*cr.lastCompletion)
}

func (cr *CronJob) GetDelay() int64 {
	return cr.GetNextExecution().UnixNano() - time.Now().UnixNano()
}
