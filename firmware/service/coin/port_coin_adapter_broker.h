#pragma once
#include <framework/broker.h>

namespace service
{
  namespace coin
  {
    struct IPortAdapterBroker
    {
      virtual bool is_connected() const = 0;
      virtual bool publish(const char *topic, const broker::Message& msg, uint32_t qos = 0) = 0;
    };
  }
}