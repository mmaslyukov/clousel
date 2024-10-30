/**
 * Build:
 * cmake --preset=default
 * cmake --build build
 */

#include <framework/broker.h>
#include <framework/core/io.h>
#include <framework/core/logger.h>
#include <framework/core/i_timestamp.h>

#include <infrastructure/config/config.h>
#include <infrastructure/status/status.h>

#include <service/mode/mode_service.h>
#include <service/coin/coin_service.h>
#include <service/web/web_service.h>

#include <platform/windows/router.h>


#include <MQTTClient.h>
#include <stdlib.h>
#include <chrono>
#include <thread>

// #define ADDRESS "tcp://192.168.0.150:1883"
// #define ADDRESS "tcp://mqtt.eclipseprojects.io:1883"
// #define CLIENTID "ExampleClientSub-app"
#define COIN_SCV_CAP 200
#define TOPIC "/test/one"
#define PAYLOAD "Hello World!"
#define QOS 1
#define TIMEOUT 10000L

class LoggerSystem : public core::logger::ILoggerSystem, public core::ITimestamp
{
public:
  virtual size_t get() const override { return 0; };
  virtual void output(const core::logger::Verbosity &verbosity, size_t tsms, const char *tag, const char *data, size_t size) const override
  {
    printf("%s (%zu) <%s> %s\n", verbosity.name(), tsms, tag, data);
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

#define BROKER_SUBSCRIBERS 1
class Paho
    : public broker::Broker<BROKER_SUBSCRIBERS>,
      // public service::mode::IPortAdapterBroker,
      public service::coin::IPortAdapterBroker
{
public:
  // constexpr Paho(const char *address, const char *client_id, const core::logger::ILogger &logger)
  constexpr Paho(const infra::Config &config, const core::logger::ILogger &logger)
      : broker::Broker<BROKER_SUBSCRIBERS>(logger),
        _client(nullptr),
        _config(config),
        _conn_opts(MQTTClient_connectOptions_initializer)
  {
    int32_t rc = 0;
    if ((rc = MQTTClient_create(
             &_client, config.broker_url()->data(), config.broker_client_id(), MQTTCLIENT_PERSISTENCE_NONE, NULL)) != MQTTCLIENT_SUCCESS)
    {
      _logger.err().log(TAG,
                        "Fail to create the client to address:%s, clentid:%s, return code: %d",
                        config.broker_url(), config.broker_client_id(), rc);
    }
    if ((rc = MQTTClient_setCallbacks(
             _client,
             static_cast<IBrokerConnectionListener *>(this),
             connlost, msgarrvd, delivered)) != MQTTCLIENT_SUCCESS)
    {
      _logger.err().log(TAG, "Fail to register callbacks, return code: %d", rc);
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
  virtual bool is_ready() const override
  {
    return true;
  }
  // ----- IBrokerClient
  bool connect() override
  {
    _conn_opts.keepAliveInterval = 20;
    _conn_opts.cleansession = 1;
    int32_t rc;

    if (is_connected())
    {
      return true;
    }

    _logger.inf().log(TAG, "Connecting to %s", _config.broker_url());
    if ((rc = MQTTClient_connect(_client, &_conn_opts)) != MQTTCLIENT_SUCCESS)
    {
      _logger.err().log(TAG, "Failed to connect to %s, return code %d", _config.broker_url(), rc);
      return false;
    }
    connected();

    for (size_t i = 0; i < _index_sub; i++)
    {
      if ((rc = MQTTClient_subscribe(_client, _subs[i].topic.get(), _subs[i].qos)) != MQTTCLIENT_SUCCESS)
      {
        _logger.err().log(TAG, "Failed to subscribe, return code %d", rc);
      }
    }

    return true;
  }

  bool disconnect() override
  {
    _logger.dbg().log(TAG, "Disconnect called");
    disconnected("By call of disconnect() function");

    int32_t rc;
    for (size_t i = 0; i < _index_sub; i++)
    {
      if ((rc = MQTTClient_unsubscribe(_client, _subs[i].topic.get())) != MQTTCLIENT_SUCCESS)
      {
        _logger.err().log(TAG, "Failed to unsubscribe, return code %d", rc);
      }
    }

    if ((rc = MQTTClient_disconnect(_client, _timeout)) != MQTTCLIENT_SUCCESS)
    {
      _logger.err().log(TAG, "Failed to disconnect, return code %d", rc);
      return false;
    }
    return true;
  }

  virtual bool is_connected() const override
  {
    return MQTTClient_isConnected(_client);
  }

  virtual broker::Token publish(const broker::ITopic &topic, const broker::Message &msg, const uint32_t qos = 0) override
  {
    MQTTClient_message pubmsg = MQTTClient_message_initializer;
    pubmsg.payload = (void *)msg.data;
    pubmsg.payloadlen = msg.size;
    pubmsg.qos = qos;
    pubmsg.retained = 0;
    int32_t rc;
    broker::Token token;
    if ((rc = MQTTClient_publishMessage(_client, topic.get(), &pubmsg, token.id_ptr())) != MQTTCLIENT_SUCCESS)
    {
      _logger.err().log(TAG, "Failed to publish message, return code %d", rc);
      return token.set_result(false);
    }

    rc = MQTTClient_waitForCompletion(_client, token.id(), _timeout);
    _logger.inf().log(TAG, "Message to topic %s with delivery token %d delivered", topic.get(), token);
    return token.set_result(true);
  }

  // virtual bool subscribe(broker::IBrokerListener *listener, const char *topic, uint32_t qos = 0)
  // {
  //   int32_t rc;
  //   if ((rc = MQTTClient_subscribe(_client, topic, qos)) != MQTTCLIENT_SUCCESS)
  //   {
  //     _logger.err().log(TAG, "Failed to subscribe, return code %d", rc);
  //     return false;
  //   }

  //   return Broker::subscribe(listener, topic);
  // }

private:
  static void delivered(void *context, MQTTClient_deliveryToken dt)
  {
    if (context)
    {
      broker::IBrokerConnectionListener *listener =
          static_cast<broker::IBrokerConnectionListener *>(context);
      listener->delivered(broker::Token(true, dt));
    }
  }

  static int msgarrvd(void *context, char *topicName, int topicLen, MQTTClient_message *message)
  {
    if (context)
    {
      broker::IBrokerConnectionListener *listener =
          static_cast<broker::IBrokerConnectionListener *>(context);
      listener->arrived(broker::TopicRef(topicName), broker::Message(message->payload, message->payloadlen));
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
  const infra::Config &_config;
  MQTTClient _client;
  MQTTClient_connectOptions _conn_opts;
  static constexpr const uint32_t _timeout = 10000;
  static constexpr const char *TAG = "broker.paho";
};

// class BrokerSubscriber : public broker::IBrokerListener
// {
// public:
//   BrokerSubscriber(const core::logger::ILogger &logger)
//       : _logger(logger) {}

// private:
//   virtual void notify(const char *topic, const broker::Message &msg)
//   {
//     _logger.inf().log("subscriber", "Received from topic: %s, message:", topic);
//     _logger.raw().dump("subscriber", (const uint8_t *)msg.data, msg.size);
//     _logger.raw().dump_ascii("subscriber", (const uint8_t *)msg.data, msg.size);
//   }

// private:
//   const core::logger::ILogger &_logger;
// };

class Wifi : public service::mode::IPortAdapterWifi
{
public:
  Wifi(const core::logger::ILogger &logger) : _logger(logger), _station(false), _softap(false) {}

  virtual bool swith_to_softap() override
  {
    _softap = true;
    _station = false;
    _logger.inf().log(TAG, "Mode has been switched to SoftAp Mode");
    return true;
  }
  virtual bool swith_to_station() override
  {

    _softap = false;
    _station = true;
    _logger.inf().log(TAG, "Mode has been switched to Station Mode");
    return true;
  }
  virtual bool is_softap() const override
  {

    return _softap;
  }
  virtual bool is_station() const override
  {
    return _station;
  }
  virtual bool is_station_connected() const override
  {
    return _station;
  }

private:
  const core::logger::ILogger &_logger;

  bool _station;
  bool _softap;
  static constexpr const char *TAG = "wifi";
};

class Button : public core::IRunnable
{
public:
  Button(service::mode::IPortButtonController &bc) : _bc(bc)
  {
  }

  virtual void run() override
  {
    char ch = getchar();
    switch (ch)
    {
    case 'p':
      _pressed = true;
      _clicked = false;
      _bc.pressed();
      break;
    case 'c':
      _pressed = false;
      _clicked = true;
      _bc.clicked();
      break;
    case 'r':
      _pressed = false;
      _clicked = false;
      break;
    default:
      break;
    }
  }

private:
  bool _pressed;
  bool _clicked;
  service::mode::IPortButtonController &_bc;
};

struct Flash : public persistency::IPersistencyFlash
{
  Flash(const core::logger::ILogger &logger)
      : _logger(logger) {}
  virtual bool load(uint8_t *memory, size_t size) const
  {
    _logger.dbg().log(TAG, "Loaded");
    return true;
  };
  virtual bool save(const uint8_t *memory, size_t size) const
  {
    _logger.dbg().log(TAG, "Saved");
    return true;
  };

private:
  const core::logger::ILogger &_logger;
  static constexpr const char *TAG = "flash";
};
class MyGpio : public core::io::IActuator<bool>
{
public:
  MyGpio(const char *name, uint32_t id, const core::logger::ILogger &logger)
      : _name(name), _value(false), _id(id), _logger(logger)
  {
    log();
  }
  virtual uint32_t id() const override
  {
    return _id;
  }
  virtual void set(const bool &value) override
  {
    if (_value != value)
    {
      _value = value;
      log();
    }
  }
  virtual bool get() const override
  {
    return _value;
  }

private:
  void log()
  {
    _logger.dbg().log("gpio", "%c GPIO%d[%s]", _value ? 'I' : 'O', _id, _name);
  }

private:
  const char *_name;
  bool _value;
  uint32_t _id;
  const core::logger::ILogger &_logger;
};

void button_run(Button *btn)
{
  while (1)
  {
    btn->run();
  }
}

struct ScenarioMaker : public strategy::IScenarioMaker
{
  ScenarioMaker(const Timestamp &ts, const core::logger::ILogger &logger, infra::Config &config): _ts(ts), _logger(logger), _config(config) {}
  virtual strategy::StepRunner &make() const override
  {
    // static MyGpio coin_pin("coin-ctrl", 1, logger);
    const infra::CoinPulseProps *props = _config.coin_pulse_props();
    if (!props->is_valid())
    {
      _config.persistency().reset_default(infra::Config::PersistencyId::COIN_PULSE_PROPS);
    }

    _logger.inf().log("TAG", "count:%d, duration:%d", props->count, props->duration);

    static MyGpio coin_pin("coin-out", 1, _logger);

    static strategy::StepWait swt(props->duration, _ts);
    static strategy::StepActuatorOff saoff(coin_pin);
    static strategy::StepActuatorOn saon(coin_pin);
    static strategy::IStep *st[17];
    memset(st, 0, sizeof(st));

    for (uint8_t count = 0, index = 0; count < props->count; count++)
    {
      if (index < sizeof(st) / sizeof(st[0]))
      {
        st[index++] = &saon;
      }
      if (index < sizeof(st) / sizeof(st[0]))
      {
        st[index++] = &swt;
      }
      if (index < sizeof(st) / sizeof(st[0]))
      {
        st[index++] = &saoff;
      }
      if (index < sizeof(st) / sizeof(st[0]))
      {
        st[index++] = &swt;
      }
    }

    static strategy::StepRunner step_runner(st, sizeof(st) / sizeof(st[0]));
    return step_runner;
  }

private:
  const Timestamp &_ts;
  const core::logger::ILogger &_logger;
  infra::Config &_config;
};



void _run(service::mode::ModeService *modsvc, infra::Status *status, service::coin::CoinService<COIN_SCV_CAP> *coinsvc)
{
  while (1)
  {
    modsvc->run();
    status->run();
    coinsvc->run();
    std::this_thread::sleep_for(std::chrono::milliseconds(100));
  }
}

#include <service/coin/messages/command.h>
int platground_json_parsers()
{
  service::coin::msg::ResponseConfig rc("321", "123", "", service::coin::msg::Config(), 1);
  
  char buf[256];
  rc.dump(buf, sizeof(buf));
  printf("json:%s\n", buf);
  printf("-----------------------------\n");
  
  auto cmd_read = R"({"Type":"MessageCommand","CarouselId":"550e8400-e29b-41d4-a716-446655440000","SequenceNum":39,"EventId":"b2ea6e51-6ffa-444f-b0e7-f3103bf5a244","Command":"ConfigRead"})";
  // auto cmd_write = R"({"Type":"MessageCommand","CarouselId":"650e8400-e29b-41d4-a716-446655440000","SequenceNum":40,"EventId":"b2ea6e51-6ffa-444f-b0e7-f3103bf5a244","Command":"ConfigWrite","Config":{"CoinPulseCnt":3,"CoinPulseDur":300}})";
  auto cmd_write = R"({"Type":"MessageCommand","CarouselId":"550e8400-e29b-41d4-a716-446655440000","SequenceNum":1,"EventId":"cedb3510-c87f-4f7d-a190-2f1f8412ff29","Command":"ConfigWrite", "Config":{"BrokerUrl":"fff"}})";
  service::coin::msg::CommandComposite cc;
  bool pr = cc.parse(cmd_read, strlen(cmd_read));

  printf("pr:%d, %s\n", pr, cc.general.type.value.data());

  cc.general.dump(buf, sizeof(buf));
  printf("general:%s\n", buf);
  cc.config.value.dump(buf, sizeof(buf));
  printf("config:%s\n", buf);
  printf("-----------------------------\n");
  
  cc.clear();
  pr = cc.parse(cmd_write, strlen(cmd_write));
  printf("pr:%d\n", pr);
  cc.general.dump(buf, sizeof(buf));
  printf("general:%s\n", buf);
  cc.config.value.dump(buf, sizeof(buf));
  printf("config:%s\n", buf);
  printf("-----------------------------\n");
  
  return 0;

}
int main()
{
  using namespace core::logger;

  Timestamp ts;
  LoggerSystem ls;
  char buff[256];
  const Logger logger(
      Printable(Configuration(buff, sizeof(buff), Verbosity("E")), ls, ts, true),
      Printable(Configuration(buff, sizeof(buff), Verbosity("W")), ls, ts, true),
      Printable(Configuration(buff, sizeof(buff), Verbosity("I")), ls, ts, true),
      Printable(Configuration(buff, sizeof(buff), Verbosity("D")), ls, ts, true),
      Printable(Configuration(buff, sizeof(buff), Verbosity("V")), ls, ts, true),
      Dumpable(Configuration(buff, sizeof(buff), Verbosity("D")), ls, ts, true));

  // const char* json_str =
  // R"({"CarouselId":"550e8400-e29b-41d4-a716-446655440000","EventId":"35215bd3-491a-49e6-8838-1dcccca58d39","Type":"MessageCommand","Command":"Play","SequenceNum":5})";
  // service::coin::msg::Command cmd;
  // char json[200];
  // if (!cmd.parse(json_str))
  // {
  //   logger.err().log("WIN", "Fail to parse");
  // }

  // if (!cmd.dump(json, sizeof(json)))
  // {
  //   logger.err().log("WIN", "Fail to dump");
  // }
  // else
  // {
  //   logger.inf().log("WIN", "Json: %s", json);
  // }
  // return 0;

  logger.raw().dump("buf", (uint8_t *)buff, 210);
  logger.err().log("test", "Hello Johny %d times", 5);
  logger.dbg().log("test", "Hello Johny %d times", 5);

  const size_t PERSISTENCE_SIZE = 256;
  static uint8_t memory[PERSISTENCE_SIZE];
  static constexpr persistency::Persistency<infra::Config::PersistencyId>::Row persistency_table[] = {
      persistency::Persistency<infra::Config::PersistencyId>::make_persistency_row<infra::WifiSettingsN>(infra::Config::PersistencyId::WIFI_CONFIG_SOFTAP),
      persistency::Persistency<infra::Config::PersistencyId>::make_persistency_row<infra::WifiSettingsN>(infra::Config::PersistencyId::WIFI_CONFIG_STATION),
      {infra::Config::PersistencyId::_LAST, 0, nullptr},
  };
  static_assert(
      persistency::Persistency<infra::Config::PersistencyId>::check_persistency_table_size(
          persistency_table,
          PERSISTENCE_SIZE),
      "Error, insufficient memory for persistency table mapping");

  persistency::Persistency<infra::Config::PersistencyId> storage(persistency_table, memory, sizeof(memory), Flash(logger));
  storage.load();
  infra::Config config(storage);

  MyGpio coin_pin("coin-ctrl", 1, logger);
  MyGpio led_coin("coin-led", 2, logger);
  MyGpio led_wifi_softap("wifi-softap-led", 3, logger);
  MyGpio led_wifi_station("wifi-station-led", 4, logger);
  MyGpio led_wifi_station_connected("wifi-softap-connected-led", 5, logger);

  Wifi wifi(logger);
  wifi.swith_to_softap();
  // wifi.swith_to_station();

  infra::Status status(logger, ts,
                       led_coin,
                       led_wifi_softap,
                       led_wifi_station,
                       led_wifi_station_connected);
  Paho broker(config, logger);
  service::mode::ModeService modsvc(logger, wifi, status);

  Button btn(modsvc);
  strategy::StepWait sw100(500, ts);
  strategy::StepActuatorOff saof(coin_pin);
  strategy::StepActuatorOn saon(coin_pin);
  strategy::IStep *st[] = {
      &saon,
      &sw100,
      &saof,
  };
  strategy::StepRunner step_runner(st, sizeof(st) / sizeof(st[0]));
        // constexpr CoinService(IPortAdapterBroker &broker,
        //                     IPortAdapterStatus &status,
        //                     const strategy::IScenarioMaker &scenario_maker,
        //                     const core::logger::ILogger &logger,
        //                     const core::ITimestamp &ts,
        //                     IPortAdapterConfig &config)
  
  service::coin::CoinService<COIN_SCV_CAP> coinsvc(broker, status, ScenarioMaker(ts,logger,config), logger, ts, config);
  broker.add_subscriber(&coinsvc, coinsvc.sub_topic());
  // broker.add_subscriber(&coinsvc, config.root_sub_topic());

  service::web::WebService ws(logger, config);

  ws.add_subscriber(&modsvc, service::web::event::EventWifiCredSaved::event_name());
  modsvc.add_subscriber(&ws, service::mode::event::EventWifiModeChanged::event_name());
  modsvc.add_subscriber(&coinsvc, service::mode::event::EventWifiModeChanged::event_name());

  platform::router::Server server;
  platform::router::Router r(server, ws, logger);

  // std::thread t1(button_run, &btn);

  std::thread t2(_run, &modsvc, &status, &coinsvc);
  r.route().listen("localhost", 8082);

  // for (;;)
  // {
  //   std::this_thread::sleep_for(std::chrono::milliseconds(1000));
  //   // logger.dbg().log("main", "Clicked: %d, Pressed: %d", btn.clicked(), btn.pressed());
  // }
  return 0;
}