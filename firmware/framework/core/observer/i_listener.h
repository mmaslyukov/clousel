#pragma once
#include <stdint.h>
#include "event.h"

namespace core
{
  namespace observer
  {
    struct IListener
    {
      virtual void notify(const Event &event) = 0;
    };
  }
}
