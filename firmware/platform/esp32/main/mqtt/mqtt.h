#pragma once

#include <framework/broker.h>
#include <framework/core/logger.h>

#include "mqtt_client.h"

#define BROKER_SUBSCRIBERS 1
class Mqtt
    : public broker::Broker<BROKER_SUBSCRIBERS>,
      public service::coin::IPortAdapterBroker,
      public IEspEvetHandler
{
  enum ConnectionState
  {
    INITIALIZED,
    INITIALIZING,
    UNINITIALIZED,
    UNINITIALIZING,
    DISCONNECTED,
    DISCONNECTING,
    CONNECTING,
    CONNECTED,
  };

public:
  constexpr Mqtt(const WifiModeStation &sta, const infra::Config &config, const core::logger::ILogger &logger)
      : broker::Broker<BROKER_SUBSCRIBERS>(logger), _wifi_sta(sta), _config(config), _client(nullptr), _connst(ConnectionState::UNINITIALIZED)
  {
    // init();
  }
  // virtual bool reinit() override
  // {
  //   esp_mqtt_client_destroy(_client);
  //   _client = nullptr;
  //   return init();
  // }

  bool connect() override
  {
    bool res = true;
    do
    {
      if (_connst == ConnectionState::CONNECTING)
      {
        res = true;
        break;
      }
      if (_connst == ConnectionState::CONNECTED)
      {
        _logger.wrn().log(TAG, "Already connected %d", _connst);
        res = true;
        break;
      }
      res = res && _wifi_sta.is_enabled();
      if (!res)
      {
        _logger.wrn().log(TAG, "Fail, wifi isn't in the station mode");
        break;
      }
      _logger.dbg().log(TAG, "Connecting to '%s'", _config.broker_url()->data());
      res = res && _wifi_sta.is_connected();
      if (!res)
      {
        _logger.wrn().log(TAG, "Fail, wifi isn't connected");
        break;
      }
      if (!_client)
      {
        res = res && init();
        if (!res)
        {
          _logger.err().log(TAG, "Fail to init to the the client");
          break;
        }
      }
      res = res && esp_mqtt_client_reconnect(_client) == ESP_OK;
      if (!res)
      {
        _logger.err().log(TAG, "Fail to connect to the broker");
      }
      _connst = ConnectionState::CONNECTING;
      _logger.dbg().log(TAG, "Switch to CONNECTING state");

    } while (false);
    return res;
  }

  bool disconnect() override
  {
    bool res = true;
    if (_connst == ConnectionState::DISCONNECTED || _connst == ConnectionState::DISCONNECTING)
    {
      return true;
    }
    _logger.dbg().log(TAG, "Disconnecting from '%s'", _config.broker_url()->data());
    res = res && esp_mqtt_client_disconnect(_client) == ESP_OK;
    if (!res)
    {
      _logger.err().log(TAG, "Fail to disconnect");
    }
    res = res && esp_mqtt_client_stop(_client) == ESP_OK;
    if (!res)
    {
      _logger.err().log(TAG, "Fail to stop");
    }
    if (_client)
    {
      res = res && esp_mqtt_client_destroy(_client) == ESP_OK;
      if (res)
      {
        _client = nullptr;
        _logger.inf().log(TAG, "Client has been destroyed");
      }
      else
      {
        _logger.err().log(TAG, "Fail to destroy the client");
      }
    }
    _connst = ConnectionState::DISCONNECTING;
    _logger.dbg().log(TAG, "Switch to DISCONNECTING state");
    return res;
  }

  virtual bool is_connected() const override
  {
    // esp_mqtt_cilent
    return _connst == ConnectionState::CONNECTED;
  }

  virtual bool is_ready() const override
  {
    return _wifi_sta.is_connected();
  }

  virtual broker::Token publish(const broker::ITopic &topic, const broker::Message &msg, const uint32_t qos = 0) override
  {
    int msg_id = esp_mqtt_client_publish(_client, topic.get(), (const char *)msg.data, msg.size, qos, 0);
    return broker::Token(true, msg_id);
  }

private:
  virtual void connected() override
  {
    _connst = ConnectionState::CONNECTED;
    for (size_t i = 0; i < _index_sub; i++)
    {
      int msg_id = esp_mqtt_client_subscribe(_client, _subs[i].topic.get(), _subs[i].qos);
      _logger.inf().log(TAG, "Subscribed to '%s' id:%d", _subs[i].topic.get(), msg_id);
    }
  }

  virtual void disconnected(const char *reason) override
  {
    _connst = ConnectionState::DISCONNECTED;
    _logger.inf().log(TAG, "Disconnected due to %s", reason);
    for (size_t i = 0; i < _index_sub; i++)
    {
      int msg_id = esp_mqtt_client_unsubscribe(_client, _subs[i].topic.get());
      _logger.inf().log(TAG, "Unsubscribed from '%s' id:%d", _subs[i].topic.get(), msg_id);
    }
  }

  virtual bool handle(esp_event_base_t event_base, int32_t event_id, void *event_data) override
  {
    bool res = true;
    esp_mqtt_event_handle_t event = reinterpret_cast<esp_mqtt_event_handle_t>(event_data);
    _logger.dbg().log(TAG, "event_id:%d", event_id);
    switch ((esp_mqtt_event_id_t)event_id)
    {
    case MQTT_EVENT_CONNECTED:
      _logger.inf().log(TAG, "CONNECTED to '%s'", _config.broker_url()->data());
      connected();
      break;
    case MQTT_EVENT_DISCONNECTED:
      _logger.inf().log(TAG, "DISONNECTED from '%s'", _config.broker_url()->data());
      disconnected("MQTT_EVENT_DISCONNECTED");
      break;
    case MQTT_EVENT_SUBSCRIBED:
      break;
    case MQTT_EVENT_UNSUBSCRIBED:
      break;
    case MQTT_EVENT_PUBLISHED:
      delivered(broker::Token(true, event->msg_id));
      break;
    case MQTT_EVENT_DATA:
    {
      broker::TopicRef topic(event->topic, event->topic_len);
      arrived(topic, broker::Message(event->data, event->data_len));
    }
    break;
    case MQTT_EVENT_BEFORE_CONNECT:
      _logger.inf().log(TAG, "About to connect to  '%s'", _config.broker_url()->data());
      break;
    case MQTT_EVENT_ERROR:
      _logger.err().log(TAG, "Got error, last errno: %s", strerror(event->error_handle->esp_transport_sock_errno));
      break;
    default:
      res = false;
      break;
    }
    return res;
  }
  static void mqtt_event_handler(void *arg, esp_event_base_t event_base, int32_t event_id, void *event_data)
  {
    if (arg)
    {
      reinterpret_cast<Mqtt *>(arg)->handle(event_base, event_id, event_data);
    }
    else
    {
      ESP_LOGE(TAG, "arg is null");
    }
  }

private:
  bool init()
  {
    bool res = true;
    esp_mqtt_client_config_t mqtt_cfg;
    memset(&mqtt_cfg, 0, sizeof(mqtt_cfg));
    mqtt_cfg.broker.address.uri = _config.broker_url()->data();
    mqtt_cfg.credentials.client_id = _config.broker_client_id();
    _client = esp_mqtt_client_init(&mqtt_cfg);
    if (_client)
    {
      _logger.inf().log(TAG, "Cliet has been created");
    }
    res = res && esp_mqtt_client_register_event(_client, MQTT_EVENT_ANY, &mqtt_event_handler, this) == ESP_OK;
    res = res && esp_mqtt_client_start(_client) == ESP_OK;
    return res;
  }

private:
  const WifiModeStation &_wifi_sta;
  const infra::Config &_config;
  esp_mqtt_client_handle_t _client;
  // esp_mqtt_client_config_t _mqtt_cfg;
  ConnectionState _connst;
  static constexpr const char *TAG = "broker.mqtt";
};