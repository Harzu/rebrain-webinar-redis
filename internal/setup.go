package internal

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"go.uber.org/multierr"

	"github.com/Harzu/rebrain-webinar-redis/internal/config"
	"github.com/Harzu/rebrain-webinar-redis/internal/connections"
	"github.com/Harzu/rebrain-webinar-redis/internal/connections/redis"
	"github.com/Harzu/rebrain-webinar-redis/internal/jobs"
	"github.com/Harzu/rebrain-webinar-redis/internal/logger"
	"github.com/Harzu/rebrain-webinar-redis/internal/services/redis/locker"
)

type Application interface {
	Run() error
}

type application struct {
	config      *config.Config
	logger      *zerolog.Logger
	jobRunner   jobs.Runner
	redisCloser connections.Closer
}

func Setup(cfg *config.Config) (Application, error) {
	appLogger, err := logger.New(cfg.LogLevel)
	if err != nil {
		return nil, err
	}

	redisClient, err := redis.New(cfg.Redis)
	if err != nil {
		return nil, err
	}

	redisLocker := locker.New(redisClient)
	jobRunner, err := jobs.New(cfg.Jobs, appLogger, redisLocker)
	if err != nil {
		return nil, err
	}

	return &application{
		config:      cfg,
		logger:      appLogger,
		jobRunner:   jobRunner,
		redisCloser: redisClient,
	}, nil
}

func (a *application) Run() error {
	a.jobRunner.Start()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	return a.shutdown()
}

func (a *application) shutdown() (err error) {
	a.jobRunner.Stop()

	if redisCloseErr := a.redisCloser.Close(); redisCloseErr != nil {
		err = multierr.Append(err, redisCloseErr)
	}

	return
}
