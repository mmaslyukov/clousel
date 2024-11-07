package main

import (
	"carousel/core/manager"
	"carousel/core/operator"
	"carousel/infrastructure/broker"
	"carousel/infrastructure/broker/topic"
	"carousel/infrastructure/config"
	"carousel/infrastructure/logger"
	"carousel/infrastructure/repository"
	"carousel/infrastructure/rest"
	"net/http"
	"time"
)

func main() {
	var err error
	log := logger.New()
	cfg := config.New()
	mqtt := broker.New(cfg.BrokerURL(), &log)
	drv := repository.DriverSQLite.New(cfg.DatabseURL())
	manRepo := repository.Carousel.New(drv, &log)
	evRepo := repository.Event.New(drv, &log)
	man := manager.New(manRepo, &log)
	op := operator.New(evRepo, mqtt, cfg, &log)
	router := rest.New(man, op, &log)
	if err = mqtt.Connect(); err != nil {
		log.Err(err).Msg("Fail to connect")
	}

	if err = mqtt.Subscribe(topic.New(cfg.RootTopicSub()), cfg.DefaultQOS(), op); err != nil {
		log.Err(err).Msg("Fail to subscribe")
	}

	ticker := time.NewTicker(time.Second)
	go op.Tick(ticker)

	log.Info().Str("URL", cfg.Server()).Msg("Listening...")
	http.ListenAndServe(cfg.Server(), router)
}
