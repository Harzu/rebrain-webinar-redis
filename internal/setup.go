package internal

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/ziflex/lecho/v2"
	"go.uber.org/multierr"

	"github.com/Harzu/rebrain-webinar-redis/internal/delivery/httpdelivery"
	"github.com/Harzu/rebrain-webinar-redis/internal/jobs"
	"github.com/Harzu/rebrain-webinar-redis/internal/services/redis/locker"
	"github.com/Harzu/rebrain-webinar-redis/internal/system/config"
	"github.com/Harzu/rebrain-webinar-redis/internal/system/connections"
	"github.com/Harzu/rebrain-webinar-redis/internal/system/connections/redis"
	"github.com/Harzu/rebrain-webinar-redis/internal/system/logger"
)

type Application interface {
	Run() error
}

type application struct {
	config      *config.Config
	logger      *zerolog.Logger
	jobRunner   jobs.Runner
	httpRouter  *echo.Echo
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

	router := echo.New()
	router.Logger = lecho.From(*appLogger)
	handlersContainer := httpdelivery.NewHandlers(appLogger)
	handlersContainer.Register(router)

	return &application{
		config:      cfg,
		logger:      appLogger,
		jobRunner:   jobRunner,
		httpRouter:  router,
		redisCloser: redisClient,
	}, nil
}

func (a *application) Run() error {
	a.jobRunner.Start()

	go func() {
		if err := a.httpRouter.Start(fmt.Sprintf(":%d", a.config.Port)); err != nil {
			a.logger.Error().Err(err).Msg("failed to start http server")
		}
	}()

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
	if httpServerCloseErr := a.httpRouter.Close(); httpServerCloseErr != nil {
		err = multierr.Append(err, httpServerCloseErr)
	}

	return
}
