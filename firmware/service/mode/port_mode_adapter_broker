#pragma once
#include <framework/broker.h>

namespace service
{
  namespace mode
  {
    struct IPortAdapterBroker
    {
      virtual bool connect() = 0;
      virtual bool disconnect() = 0;
      virtual bool is_connected() const = 0;
      virtual bool publish(const char *topic, const broker::Message& msg, uint32_t qos = 0) = 0;
    };

  }
}