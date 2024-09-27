#pragma once
#include <stdint.h>

namespace core
{
  namespace io
  {
    template<typename T>
    struct IActuator
    {
      virtual uint32_t id() const = 0;
      virtual void set(const T& value) = 0;
      virtual T get() const = 0;
    };

  }
}