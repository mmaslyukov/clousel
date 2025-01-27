package main

import (
	"gateway/core/dispatcher"
	"gateway/infra/broker"
	"gateway/infra/config"
	"gateway/infra/ipc"
	"gateway/infra/logger"
)

func main() {
	var err error
	log := logger.New()
	cfg := config.New()
	mqtt := broker.New(cfg, log)
	if err = mqtt.Connect(); err != nil {
		log.Err(err).Msg("Fail to connect")
	}
	ipc := ipc.IpcCreate(cfg, log)
	gw := dispatcher.DispatcherCreate(mqtt, ipc, cfg, log)

	gw.Run()

}
