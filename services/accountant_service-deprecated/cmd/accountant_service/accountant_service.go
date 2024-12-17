package main

import (
	"accountant_service/domain/accountment"
	"accountant_service/domain/accountment/accountment_event"
	"accountant_service/domain/accountment/accountment_service"
	"accountant_service/domain/carousel/carousel_service"
	"accountant_service/framework/logger"
	"accountant_service/infrastructure/gateway"
	"accountant_service/repository/accountment_repository"
	"accountant_service/repository/carousel_repository"

	"github.com/google/uuid"
)

// type INamed interface {
// 	Name() string
// }
// type IEvent interface {
// 	INamed
// }

// type ICommand interface {
// 	INamed
// }

// type EventOne struct {
// }

// func (e *EventOne) Name() string {
// 	return "event.one"
// }

// type EventTwo struct {
// }

// func (e *EventTwo) Name() string {
// 	return "event.two"
// }

// type IEventListener interface {
// 	Notify(IEvent)
// }

// type IEventSubscribable interface {
// 	Subscribe(event IEvent, listener IEventListener)
// }

// type ICommandExecuter interface {
// 	Execute(ICommand)
// }

// type TestA struct {
// 	a       int
// 	clients map[string]IEventListener
// }

// func NewTestA() *TestA {
// 	return &TestA{a: 1, clients: make(map[string]IEventListener)}
// }

// type TestB struct {
// 	b int
// }

// func NewTestB() *TestB {
// 	return &TestB{}
// }

// func (t *TestA) demo() {
// 	one := EventOne{}
// 	for k, v := range t.clients {
// 		if k == one.Name() {
// 			v.Notify(&one)
// 		}
// 	}
// }

// func (t *TestA) Subscribe(event IEvent, listener IEventListener) {
// 	t.clients[event.Name()] = listener
// }

//	func (t *TestB) Notify(event IEvent) {
//		fmt.Printf("Notify by %s", event.Name())
//	}
// func main() {
// 	a := NewTestA()
// 	b := NewTestB()
// 	a.Subscribe(&EventOne{}, b)
// 	a.demo()
// 	// fmt.Printf("")
// }

func main() {
	var err error
	logger := logger.LoggerCreate()

	repo_sales := accountment_repository.StubSalesRepositoryCreate()
	repo_analytics := accountment_repository.StubAnalyticsRepositoryCreate()
	repo_carousel := carousel_repository.StubRideRepositoryCreate()

	gw := gateway.StubPublisherGatewatCreate()

	service_sales := accountment_service.ServiceSalesCreate(logger, repo_sales)
	service_analytics := accountment_service.ServiceAnalyticsCreate(logger, repo_analytics)
	service_ride := carousel_service.ServiceRideCreate(gw, repo_carousel, logger)

	service_sales.Subscribe(accountment_event.EventRidesUpdatedCreateEmpty(), service_ride)
	if err = service_sales.ApplyAndSaveReceipt(accountment.ReceiptDetails{Rides: 1}); err != nil {
		logger.Err().Println(err)
	}

	tags, _ := service_sales.ReadPriceTags(uuid.New())
	receipts, _ := service_analytics.LoadReceipts(uuid.New())
	logger.Dbg().Printf("rides: %v", tags)
	logger.Dbg().Printf("receipts: %v", receipts)

}
