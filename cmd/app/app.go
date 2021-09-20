package main

import (
	"log"

	"github.com/Harzu/rebrain-webinar-redis/internal"
	"github.com/Harzu/rebrain-webinar-redis/internal/system/config"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatalln(err)
	}

	worker, err := internal.Setup(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	if err := worker.Run(); err != nil {
		log.Fatalln(err)
	}
}
