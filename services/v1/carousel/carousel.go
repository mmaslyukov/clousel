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
	crRepo := repository.Carousel.New(drv, &log)
	evRepo := repository.Event.New(drv, crRepo, &log)
	snRepo := repository.Snaphsot.New(drv, &log)
	man := manager.New(crRepo, snRepo, evRepo, &log)
	op := operator.New(evRepo, crRepo, snRepo, mqtt, cfg, &log)
	router := rest.New(man, op, &log)
	if err = mqtt.Connect(); err != nil {
		log.Err(err).Msg("Fail to connect")
	}

	if err = mqtt.Subscribe(topic.New(cfg.RootTopicSub()), cfg.DefaultQOS(), op); err != nil {
		log.Err(err).Msg("Fail to subscribe")
	}

	ticker := time.NewTicker(time.Second)
	go op.Tick(ticker)

	log.Info().Str("URL", cfg.ServerAddress()).Str("Key", cfg.ServerKeyPath()).Str("Cert", cfg.ServerCertPath()).Msg("Listening...")
	http.ListenAndServe(cfg.ServerAddress(), router)
	// err = http.ListenAndServeTLS(cfg.ServerAddress(), cfg.ServerCertPath(), cfg.ServerKeyPath(), router)
	log.Err(err).Send()
}
