#pragma once
#include <stdint.h>
#include "verbosity.h"
#include "i_printable.h"
#include "i_dumpable.h"

namespace core
{
  namespace logger
  {
   struct ILoggerSystem
    {
      virtual void output(const Verbosity &verbosity, size_t tsms, const char *tag, const char *data, size_t size) const = 0;
    };
  }
}