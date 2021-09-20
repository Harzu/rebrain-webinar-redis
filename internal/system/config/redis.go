package config

type Redis struct {
	URL        string `envconfig:"REDIS_URL"`
	MasterName string `envconfig:"REDIS_MASTER_NAME"`
	Password   string `envconfig:"REDIS_PASSWORD,optional"`
}
