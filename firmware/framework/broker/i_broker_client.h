#pragma once
#include <stdint.h>
#include "i_broker_listener.h"
#include "message.h"

namespace broker
{
  struct IBrokerConnectionListener
  {
    virtual ~IBrokerConnectionListener() {}
    virtual void disconnected(const char* reason) = 0;
    virtual void connected() = 0;
    virtual void arrived(const char *topic, const Message& msg) = 0;
    virtual void delivered(const Token& token) = 0;
  };

  struct IBrokerClient
  {
    virtual ~IBrokerClient() {}
    virtual bool connect()  = 0;
    virtual bool disconnect()  = 0;
    virtual bool is_connected() const = 0;
    virtual bool publish(const char *topic, const Message& msg, uint32_t qos = 0) = 0;
    virtual bool add_subscriber(broker::IBrokerListener *listener, const char *topic, const uint32_t qos) = 0;
  };
}
