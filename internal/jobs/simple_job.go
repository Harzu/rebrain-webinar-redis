package jobs

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	"github.com/Harzu/rebrain-webinar-redis/internal/config"
	"github.com/Harzu/rebrain-webinar-redis/internal/services/redis/locker"
)

const simpleJobName = "simple_job"

type simpleJob struct {
	spec    string
	lockTTL time.Duration
	locker  locker.Locker
	logger  *zerolog.Logger
}

func newSimpleJob(
	cfg *config.Jobs,
	logger *zerolog.Logger,
	locker locker.Locker,
) job {
	return &simpleJob{
		spec:    cfg.SimpleJobSpec,
		lockTTL: cfg.SimpleJobSpecLockTTl,
		logger:  logger,
		locker:  locker,
	}
}

func (j *simpleJob) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobLogger := j.logger.With().Str("job_name", simpleJobName).Logger()

	if _, err := j.locker.Obtain(ctx, simpleJobName, j.lockTTL); err != nil {
		jobLogger.Debug().Err(err).Msg("failed to obtain log")
		return
	}

	jobLogger.Info().Msg("job start")
	jobLogger.Info().Msg("job work")
	time.Sleep(time.Second)
	jobLogger.Info().Msg("job success")
}

func (j *simpleJob) Spec() string {
	return j.spec
}
