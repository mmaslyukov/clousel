#pragma once
#include <stdint.h>
#include "verbosity.h"
#include "i_enabable.h"

namespace core
{
  namespace logger
  {
    
    struct IPrintable : public IEnabable
    {
      virtual void log(const char *tag, const char *format, ...) const = 0;
    };

  }
}