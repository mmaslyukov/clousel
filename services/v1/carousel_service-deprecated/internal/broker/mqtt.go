package broker

import (
	"encoding/json"
	"fmt"
	"time"

	"carousel_service/internal/config"
	"carousel_service/internal/logger"
	. "carousel_service/internal/ports"

	pbr "carousel_service/internal/ports/port_broker"
	paho "github.com/eclipse/paho.mqtt.golang"
)

const (
	mqtt_con_wait   = 10 * time.Second
	mqtt_check_conn = 60 * time.Second
)

type MessageWrapper struct {
	topic   string
	payload []byte
}

type BrokerRunnerInterface interface {
	Connect() error
	Run()
}

func NewBroker(runner *BrokerRunner) BrokerAdapter {

	return BrokerAdapter{
		portCmd: runner.portBrokerCmd,
		portEvt: runner.portBrokerEvt,
	}
}
func NewBrokerRunner() BrokerRunner {
	return BrokerRunner{
		sequenceTx:    make(map[string]int),
		sequenceRx:    make(map[string]int),
		portBrokerCmd: NewPort[any](100),
		portBrokerEvt: NewPort[any](100),
		portMqttRx:    NewPort[MessageWrapper](100),
	}
}

type BrokerAdapter struct {
	portCmd PortInterface[any]
	portEvt PortInterface[any]
}

func (b *BrokerAdapter) SetSubscribers(subs pbr.Subscribers) {
	b.portCmd.Send(subs)
}
func (b *BrokerAdapter) PublishQueue(data any) {
	b.portCmd.Send(data)
}
func (b *BrokerAdapter) ListenQueue() PortInterface[any] {
	return b.portEvt
}

type BrokerRunner struct {
	sequenceTx map[string]int
	sequenceRx map[string]int
	url        string
	client     paho.Client
	//ports:
	portBrokerCmd PortInterface[any]
	portBrokerEvt PortInterface[any]
	portMqttRx    PortInterface[MessageWrapper]
}

func (b *BrokerRunner) init() *BrokerRunner {
	b.url = config.GetMQTTUrl()
	logger.Info.Printf("Mqtt Broker URL: '%s'", b.url)
	opts := paho.NewClientOptions()
	opts.AddBroker(b.url)
	opts.SetClientID("carousel-service")
	// opts.SetUsername("emqx")
	// opts.SetPassword("public")
	opts.SetDefaultPublishHandler(func(client paho.Client, msg paho.Message) {
		// logger.Info.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
		b.portMqttRx.Send(MessageWrapper{
			topic:   msg.Topic(),
			payload: msg.Payload(),
		})

	})
	opts.OnConnect = func(client paho.Client) {
		logger.Info.Printf("Connected to '%s'", config.GetMQTTUrl())
	}
	opts.OnConnectionLost = func(client paho.Client, err error) {
		logger.Info.Printf("Connect lost: %v\n", err)
	}

	b.client = paho.NewClient(opts)
	return b
}

func (b *BrokerRunner) keepConnected() error {
	const timeoutSeconds = time.Second * 5
	if !b.client.IsConnectionOpen() {
		token := b.client.Connect()
		ok := token.WaitTimeout(timeoutSeconds)
		if token.Error() != nil {
			logger.Error.Printf("Connection error: %s", token.Error())
			return token.Error()
		}
		if !ok {
			logger.Error.Printf("Unable to connect to '%s'", config.GetMQTTUrl())
			return fmt.Errorf("Connection timeout %d ms", timeoutSeconds)
		}
	}
	return nil
}

func (b *BrokerRunner) Connect() error {
	b.init()
	return b.keepConnected()
}

func (b *BrokerRunner) commandHandler(cmd any) {
	switch c := cmd.(type) {
	case nil:
		break
	case pbr.Subscribers:
		topics := make(map[string]byte)
		for k, v := range c.Carousels {
			topic := fmt.Sprintf("%s/%s", config.GetTopicCarousel(), k)
			topics[topic] = v
			logger.Info.Printf("Subscribed on '%s'", topic)
		}
		b.client.SubscribeMultiple(topics, nil)
		break
	case pbr.MessageCommand:
		b.sequenceTx[c.CarouselId]++
		c.SequenceNum = b.sequenceTx[c.CarouselId]
		json, err := json.Marshal(c)
		logger.Debug.Printf("Message is about to be sent: %s", json)
		if err == nil {
			topic := fmt.Sprintf("%s/%s", config.GetTopicCloud(), c.CarouselId)
			token := b.client.Publish(topic, 0, false, json)
			if !token.WaitTimeout(2000*time.Millisecond) && token.Error() != nil {
				logger.Error.Printf("Fail to publish the message")
			}
		} else {
			logger.Error.Printf("Fail to marshal to json")
		}
		break
	case pbr.MessageFeedDelme:
		json, err := json.Marshal(c)
		logger.Debug.Printf("Command Received: %s", json)
		if err == nil {
			token := b.client.Publish("/10169C25BF58/debug", 0, false, json)
			if !token.WaitTimeout(2000*time.Millisecond) && token.Error() != nil {
				logger.Error.Printf("Fail to publish the message")
			}
		} else {
			logger.Error.Printf("Fail to marshal to json")

		}
		break
	default:
		logger.Error.Printf("Type unknown: %+v", c) // here v has type interface{}
		break
	}
}

// func unmarshal[T any](msg MessageWrapper, f func(v T)) error {
// 	var v T
// 	var err error
// 	if err = json.Unmarshal(msg.payload, &v); err == nil {
// 		f(v)
// 	}
// 	return err
// }

func (b *BrokerRunner) mqttIncommingHandler(msg MessageWrapper) error {
	// Topic exmaples
	// - /clausel/carousel/UUID

	logger.Debug.Printf("Received topic:%s, paylod:%s", msg.topic, msg.payload)
	// err := errors.New("Default error for mqtt incoming handler")
	var err error
	var evm pbr.EventMinimal

	if err = json.Unmarshal(msg.payload, &evm); err != nil {
		return err
	}
	if b.sequenceRx[evm.CarouselId] == evm.SequenceNum {
		err = fmt.Errorf("Message is ignored due to duplication, seq:%d", evm.SequenceNum)
		return err

	}
	b.sequenceRx[evm.CarouselId] = evm.SequenceNum
	switch evm.Type {
	case pbr.NewEventAck().Type:
		var e pbr.EventAck
		if err = json.Unmarshal(msg.payload, &e); err != nil {
			return err
		}
		b.portBrokerEvt.Send(e)
		break
	case pbr.NewEventHeartbeat().Type:
		var e pbr.EventHeartbeat
		if err = json.Unmarshal(msg.payload, &e); err != nil {
			return err
		}
		b.portBrokerEvt.Send(e)
		break
	default:
		err = fmt.Errorf("Unknown type:%s", evm.Type)
	}

	// if err != nil {
	// 	err = unmarshal[pbr.EventAck](msg, func(v pbr.EventAck) {
	// 		if b.sequenceRx[v.CarouselId] != v.SequenceNum {
	// 			b.sequenceRx[v.CarouselId] = v.SequenceNum
	// 			logger.Debug.Printf("Send EventAck: %+v", v)
	// 			b.portBrokerEvt.Send(v)
	// 		} else {
	// 			logger.Warning.Printf("Message is ignored due to duplication")
	// 		}
	// 	})
	// }

	// if err != nil {
	// 	err = unmarshal[pbr.EventHeartbeat](msg, func(v pbr.EventHeartbeat) {
	// 		if b.sequenceRx[v.CarouselId] != v.SequenceNum {
	// 			b.sequenceRx[v.CarouselId] = v.SequenceNum
	// 			logger.Debug.Printf("Send Heartbeat: %+v", v)
	// 			b.portBrokerEvt.Send(v)
	// 		} else {
	// 			logger.Warning.Printf("Message is ignored due to duplication")
	// 		}
	// 	})
	// }
	return err
	//		logger.Warning.Printf("Unknown message topic:%s, paylod:%s", msg.topic, msg.payload)
}

func (b *BrokerRunner) Run() {
	ticker := time.NewTicker(mqtt_check_conn)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case cmd := <-b.portBrokerCmd.Receiver():
			b.commandHandler(cmd)
			break
			// logger.Debug.Print("Received: %+v", evt)
		case msg := <-b.portMqttRx.Receiver():
			if err := b.mqttIncommingHandler(msg); err != nil {
				logger.Error.Printf("Incomming message has been handled improperly: %s", err)
			}
			break
		case <-ticker.C:
			if err := b.keepConnected(); err != nil {
				logger.Error.Printf("Keep connection logic has been finished improperly: %s", err)

			}
			break
		}
	}
}

// func NewBrokerRunner() *BrokerRunner {
// 	return &BrokerRunner{
// 		portBrokerEvt: pbr.Prototype,
// 		portBrokerMsg: NewPort[any](100),
// 		portMqttRx:    NewPort[MessageWrapper](100),
// 	}
// }
