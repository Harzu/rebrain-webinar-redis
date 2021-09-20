package config

import "time"

type Jobs struct {
	// cron spec https://en.wikipedia.org/wiki/Cron
	SimpleJobSpec        string        `envconfig:"SIMPLE_JOB_SPEC"`
	SimpleJobSpecLockTTl time.Duration `envconfig:"SIMPLE_JOB_LOCK_TTL,default=5m"`
	// cron spec https://en.wikipedia.org/wiki/Cron
	SimpleLinearJobSpec        string        `envconfig:"SIMPLE_LINEAR_JOB_SPEC"`
	SimpleLinearJobSpecLockTTl time.Duration `envconfig:"SIMPLE_LINEAR_JOB_LOCK_TTL,default=5m"`
}
