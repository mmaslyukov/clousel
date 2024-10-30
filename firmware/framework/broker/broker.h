#pragma once
#include <string.h>
#include <framework/core/logger.h>
#include "i_broker_client.h"
#include "i_broker_listener.h"
#include "topic.h"

namespace broker
{
  template <size_t N>
  class Broker : public IBrokerClient, public IBrokerConnectionListener
  {
    struct Subscriber
    {
      constexpr Subscriber() : listener(nullptr), topic(), qos(0) {}
      constexpr Subscriber(const ITopic &topic, IBrokerMessageListener *listener, uint32_t qos = 0)
          : listener(listener), topic(topic.get(), topic.len()), qos(0) {}
      void foor() {}
      IBrokerMessageListener *listener;
      TopicContainer<100> topic;
      // const char *topic;
      uint32_t qos;
    };

    struct Delivery
    {
      constexpr Delivery() : token(), listener(nullptr) {}
      constexpr Delivery(Token token, IBrokerDeliveryListener *listener) : token(token), listener(listener) {}

      bool is_valid()
      {
        return listener != nullptr && token.is_valid();
      }

      void invalidate()
      {
        listener = nullptr;
        token.invalidate();
      }
      Token token;
      IBrokerDeliveryListener *listener;
    };

  public:
    constexpr Broker(const core::logger::ILogger &logger) : _logger(logger), _index_sub(0) {}

    virtual ~Broker() {}

    virtual void disconnected(const char *reason) override
    {
      _logger.inf().log(TAG, "Disconnected: %s", reason);
    }

    virtual void connected() override
    {
      _logger.inf().log(TAG, "Connected");
    }

    virtual void arrived(const ITopic &topic, const Message &msg) override
    {
      _logger.dbg().log(TAG, "Arrived message on topic: %.*s", topic.len(), topic.get());
      for (size_t i = 0; i < _index_sub; i++)
      {
        if (_subs[i].topic.part_of(topic))
        {
          _subs[i].listener->notify(topic, msg);
        }
      }
    }

    virtual void delivered(const Token &token) override
    {
      _logger.inf().log(TAG, "Delivered message with token: %d", token.id());
      for (size_t i = 0; i < sizeof(_deliveries) / sizeof(_deliveries[0]); i++)
      {
        if (_deliveries[i].is_valid() && (_deliveries[i].token.id() == token.id()) && _deliveries[i].listener)
        {
          _logger.dbg().log(TAG, "del: deliveries index:%d", i);
          _deliveries[i].listener->delivered(token);
          _deliveries[i].invalidate();
          break;
        }
      }
    }

    virtual bool add_subscriber(broker::IBrokerMessageListener *listener, const ITopic &topic, uint32_t qos = 0) override
    {
      if (_index_sub < N)
      {
        if (find_subscriber(listener, topic) == nullptr)
        {
          _logger.inf().log(TAG, "Add subscriber qos:%d on topic: %.*s", qos, topic.len(), topic.get());
          _subs[_index_sub++] = Subscriber(topic, listener, qos);
          return true;
        }
        else
        {
          _logger.wrn().log(TAG, "Add subscriber failed due to duplication of topic/listener: %.*s", topic.len(), topic.get());
          return false;
        }
      }
      else
      {
        _logger.err().log(TAG, "Failed to subscribe, due to limit of subscribers has been reached %d/%d", _index_sub, N);
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

    virtual Token publish_confirm(const ITopic &topic, const Message &msg, IBrokerDeliveryListener *delivery_listener, uint32_t qos = 0) override
    {
      Token token = publish(topic, msg, qos);
      for (size_t i = 0; i < sizeof(_deliveries) / sizeof(_deliveries[0]); i++)
      {
        if (!_deliveries[i].is_valid())
        {
          _logger.dbg().log(TAG, "pub: deliveries index:%d", i);
          _deliveries[i].token = token;
          _deliveries[i].listener = delivery_listener;
          break;
        }
      }
      return token;
    }

  private:
    Subscriber *find_subscriber(broker::IBrokerMessageListener *listener, const ITopic &topic)
    {
      for (size_t i = 0; i < _index_sub; i++)
      {
        if ((_subs[i].listener == listener) &&
            (strcmp(_subs[i].topic.get(), topic.get()) == 0))
        {
          return &_subs[i];
        }
      }
      return nullptr;
    }

  protected:
    const core::logger::ILogger &_logger;
    size_t _index_sub;
    Subscriber _subs[N];
    Delivery _deliveries[10];

  private:
    static constexpr const char *TAG = "broker";
  };
}
