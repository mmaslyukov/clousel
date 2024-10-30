#pragma once
#include <stdint.h>
#include "message.h"
#include "topic.h"

namespace broker
{
  struct IBrokerMessageListener
  {
    virtual void notify(const ITopic &topic, const Message& msg) = 0;
  };
  struct IBrokerDeliveryListener
  {
    virtual void delivered(const Token& token) = 0;
  };
} // namespace broker
