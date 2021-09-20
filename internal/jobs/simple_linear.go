package jobs

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	"github.com/Harzu/rebrain-webinar-redis/internal/services/redis/locker"
	"github.com/Harzu/rebrain-webinar-redis/internal/system/config"
)

const (
	simpleLinearJobName = "simple_linear_job"
	getLockInterval     = time.Second
)

type simpleLinearJob struct {
	spec    string
	lockTTL time.Duration
	locker  locker.Locker
	logger  *zerolog.Logger
}

func newSimpleLinearJob(
	cfg *config.Jobs,
	logger *zerolog.Logger,
	locker locker.Locker,
) job {
	return &simpleLinearJob{
		spec:    cfg.SimpleLinearJobSpec,
		lockTTL: cfg.SimpleLinearJobSpecLockTTl,
		logger:  logger,
		locker:  locker,
	}
}

func (j *simpleLinearJob) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobLogger := j.logger.With().Str("job_name", simpleLinearJobName).Logger()

	lock, err := j.locker.ObtainLinear(ctx, simpleLinearJobName, j.lockTTL, getLockInterval)
	if err != nil {
		jobLogger.Debug().Err(err).Msg("failed to obtain lock")
		return
	}

	defer func() {
		if err := lock.Release(ctx); err != nil {
			jobLogger.Debug().Err(err).Msg("failed to release lock")
		}
	}()

	jobLogger.Info().Msg("job start")
	jobLogger.Info().Msg("job work")
	time.Sleep(time.Second)
	jobLogger.Info().Msg("job success")
}

func (j *simpleLinearJob) Spec() string {
	return j.spec
}
