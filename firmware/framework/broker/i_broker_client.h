#pragma once
#include <stdint.h>
#include "i_broker_listener.h"
#include "message.h"
#include "topic.h"

namespace broker
{
  struct IBrokerConnectionListener
  {
    virtual ~IBrokerConnectionListener() {}
    virtual void disconnected(const char *reason) = 0;
    virtual void connected() = 0;
    virtual void arrived(const ITopic &topic, const Message &msg) = 0;
    virtual void delivered(const Token &token) = 0;
  };

  struct IBrokerClient
  {
    virtual ~IBrokerClient() {}
    virtual bool connect() = 0;
    virtual bool disconnect() = 0;
    virtual bool is_connected() const = 0;
    virtual Token publish(const ITopic &topic, const Message &msg, uint32_t qos) = 0;
    virtual Token publish_confirm(const ITopic &topic, const Message &msg, IBrokerDeliveryListener *delivery_listener, uint32_t qos) = 0;
    virtual bool add_subscriber(IBrokerMessageListener *listener, const ITopic &topic, const uint32_t qos) = 0;
  };
}
