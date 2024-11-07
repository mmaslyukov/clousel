package broker

import (
	"carousel/core/operator"
	"carousel/infrastructure/broker/topic"
	"encoding/json"
	"fmt"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
)

const (
	tmMqttOpWait    = 1 * time.Second
	tmMqttConnWait  = 5 * time.Second
	tmMqttConnCheck = 60 * time.Second
)

type Broker struct {
	subscribers map[string]byte
	listeners   map[string]operator.IPortOperatorControllerMqtt
	sequenceTx  map[string]int
	sequenceRx  map[string]int
	url         string
	client      paho.Client
	log         *zerolog.Logger
}

func New(url string, log *zerolog.Logger) *Broker {
	b := &Broker{
		subscribers: make(map[string]byte),
		listeners:   make(map[string]operator.IPortOperatorControllerMqtt),
		sequenceTx:  make(map[string]int),
		sequenceRx:  make(map[string]int),
		url:         url,
		log:         log,
	}
	b.init()
	return b
}

func (b *Broker) init() {
	opts := paho.NewClientOptions()
	opts.AddBroker(b.url)
	opts.ClientID = "carousel-service-sg"
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

	// if e := b.Publish("/clousel", operator.CreateRequestPlay("123", "321"), 0); e != nil {
	// 	b.log.Err(e).Msg("Publish")
	// }
	return err
}

func (b *Broker) Publish(topic operator.Itopic, msg operator.IMessageGeneric, qos byte) error {
	var payload []byte
	var err error

	if !b.client.IsConnected() {
		err = fmt.Errorf("Broker is not connected")
		return err
	}

	err = fmt.Errorf("Unknown msg, Type=%s. Will not publish", msg.Name())
	b.sequenceTx[msg.Target()] += 1
	msg.SetSequenceId(b.sequenceTx[msg.Target()])
	switch m := msg.(type) {
	case *operator.RequestPlay:
		if payload, err = json.Marshal(m); err == nil {
			b.log.Debug().Str("Topic", topic.Get()).Int("QOS", int(qos)).Msgf("About to publish: %s", payload)
			token := b.client.Publish(topic.Get(), qos, false, payload)
			if !token.WaitTimeout(tmMqttOpWait) && token.Error() != nil {
				if token.Error() != nil {
					err = token.Error()
				} else {
					err = fmt.Errorf("WaitTimeout functions failed")
				}
			}
		} else {
			b.log.Err(err).Str("topic", topic.Get()).Str("Type", msg.Name()).Msg("Fail to unmarshall")
		}
	}
	return err
}

func (b *Broker) msgHandler(m paho.Message) {
	b.log.Debug().Str("topic", m.Topic()).Msg("Received")
	var msg operator.IMessageGeneric
	var err error
	for ok := true; ok; ok = false {

		if m.Duplicate() {
			b.log.Warn().Str("topic", m.Topic()).Msg("Got duplicated msg")
			break
		}

		var mg operator.MessageGeneric
		if err := json.Unmarshal(m.Payload(), &mg); err != nil {
			b.log.Err(err).Str("topic", m.Topic()).Msg("Fail to unmarshall")
			break
		}

		if mg.SeqId == b.sequenceRx[mg.CarId] {
			b.log.Warn().Str("topic", m.Topic()).Int("SeqNum", mg.SeqId).Msg("Message is handled before, skip")
			break
		}

		b.sequenceRx[mg.CarId] = mg.SeqId

		switch {
		case mg.MsgType == operator.MsgTypeEventHeartbeat:
			var ehb operator.EventHeartbeat
			if err = json.Unmarshal(m.Payload(), &ehb); err == nil {
				msg = &ehb
			}
		case mg.MsgType == operator.MsgTypeResponseAck:
			var ra operator.ResponseAck
			if err = json.Unmarshal(m.Payload(), &ra); err == nil {
				msg = &ra
			}
		case mg.MsgType == operator.MsgTypeEventCompleted:
			var ec operator.EventCompleted
			if err = json.Unmarshal(m.Payload(), &ec); err == nil {
				msg = &ec
			}
		}

		if err != nil {
			b.log.Err(err).Str("topic", m.Topic()).Str("Type", mg.MsgType).Msg("Fail to unmarshall")
			break
		}

		if msg == nil {
			b.log.Err(err).Str("topic", m.Topic()).Str("Type", mg.MsgType).Msg("Unknonw message")
			break
		}

		var listener operator.IPortOperatorControllerMqtt

		parent_topic := topic.New(m.Topic()).Parent()
		if listener = b.listeners[parent_topic]; listener == nil {
			b.log.Warn().Str("topic", parent_topic).Msg("Have no listener")
			break
		}
		listener.Notify(msg)
	}
}
func (b *Broker) subscribeInternal() error {
	var err error
	b.log.Debug().Msgf("Subscribe accroding to %v", b.subscribers)
	token := b.client.SubscribeMultiple(b.subscribers, func(c paho.Client, m paho.Message) {
		b.msgHandler(m)
	})
	if !token.WaitTimeout(tmMqttOpWait) && token.Error() != nil {
		err = token.Error()
		b.log.Err(err).Msgf("Fail to subscribe %v", b.subscribers)
	}
	return err
}

func (b *Broker) Subscribe(topic operator.Itopic, qos byte, listener operator.IPortOperatorControllerMqtt) error {
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
