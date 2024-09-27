#pragma once
#include <framework/persistency.h>
#include <framework/strategy.h>
#include <framework/core/i_runnable.h>
#include <framework/core/i_timestamp.h>
#include "messages/heartbeat.h"
#include "messages/command.h"
#include "messages/response_ack.h"
#include "port_coin_adapter_broker.h"
#include "port_coin_adapter_status.h"
#include "port_coin_adapter_config.h"
#include "port_coin_controller_broker.h"

namespace service
{
  namespace coin
  {
    template <size_t N>
    class CoinService
        : public IPortControllerBroker,
          public core::IRunnable
    {
    public:
      constexpr CoinService(IPortAdapterBroker &broker,
                            IPortAdapterStatus &status,
                            strategy::StepRunner &scenario,
                            const core::logger::ILogger &logger,
                            const core::ITimestamp &ts,
                            const IPortAdapterConfig &config)
          : _broker(broker),
            _status(status),
            _scenario(scenario),
            _logger(logger),
            _config(config),
            _ts(ts),
            _message_inomming_len(0) {}

      virtual void notify(const char *topic, const broker::Message &msg) override
      {
        _logger.inf().log(TAG, "Arrived message from topic %s", topic);
        broker::Topic<100> sub(_config.root_sub_topic());
        // sub.append(_config.carousel_id());
        if (!sub.part_of(topic))
        {
          return;
        }
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
      }

      virtual void run() override
      {
        scenario();
        hearbeat();
      }
      inline broker::Topic<100> pub_topic() const
      {
        broker::Topic<100> topic(_config.root_pub_topic());
        topic.append(_config.carousel_id()->value());
        return topic;
      }
      inline broker::Topic<100> sub_topic() const
      {
        broker::Topic<100> topic(_config.root_sub_topic());
        topic.append(_config.carousel_id()->value());
        return topic;
      }
    private:
      void scenario()
      {
        const char *err = nullptr;
        if (_message_inomming_len > 0)
        {
          do
          {
            if (!_last_cmd.parse(_message_buffer))
            {
              err = "Fail to parse json message";
              _logger.err().log(TAG, err);
              break;
            }
            if (!_last_cmd.type.value.eq(BMT_COMMAND))
            {
              err = "Unexpected message type";
              _logger.err().log(TAG, "%s: %s", err, _last_cmd.type.value.data());
              break;
            }
            if (!_last_cmd.command.value.eq(BMTC_PLAY))
            {
              err = "Unexpected command";
              _logger.err().log(TAG, "%s: %s", err, _last_cmd.command.value.data());
              break;
            }
            _logger.dbg().log(TAG, "Scenario reset");
            _scenario.reset();
            _status.led_coin_blink();
          } while (false);
          _message_inomming_len = 0;
        }
        if (err)
        {
          publish(msg::ResponseAck(_last_cmd.carousel_id.value.data(), _last_cmd.event_id.value.data(), err, sequence_num()));
        }

        if (!_scenario.finished())
        {
          _scenario.run();
          if (_scenario.finished())
          {
            _logger.dbg().log(TAG, "Scenario finished");
            publish(msg::ResponseAck(_last_cmd.carousel_id.value.data(), _last_cmd.event_id.value.data(), nullptr, sequence_num()));
          }
        }
      }

      void hearbeat()
      {
        static size_t tsp = 0;
        if (_ts.get() > tsp + _config.heartbeat_tm_ms() || !tsp)
        {
          tsp = _ts.get();
          // publish(msg::Heartbeat(_config.carousel_id()->value()));
          publish(msg::Heartbeat(_config.carousel_id()->value(), sequence_num()));
        }
      }

      template <typename T> bool publish(const T &msg) 
      {
        bool res = false;
        broker::Topic<100> pub_topic(_config.root_pub_topic());
        char json_buffer[200];// = _message_buffer; //[200];
        do
        {
          if (!pub_topic.append(_config.carousel_id()->value()))
          {
            _logger.err().log(TAG, "Fail append to topic");
            break;
          }
          // msg.sequence_num.value = sequence_num();
          if (!msg.dump(json_buffer, sizeof(json_buffer)))
          {
            _logger.err().log(TAG, "Fail to dump message");
            break;
          }
          _logger.dbg().log(TAG,  "Publishing %s", json_buffer);
          if (!_broker.publish(pub_topic.get(), json_buffer))
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
      strategy::StepRunner &_scenario;
      const core::logger::ILogger &_logger;
      const core::ITimestamp &_ts;
      const IPortAdapterConfig &_config;
      msg::Command _last_cmd;
      size_t _message_inomming_len;
      char _message_buffer[N];
      // broker::Topic<100> _topic;
      static constexpr const char *TAG = "coin";
    };
  }
}