#pragma once
#include <stdint.h>
#include "verbosity.h"
#include "i_enabable.h"

namespace core
{
  namespace logger
  {
    
    struct IDumpable : public IEnabable
    {
      virtual void dump(const char *tag, const uint8_t* data, size_t size) const = 0;
      virtual void dump_ascii(const char *tag, const uint8_t* data, size_t size) const = 0;
    };
    
  }
}