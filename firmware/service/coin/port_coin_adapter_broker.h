#pragma once
#include <framework/broker.h>

namespace service
{
  namespace coin
  {
    struct IPortAdapterBroker
    {
      virtual bool connect() = 0;
      virtual bool disconnect() = 0;
      virtual bool is_connected() const = 0;
      virtual bool is_ready() const = 0;
      // virtual bool reinit() = 0;
      virtual broker::Token publish(const broker::ITopic &topic, const broker::Message& msg, uint32_t qos = 0) = 0;
    };
  }
}