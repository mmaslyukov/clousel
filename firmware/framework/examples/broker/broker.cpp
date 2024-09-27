#include <framework/broker.h>
#include <MQTTClient.h>
#include <stdlib.h>
#include <chrono>

#define ADDRESS "tcp://192.168.0.150:1883"
// #define ADDRESS "tcp://mqtt.eclipseprojects.io:1883"
#define CLIENTID "ExampleClientSub-app"
#define TOPIC "/test/one"
#define PAYLOAD "Hello World!"
#define QOS 1
#define TIMEOUT 10000L

class LoggerSystem : public core::logger::ILoggerSystem, public core::ITimestamp
{
public:
  virtual size_t get() const override { return 0; };
  virtual void output(const core::logger::Verbosity &verbosity, const char *tag, const char *data, size_t size) const override
  {
    printf("%d %s\n", verbosity.id(), data);
  };
};

class Timestamp : public core::ITimestamp
{
public:
  virtual size_t get() const override
  {
    return std::chrono::duration_cast<std::chrono::milliseconds>(std::chrono::system_clock::now().time_since_epoch()).count();
  }
};

// struct BrokerClientNone : public broker::IBrokerClient
// {
//   bool connect() const
//   {
//     return false;
//   }
//   bool disconnect() const
//   {
//     return false;
//   }
//   bool is_connected() const
//   {
//     return false;
//   }
//   bool publish(const char *topic, const broker::Message &msg)
//   {
//     return false;
//   }
//   bool subscribe(const char *topic, broker::IBrokerListener &listener)
//   {
//     return false;
//   }
// };

#define BROKER_SUBSCRIBERS 1
class Paho : public broker::Broker<BROKER_SUBSCRIBERS>
{
public:
  constexpr Paho(const char *address, const char *client_id, const core::logger::ILogger &logger)
      : broker::Broker<BROKER_SUBSCRIBERS>(logger),
        _client(nullptr),
        _conn_opts(MQTTClient_connectOptions_initializer)
  {
    int32_t rc = 0;
    if ((rc = MQTTClient_create(
             &_client, address, client_id, MQTTCLIENT_PERSISTENCE_NONE, NULL)) != MQTTCLIENT_SUCCESS)
    {
      _logger.err().log(TAG,
                        "Fail to create the client to address:%s, clentid:%s, return code: %d",
                        address, client_id, rc);
    }
    if ((rc = MQTTClient_setCallbacks(
             _client,
             static_cast<IBrokerConnectionListener *>(this),
             connlost, msgarrvd, delivered)) != MQTTCLIENT_SUCCESS)
    {
      _logger.err().log(TAG, "Fail to register callnacks, return code: %d", rc);
    }
  }

  virtual ~Paho()
  {
    int32_t rc;
    if (Paho::is_connected())
    {
      Paho::disconnect();
    }
    MQTTClient_destroy(&_client);
  }

  // ----- IBrokerClient
  bool connect()
  {
    _conn_opts.keepAliveInterval = 20;
    _conn_opts.cleansession = 1;
    int32_t rc;
    if ((rc = MQTTClient_connect(_client, &_conn_opts)) != MQTTCLIENT_SUCCESS)
    {
      _logger.err().log(TAG, "Failed to connect, return code %d", rc);
      return false;
    }
    return true;
  }

  bool disconnect()
  {
    int32_t rc;
    if ((rc = MQTTClient_disconnect(_client, _timeout)) != MQTTCLIENT_SUCCESS)
    {
      _logger.err().log(TAG, "Failed to disconnect, return code %d", rc);
      return false;
    }
    return true;
  }

  virtual bool is_connected() const
  {
    return MQTTClient_isConnected(_client);
  }

  virtual bool publish(const char *topic, const broker::Message &msg, const uint32_t qos = 0)
  {
    MQTTClient_message pubmsg = MQTTClient_message_initializer;
    pubmsg.payload = (void *)msg.data;
    pubmsg.payloadlen = msg.size;
    pubmsg.qos = qos;
    pubmsg.retained = 0;
    int32_t rc;
    broker::Token token;
    if ((rc = MQTTClient_publishMessage(_client, topic, &pubmsg, &token)) != MQTTCLIENT_SUCCESS)
    {
      _logger.err().log(TAG, "Failed to publish message, return code %d", rc);
      return false;
    }

    rc = MQTTClient_waitForCompletion(_client, token, _timeout);
    _logger.dbg().log(TAG, "Message to topic %s with delivery token %d delivered", topic, token);
    return false;
  }

  virtual bool subscribe(broker::IBrokerListener *listener, const char *topic, uint32_t qos = 0)
  {
    int32_t rc;
    if ((rc = MQTTClient_subscribe(_client, topic, qos)) != MQTTCLIENT_SUCCESS)
    {
      _logger.err().log(TAG, "Failed to subscribe, return code %d", rc);
      return false;
    }

    return Broker::add_subscriber(listener, topic);
  }

private:
  static void delivered(void *context, MQTTClient_deliveryToken dt)
  {
    if (context)
    {
      broker::IBrokerConnectionListener *listener =
          static_cast<broker::IBrokerConnectionListener *>(context);
      listener->delivered(static_cast<broker::Token>(dt));
    }
  }

  static int msgarrvd(void *context, char *topicName, int topicLen, MQTTClient_message *message)
  {
    if (context)
    {
      broker::IBrokerConnectionListener *listener =
          static_cast<broker::IBrokerConnectionListener *>(context);
      listener->arrived(topicName, broker::Message(message->payload, message->payloadlen));
    }
    MQTTClient_freeMessage(&message);
    MQTTClient_free(topicName);
    return 1;
  }

  static void connlost(void *context, char *cause)
  {
    if (context)
    {
      broker::IBrokerConnectionListener *listener =
          static_cast<broker::IBrokerConnectionListener *>(context);
      listener->disconnected(cause);
    }
  }

private:
  MQTTClient_connectOptions _conn_opts;
  MQTTClient _client;
  static constexpr const uint32_t _timeout = 10000;
};

class BrokerSubscriber : public broker::IBrokerListener
{
public:
  BrokerSubscriber(const core::logger::ILogger &logger)
      : _logger(logger) {}

private:
  virtual void notify(const char *topic, const broker::Message &msg)
  {
    _logger.inf().log("subscriber", "Received from topic: %s, message:", topic);
    _logger.raw().dump("subscriber", (const uint8_t *)msg.data, msg.size);
    _logger.raw().dump_ascii("subscriber", (const uint8_t *)msg.data, msg.size);
  }

private:
  const core::logger::ILogger &_logger;
};

void topic(const core::logger::ILogger &logger)
{
  const int cap = 20;
  broker::Topic<cap> t("/root/leaf");
  logger.raw().dump("topic", (const uint8_t *)t.get(), cap);
  logger.dbg().log("topic", "%d/%d|%s", t.len(), cap, t.get());
  t.append("flower");
  logger.raw().dump("topic", (const uint8_t *)t.get(), cap);
  logger.dbg().log("topic", "%d/%d|%s", t.len(), cap, t.get());
  logger.dbg().log("topic", "contain %d", t.contain("/root"));
  logger.dbg().log("topic", "contain %d", t.contain("/root/leaf"));
  logger.dbg().log("topic", "contain %d", t.contain("/root/leaf/flower"));
  logger.dbg().log("topic", "contain %d", t.contain("root/leaf/flower"));
  logger.dbg().log("topic", "contain %d", t.part_of("/root/leaf/flower/bee"));
}

int main()
{

  using namespace core::logger;
  Timestamp ts;
  LoggerSystem ls;
  char buff[128];
  const Logger logger(
      Printable(Configuration(buff, sizeof(buff), Verbosity("E")), ls, ts, true),
      Printable(Configuration(buff, sizeof(buff), Verbosity("W")), ls, ts, true),
      Printable(Configuration(buff, sizeof(buff), Verbosity("I")), ls, ts, true),
      Printable(Configuration(buff, sizeof(buff), Verbosity("D")), ls, ts, true),
      Dumpable(Configuration(buff, sizeof(buff), Verbosity("D")), ls, ts, true));

  topic(logger);
  return 0;

  if (true)
  {
    Paho broker(ADDRESS, CLIENTID, logger);
    broker.connect();
    // broker::IBrokerConnectionListener *pb = &broker;
    // pb->arrived("ssss", broker::Message());
    BrokerSubscriber bs(logger);
    broker.subscribe(&bs, "/test/one");

    int ch;
    do
    {
      ch = getchar();
    } while (ch != 'Q' && ch != 'q');
    broker.publish("/test/one", broker::Message("Hello1"));
    broker.publish("/test/one", broker::Message("Hello2"));
  }
}
