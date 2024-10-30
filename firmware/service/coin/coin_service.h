#pragma once
#include <framework/persistency.h>
#include <framework/core/observer.h>
#include <framework/strategy.h>
#include <framework/core/i_runnable.h>
#include <framework/core/i_timestamp.h>
#include "messages/heartbeat.h"
#include "messages/command_composite.h"
#include "messages/response_ack.h"
#include "messages/response_config.h"
#include "messages/response_scenario_completed.h"
#include "port_coin_adapter_broker.h"
#include "port_coin_adapter_status.h"
#include "port_coin_adapter_config.h"
#include "port_coin_controller_broker.h"

/**
 * mosquitto_pub -h 192.168.0.150 -p 1883 -t '/clousel/cloud/550e8400-e29b-41d4-a716-446655440000' -m '{"Type":"MessageCommand","CarouselId":"550e8400-e29b-41d4-a716-446655440000","SequenceNum":1,"EventId":"cedb3510-c87f-4f7d-a190-2f1f8412ff29","Command":"ConfigWrite", "Config":{"BrokerUrl":"mqtt://192.168.0.150:1883","BrokerUsername":"CLOUSEL","BrokerPassword":"123wqeqwsaddsa","CoinPulseCnt":1,"CoinPulseDur":100}}'
 * mosquitto_pub -h 192.168.0.150 -p 1883 -t '/clousel/cloud/550e8400-e29b-41d4-a716-446655440000' -m '{"Type":"MessageCommand","CarouselId":"550e8400-e29b-41d4-a716-446655440000","SequenceNum":1,"EventId":"cedb3510-c87f-4f7d-a190-2f1f8412ff29","Command":"ConfigRead"}'
 */
namespace service
{
  namespace coin
  {
    template <size_t N>
    class CoinService
        : public IPortControllerBroker,
          public core::IRunnable,
          public core::observer::IListener

    {
    public:
      constexpr CoinService(IPortAdapterBroker &broker,
                            IPortAdapterStatus &status,
                            const strategy::IScenarioMaker &scenario_maker,
                            const core::logger::ILogger &logger,
                            const core::ITimestamp &ts,
                            IPortAdapterConfig &config)
          : _broker(broker),
            _status(status),
            _scenario_maker(scenario_maker),
            _scenario(scenario_maker.make()),
            _logger(logger),
            _config(config),
            _ts(ts),
            _enabled(false),
            _message_inomming_len(0)
      {
        _last_cmd.clear();
      }

      virtual void notify(const broker::ITopic &topic, const broker::Message &msg) override
      {
        _logger.inf().log(TAG, "Arrived message %d bytes from topic %.*s", msg.size, topic.len(), topic.get());
        for (_message_inomming_len = 0; _message_inomming_len < msg.size; _message_inomming_len++)
        {
          if (_message_inomming_len < N)
          {
            _message_buffer[_message_inomming_len] = msg.data[_message_inomming_len];
          }
          else
          {
            _logger.err().log(TAG, "Not enought space for saving message %d/%d", msg.size, N);
            _message_inomming_len = 0;
            break;
          }
        }
        _logger.dbg().log(TAG, "Got message (%d): %.*s", _message_inomming_len, _message_inomming_len, _message_buffer);
        _logger.inf().log(TAG, "Copied message of %d bytes", _message_inomming_len);
      }

      virtual void run() override
      {
        handle_incoming_message();
        control_broker_connection();
        control_scenario_execution();
        hearbeat();
      }
      inline broker::TopicContainer<100> pub_topic() const
      {
        broker::TopicContainer<100> topic(_config.root_pub_topic());
        topic.append(_config.carousel_id()->value());
        return topic;
      }
      inline broker::TopicContainer<100> sub_topic() const
      {
        broker::TopicContainer<100> topic(_config.root_sub_topic());
        topic.append(_config.carousel_id()->value());
        return topic;
      }

    private:
      virtual void notify(const core::observer::Event &event) override
      {
        if (event.name() == mode::event::EventWifiModeChanged::event_name())
        {
          const mode::event::EventWifiModeChanged &ewmc = reinterpret_cast<const mode::event::EventWifiModeChanged &>(event);
          _enabled = ewmc.is_station();
          _logger.inf().log(TAG, "Has got a new WifiMode event. Service is now %s", _enabled ? "enabled" : "disabled");
        }
      }

      void control_broker_connection()
      {
        if (_enabled && _broker.is_ready() && !_broker.is_connected())
        {
          _broker.connect();
        }
        else if (!_enabled && _broker.is_connected())
        {
          _broker.disconnect();
        }
      }

      void handle_incoming_message()
      {
        const char *err = nullptr;
        bool reconnect = false;
        if (_message_inomming_len > 0)
        {
          do
          {
            if (!_last_cmd.parse(_message_buffer, _message_inomming_len))
            {
              err = "Fail to parse json message";
              _logger.err().log(TAG, err);
              break;
            }
            if (!_last_cmd.general.carousel_id.value.eq(_config.carousel_id()->value()))
            {
              err = "Carousel id mismatch";
              _logger.err().log(TAG, "%s, got:%s, own:%s", _last_cmd.general.carousel_id.value.data(), _config.carousel_id()->value());
              break;
            }

            if (!_last_cmd.general.type.value.eq(BMT_COMMAND))
            {
              err = "Unexpected message type";
              _logger.err().log(TAG, "%s: '%s'", err, _last_cmd.general.type.value.data());
              break;
            }

            _logger.inf().log(TAG, "Message is command message with exact command: %s", _last_cmd.general.command.value.data());

            if (_last_cmd.general.command.value.eq(BMTC_PLAY))
            {
              _logger.inf().log(TAG, "Scenario has been reseted");
              _scenario.reset();
              _status.led_coin_blink();
              publish(msg::ResponseAck(_last_cmd.general.carousel_id.value.data(), _last_cmd.general.event_id.value.data(), nullptr, sequence_num()));
            }
            else if (_last_cmd.general.command.value.eq(BMTC_CONFIG_WRITE))
            {
              if (!_last_cmd.config.value.broker_url.value.empty())
              {
                _logger.inf().log(TAG, "Set Broker URL: '%s'", _last_cmd.config.value.broker_url.value.data());
                _config.set_broker_url(_last_cmd.config.value.broker_url.value);
                reconnect = true;
              }
              if (!_last_cmd.config.value.broker_username.value.empty())
              {
                _logger.inf().log(TAG, "Set Broker Username: '%s'", _last_cmd.config.value.broker_username.value.data());
                _config.set_broker_username(_last_cmd.config.value.broker_username.value);
                reconnect = true;
              }
              if (!_last_cmd.config.value.broker_password.value.empty())
              {
                _logger.inf().log(TAG, "Set Broker Password: '%s'", _last_cmd.config.value.broker_password.value.data());
                _config.set_broker_password(_last_cmd.config.value.broker_password.value);
                reconnect = true;
              }
              infra::CoinPulseProps props(_last_cmd.config.value.coin_pulse_count.value, _last_cmd.config.value.coin_pulse_duration.value);
              if (props.is_valid())
              {
                _config.set_coin_pulse_props(props);
                _scenario = _scenario_maker.make();
                publish(msg::ResponseAck(_last_cmd.general.carousel_id.value.data(), _last_cmd.general.event_id.value.data(), nullptr, sequence_num()));
              }
              else
              {
                err = "Invalid coin pulse data";
              }
              if (reconnect)
              {
                _broker.disconnect();
              }
              if (!_config.save())
              {
                err = "Fail to save new configuration";
              }
            }
            else if (_last_cmd.general.command.value.eq(BMTC_CONFIG_READ))
            {
              auto *props = _config.coin_pulse_props();
              publish(msg::ResponseConfig(
                  _last_cmd.general.carousel_id.value.data(),
                  _last_cmd.general.event_id.value.data(),
                  nullptr,
                  msg::Config(props->count, props->duration, _config.broker_url()),
                  sequence_num()));
            }
            else
            {
              err = "Unexpected command";
              _logger.err().log(TAG, "%s: '%s'", err, _last_cmd.general.command.value.data());
              break;
            }
          } while (false);
          if (err)
          {
            _logger.err().log(TAG, "About to publish ack with an error: %s", err);
            publish(msg::ResponseAck(_last_cmd.general.carousel_id.value.data(),
                                     _last_cmd.general.event_id.value.data(),
                                     err,
                                     sequence_num()));
          }
          _last_cmd.clear();
          _message_inomming_len = 0;
        }
        if (reconnect)
        {
          bool res = true;
          res = res && _broker.disconnect();
          if (!res)
          {
            _logger.err().log(TAG, "Fail to disconnect");
          }
          // res = res && _broker.reinit();
          // if (!res)
          // {
          //   _logger.err().log(TAG, "Fail to reinit");
          // }
        }
      }

      void control_scenario_execution()
      {
        if (!_scenario.finished())
        {
          _scenario.run();
          if (_scenario.finished())
          {
            _logger.inf().log(TAG, "Scenario finished");
            publish(msg::ResponseScenarioCompleted(_last_cmd.general.carousel_id.value.data(), _last_cmd.general.event_id.value.data(), nullptr, sequence_num()));
          }
        }
      }

      void hearbeat()
      {
        static size_t tsp = 0;
        if (_broker.is_connected() && _ts.get() > tsp)
        {
          tsp = _ts.get() + _config.heartbeat_tm_ms();
          publish(msg::Heartbeat(_config.carousel_id()->value(), sequence_num()));
        }
      }

      template <typename T>
      bool publish(const T &msg)
      {
        bool res = false;
        broker::TopicContainer<100> pub_topic(_config.root_pub_topic());
        // char json_buffer[N];
        do
        {
          if (!_broker.is_connected())
          {
            _logger.wrn().log(TAG, "Can't publish, broker isn't connected");
            break;
          }
          if (!pub_topic.append(_config.carousel_id()->value()))
          {
            _logger.err().log(TAG, "Fail append to topic");
            break;
          }
          if (!msg.dump(_message_buffer, sizeof(_message_buffer)))
          {
            _logger.err().log(TAG, "Fail to dump message");
            break;
          }
          _logger.inf().log(TAG, "Publishing %s", _message_buffer);
          if (!_broker.publish(pub_topic, _message_buffer).is_ok())
          {
            _logger.err().log(TAG, "Fail to publish message to topic: %s", pub_topic.get());
            break;
          }
          res = true;
        } while (false);
        return res;
      }
      inline size_t sequence_num() const
      {
        static size_t _sequence_num = 0;
        return ++_sequence_num;
      }

    private:
      IPortAdapterBroker &_broker;
      IPortAdapterStatus &_status;
      const strategy::IScenarioMaker& _scenario_maker;
      strategy::StepRunner &_scenario;
      const core::logger::ILogger &_logger;
      IPortAdapterConfig &_config;
      const core::ITimestamp &_ts;
      bool _enabled;
      size_t _message_inomming_len;
      char _message_buffer[N];
      msg::CommandComposite _last_cmd;

    public:
      static constexpr const char *TAG = "coin";
    };
  }
}