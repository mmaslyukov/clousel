#pragma once
#include <stdint.h>
#include "event.h"
#include "i_listener.h"

namespace core
{
  namespace observer
  {
    struct ISubscribable
    {
      virtual bool add_subscriber(IListener *listener, const EventBase &event) = 0;
    };
  }
}
