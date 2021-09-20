package jobs

import (
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"

	"github.com/Harzu/rebrain-webinar-redis/internal/services/redis/locker"
	"github.com/Harzu/rebrain-webinar-redis/internal/system/config"
)

type job interface {
	Run()
	Spec() string
}

type Runner interface {
	Start()
	Stop()
}

type jobsRunner struct {
	scheduler *cron.Cron
}

func New(cfg *config.Jobs, logger *zerolog.Logger, locker locker.Locker) (Runner, error) {
	scheduler := cron.New(cron.WithLogger(newCronLogger(logger)))

	simpleJob := newSimpleJob(cfg, logger, locker)
	if _, err := scheduler.AddJob(simpleJob.Spec(), simpleJob); err != nil {
		return nil, fmt.Errorf("failed to add simple job: %w", err)
	}

	simpleLinearJob := newSimpleLinearJob(cfg, logger, locker)
	if _, err := scheduler.AddJob(simpleLinearJob.Spec(), simpleLinearJob); err != nil {
		return nil, fmt.Errorf("failed to add simple linear job: %w", err)
	}

	return &jobsRunner{scheduler: scheduler}, nil
}

func (j jobsRunner) Start() {
	j.scheduler.Start()
}

func (j jobsRunner) Stop() {
	j.scheduler.Stop()
}
