#pragma once
#include <stdint.h>
#include "message.h"
namespace broker
{
  struct IBrokerListener
  {
    virtual void notify(const char *topic, const Message& msg) = 0;
  };

} // namespace broker
