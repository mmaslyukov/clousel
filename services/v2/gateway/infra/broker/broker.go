package broker

import (
	"encoding/json"
	"fmt"
	"gateway/core/dispatcher"
	"gateway/infra/broker/topic"
	"gateway/lib/fault"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	tmMqttOpWait    = 1 * time.Second
	tmMqttConnWait  = 5 * time.Second
	tmMqttConnCheck = 60 * time.Second
)

type Broker struct {
	subscribers map[string]byte
	listeners   map[string]dispatcher.IDispatcherMqttController
	sequenceTx  map[uuid.UUID]int
	sequenceRx  map[uuid.UUID]int
	url         string
	client      paho.Client
	log         *zerolog.Logger
}

func New(url string, log *zerolog.Logger) *Broker {
	b := &Broker{
		subscribers: make(map[string]byte),
		listeners:   make(map[string]dispatcher.IDispatcherMqttController),
		sequenceTx:  make(map[uuid.UUID]int),
		sequenceRx:  make(map[uuid.UUID]int),
		url:         url,
		log:         log,
	}
	b.init()
	return b
}

func (b *Broker) init() {
	opts := paho.NewClientOptions()
	opts.AddBroker(b.url)
	opts.ClientID = "gateway-service-sg"
	// opts.AutoReconnect = true
	opts.DefaultPublishHandler = func(c paho.Client, m paho.Message) {
		b.log.Warn().Str("topic", m.Topic()).Msg("Received message from the unexpected topic")
	}
	opts.OnConnect = func(client paho.Client) {
		b.log.Info().Str("URL", b.url).Msg("Connected")
		b.subscribeInternal()
	}
	opts.OnConnectionLost = func(client paho.Client, err error) {
		b.log.Info().Str("URL", b.url).Msg("Disconnected")
		b.Connect()
	}
	b.client = paho.NewClient(opts)
	// mqtt.ERROR = log.New(os.Stdout, "E", 0)
	// mqtt.CRITICAL = log.New(os.Stdout, "C", 0)
	// mqtt.WARN = log.New(os.Stdout, "W", 0)
	// mqtt.DEBUG = log.New(os.Stdout, "D", 0)
}

func (b *Broker) Connect() error {
	var err error
	b.log.Info().Str("URL", b.url).Msg("Connecting")
	if !b.client.IsConnectionOpen() {
		token := b.client.Connect()
		ok := token.WaitTimeout(tmMqttConnWait)
		err = token.Error()
		if !ok {
			err = fmt.Errorf("Connection timeout %d ms", tmMqttConnWait)
		}
	} else {
		b.log.Info().Msgf("Excessive connect request, aslready connected '%s'", b.url)
	}
	if err != nil {
		b.log.Err(err).Str("URL", b.url).Msg("Fail to connect")
	}

	// if e := b.Publish("/clousel", dispatcher.CreateRequestPlay("123", "321"), 0); e != nil {
	// 	b.log.Err(e).Msg("Publish")
	// }
	return err
}

func (b *Broker) Publish(topic dispatcher.ITopic, msg dispatcher.IMessageGeneric, qos byte) fault.IError {
	var err fault.IError

	if !b.client.IsConnected() {
		return fault.New(EBrokerNotConnected).Msg("Boker is not connected")
	}

	b.sequenceTx[msg.Target()] += 1
	msg.SetSequenceId(b.sequenceTx[msg.Target()])
	switch m := msg.(type) {
	case *dispatcher.RequestPlay:
		if payload, e := json.Marshal(m); e == nil {
			b.log.Debug().Str("Topic", topic.Get()).Int("QOS", int(qos)).Msgf("About to publish: %s", payload)
			token := b.client.Publish(topic.Get(), qos, false, payload)
			if !token.WaitTimeout(tmMqttOpWait) && token.Error() != nil {
				if token.Error() != nil {
					err = fault.New(EBrokerPaho).Msgf("%s", token.Error())
				} else {
					err = fault.New(EBrokerTimeout).Msg("WaitTimeout functions failed")
				}
			}
		} else {
			err = fault.New(EBrokerMarshall).Msg(e.Error())
			b.log.Err(e).Str("topic", topic.Get()).Str("Type", msg.Name()).Msg("Fail to marshall")
		}
	default:
		err = fault.New(EBrokerNotConnected).Msgf("Unknown msg, Type=%s. Will not publish", msg.Name())

	}
	return err
}

func (b *Broker) msgHandler(m paho.Message) {
	b.log.Debug().Str("topic", m.Topic()).Msg("Received")
	var msg dispatcher.IMessageGeneric
	var err error
	for ok := true; ok; ok = false {

		if m.Duplicate() {
			b.log.Warn().Str("topic", m.Topic()).Msg("Got duplicated msg")
			break
		}

		var mg dispatcher.MessageGeneric
		if err := json.Unmarshal(m.Payload(), &mg); err != nil {
			b.log.Err(err).Str("topic", m.Topic()).Str("Payload", string(m.Payload())).Msg("Fail to unmarshall MessageGeneric")
			break
		}

		if mg.SeqId == b.sequenceRx[mg.CarId] {
			b.log.Warn().Str("topic", m.Topic()).Int("SeqNum", mg.SeqId).Msg("Message is handled before, skip")
			break
		}

		b.sequenceRx[mg.CarId] = mg.SeqId

		switch {
		case mg.MsgType == dispatcher.MsgTypeEventHeartbeat:
			var ehb dispatcher.EventHeartbeat
			if err = json.Unmarshal(m.Payload(), &ehb); err == nil {
				msg = &ehb
			}
		case mg.MsgType == dispatcher.MsgTypeResponseAck:
			var ra dispatcher.ResponseAck
			if err = json.Unmarshal(m.Payload(), &ra); err == nil {
				msg = &ra
			}
		case mg.MsgType == dispatcher.MsgTypeEventCompleted:
			var ec dispatcher.EventCompleted
			if err = json.Unmarshal(m.Payload(), &ec); err == nil {
				msg = &ec
			}
		}

		if err != nil {
			b.log.Err(err).Str("topic", m.Topic()).Str("Type", mg.MsgType).Msgf("Fail to unmarshall %s", mg.MsgType)
			break
		}

		if msg == nil {
			b.log.Err(err).Str("topic", m.Topic()).Str("Type", mg.MsgType).Msg("Unknonw message")
			break
		}

		var listener dispatcher.IDispatcherMqttController

		parent_topic := topic.New(m.Topic()).Parent()
		if listener = b.listeners[parent_topic]; listener == nil {
			b.log.Warn().Str("topic", parent_topic).Msg("Have no listener")
			break
		}
		listener.BrokerNotify(msg)
	}
}
func (b *Broker) subscribeInternal() fault.IError {
	const fn = "Infra.Broker.subscribeInternal"
	var err fault.IError
	b.log.Debug().Msgf("%s: Subscribe according to %v", fn, b.subscribers)
	token := b.client.SubscribeMultiple(b.subscribers, func(c paho.Client, m paho.Message) {
		b.msgHandler(m)
	})
	if !token.WaitTimeout(tmMqttOpWait) && token.Error() != nil {
		err = fault.New(EBrokerPaho).Msg(token.Error().Error())
		b.log.Err(err).Msgf("Fail to subscribe %v", b.subscribers)
	}
	return err
}

func (b *Broker) Subscribe(topic dispatcher.ITopic, qos byte, listener dispatcher.IDispatcherMqttController) fault.IError {
	if false {
		var err error
		for t := range b.subscribers {
			token := b.client.Unsubscribe(t)
			if !token.WaitTimeout(tmMqttOpWait) && token.Error() != nil {
				err = token.Error()
				b.log.Err(err).Str("Topic", t).Int("QOS", int(qos)).Msg("Fail to unssubscribe")
			}
		}
	}

	b.log.Info().Str("Topic", topic.Subscribable()).Int("QOS", int(qos)).Msg("Adding subscriber")
	b.subscribers[topic.Subscribable()] = qos
	b.listeners[topic.Get()] = listener
	return b.subscribeInternal()
}
