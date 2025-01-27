- **Dummy** objects are passed around but never actually used. Usually they are just used to fill parameter lists.
- **Fake** objects actually have working implementations, but usually take some shortcut which makes them not suitable for production (an in memory database is a good example).
- **Stubs** provide canned answers to calls made during the test, usually not responding at all to anything outside what's programmed in for the test. Stubs may also record information about calls, such as an email gateway stub that remembers the messages it 'sent', or maybe only how many messages it 'sent'.
- **Mocks** are what we are talking about here: objects pre-programmed with expectations which form a specification of the calls they are expected to receive



/
/cmd
/domain
 /carousel
  /internal
    ride.go
    manage.go
  /service
    service_manage.go
    service_ride.go
  port_ride.go
 /accountment
  /internal
    sales.go
    analytics.go
  /service
    service_sales.go
  /event
    event_rides_updated.go
  port_analytics.go
  port_sales.go

/framework
 /core
  i_event.go
  i_command.go
  i_message.go
  i_event_listener.go
  i_event_subscriber.go
  i_command_executer.go
  i_persistency.go
  i_logger.go
 /persistency
  /sqlite
    sqlite_persistency.go
  /postgre
    postgre_persistency.go
  persistency_factory.go
 /logger
 /utils

/repository
 /carousel
  stub_ride.go
  stub_manage.go
  adapter_ride.go
  adapter_manage.go
 /accountment
  stub_analytics.go
  stub_sales.go
  adapter_analytics.go
  adapter_sales.go
 
/infrastructure
 /carousel
  /router
   statistics.go
   manage_owner.go
   manage_carousel.go
  /gateway
   stub_refill.go
   adapter_refill.go
 /accountment
  /router
   initiate.go
   complete.go
   failure.go

/test
 /carousel
 /accountment

IMessage
 .
/_\
 |
 +--- ICommand
 |
 +--- IEvent

IMessage : {Name() -> String}
IEventListener : Notify(IEvent)
IEventSubscribable : Subscribe(IEvent, IEventListener)
ICommandExecuter : Execute(ICommand)

//core
service/application level
pyament.go
  func:
  //- Execute(IMessage) -> List[Event], error
  //- Subscribe(IMessage, IMessageListener)
payment/balance.go: 
  func:
  //- execute(IMessage) -> List[Event], error
  - update(carouselId, amount) -> List[Event], error 
  - refund(token, reason) -> List[Event], error
  - read_owned() -> List[History], error
  evt sub:
  - PaymentSucceded
  - RefundSucceded
  evt pub: 
  - BalanceUpdated
  

//payment/refund.go:
  //sub: RidesPaymentFailed, CarouselRefillFailed

carousel/refill.go: {port}
  func: 
  - execute(Event) -> List[Event], error
  - refill(carouselId, amount) -> error
  sub: BalanceUpdated 
  //pub: CarouselRefillSucceded, CarouselRefillFailed


//infrastrucure
stripe/hook.go {port}
rest/requests/refill.go {port}
rest/router/payment/initiate.go
rest/router/payment/complete.go
rest/router/payment/failure.go


use cases
1. refill:
Given carousel with id
When User scans QR code
Then Service returns structure with prices for each rounds set (1:1, 2:3, 3:5)

Given User chooses price entry
When User submits payment request with selected rounds set
Then Request activate page with card field

Given User filled in cards data
When User submits data
Then Service communicate with Stripe for payment

Given Payment successfull
When Service callback has been envoked by Stripe
Then Service creates an event publishes it and saves a copy to record table
Then Carousel controller part sents REST API call to tell Carousel sevice about paid rounds
Then updated record table on sucessfull response
