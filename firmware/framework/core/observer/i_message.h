#pragma once
#include <stdint.h>

namespace core
{
  namespace observer
  {

    struct IMessage
    {
      virtual const char *name() const = 0;
    };
  }
}
