#pragma once
#include <stdint.h>

namespace core
{
  namespace io
  {
    struct IButton
    {
      virtual uint32_t id() const = 0;
      virtual bool clicked() const = 0;
      virtual size_t pressed() const = 0;
    };

  }
}