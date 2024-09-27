#pragma once
#include <stdint.h>

namespace core
{
  struct ITimestamp
  {
    virtual size_t get() const = 0;
  };
}