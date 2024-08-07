package main

import (
	// "context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"carousel_service/api"
	"carousel_service/internal/broker"
	"carousel_service/internal/carousel"
	"carousel_service/internal/config"
	"carousel_service/internal/logger"
	"carousel_service/internal/persistency"
	"carousel_service/internal/ports"

	_ "carousel_service/api"
	_ "carousel_service/internal/logger"
	_ "carousel_service/internal/persistency"
	_ "carousel_service/internal/utils"
)

func main1() {
	p := ports.NewPort[int](1)
	var j ports.PortInterface[int]
	j = p
	go func() {
		for {
			time.Sleep(time.Second)
			j.Send(1)
		}
	}()
	var i ports.PortInterface[int]
	i = p
	for {
		select {
		case x := <-i.Receiver():
			logger.Debug.Printf("rx: %d", x)
			break
		}
	}
}
func main2() {
	// context.Background().Done()
	type EventGeneral struct {
		CarouselId string `json:CarouselId`
		EventId    string `json:EventId`
	}
	type EventAck struct {
		CarouselId    string `json:CarouselId`
		CorrelationId string `json:CorrelationId`
	}
	type EventPlay struct {
		EventGeneral
		Command string `json:Command`
	}
	e := EventPlay{}
	result, _ := json.Marshal(e)
	fmt.Printf("json: %s", result)
}
func main3() {
	// layout := "2000-01-01T00:00:00Z"
	tf := "2024-08-02T09:31:26Z"
	t, err := time.Parse(time.RFC3339, tf)
	// t, err := time.Parse(layout, tf)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	logger.Info.Printf("Time %v", t)

}
func main() {

	db := persistency.NewPersistency()
	brokerRunner := broker.NewBrokerRunner()
	brokerApi := broker.NewBroker(&brokerRunner)
	carouselRunner := carousel.NewCarouselRunner(&brokerApi, &db)
	carouselApi := carousel.NewCarouselHandler(&carouselRunner, &db)

	brokerRunner.Connect()
	go carouselRunner.Run()
	go brokerRunner.Run()

	router := api.Router(&carouselApi)
	logger.Info.Printf("Starting")
	http.ListenAndServe(config.GetServerDest(), router)

}
