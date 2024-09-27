#pragma once
#include <string.h>
#include <framework/core/logger.h>
#include "i_broker_client.h"
#include "i_broker_listener.h"

namespace broker
{
  template <size_t N>
  class Broker : public IBrokerClient, public IBrokerConnectionListener
  {
    struct Subscriber
    {
      constexpr Subscriber() : listener(nullptr), topic(nullptr), qos(0)  {}
      constexpr Subscriber(const char *topic, IBrokerListener *listener, uint32_t qos=0)
          : listener(listener), topic(topic), qos(0) {}
      void foor() {}
      IBrokerListener *listener;
      const char *topic;
      uint32_t qos;
    };

  public:
    constexpr Broker(const core::logger::ILogger &logger) : _index(0), _logger(logger) {}

    virtual ~Broker() {}

    virtual void disconnected(const char *reason) override
    {
      _logger.inf().log(TAG, "Disconnected: %s", reason);
    }

    virtual void connected() override
    {
      _logger.inf().log(TAG, "Connected");
    }

    virtual void arrived(const char *topic, const Message &msg) override
    {
      // printf("1\n");
      _logger.dbg().log(TAG, "Arrived message on topic: %s", topic);

      for (size_t i = 0; i < _index; i++)
      {
        if (strcmp(_subs[i].topic, topic) == 0)
        {
          // _logger.dbg().log(TAG, "Notify");

          _subs[i].listener->notify(topic, msg);
        }
      }
    }

    virtual void delivered(const Token &token) override
    {
      _logger.inf().log(TAG, "Delivered message with token: %d", token);
    }

    virtual bool add_subscriber(broker::IBrokerListener *listener, const char *topic, uint32_t qos = 0) override
    {
      if (_index < N)
      {
        if (find_subscriber(listener, topic) == nullptr)
        {
          _logger.inf().log(TAG, "Add subscriber qos:%d on topic: %s", qos, topic);
          _subs[_index++] = Subscriber(topic, listener, qos);
          return true;
        }
        else
        {
          _logger.wrn().log(TAG, "Add subscriber failed due to duplication of topic/listener: %s", topic);
          return false;
        }
      }
      else
      {
        _logger.err().log(TAG, "Failed to subscribe, due to limit of subscribers has been reached %d/%d", _index, N);
      }
      return false;
    }

    // virtual bool connect() const override
    // {
    //   return false;
    // }
    // virtual bool disconnect() const override
    // {
    //   return false;
    // }
    // virtual bool is_connected() const override
    // {
    //   return false;
    // }

    // virtual bool publish(const char *topic, const Message &msg) override
    // {
    //   return false;
    // }

  private:
    Subscriber *find_subscriber(broker::IBrokerListener *listener, const char *topic)
    {
      for (size_t i = 0; i < _index; i++)
      {
        if ((_subs[i].listener == listener) &&
            (_subs[i].topic == topic))
        {
          return &_subs[i];
        }
      }
      return nullptr;
    }

  protected:
    const core::logger::ILogger &_logger;
    Subscriber _subs[N];
    size_t _index;
    static constexpr const char *TAG = "broker";
  };
}
