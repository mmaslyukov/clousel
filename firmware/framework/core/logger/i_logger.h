#pragma once
#include <stdint.h>
#include "verbosity.h"
#include "i_printable.h"
#include "i_dumpable.h"

namespace core
{
  namespace logger
  {
    
    struct ILogger
    {
      virtual const IPrintable &err() const = 0;
      virtual const IPrintable &wrn() const = 0;
      virtual const IPrintable &inf() const = 0;
      virtual const IPrintable &dbg() const = 0;
      virtual const IPrintable &vrb() const = 0;
      virtual const IDumpable &raw() const = 0;
    };

  }
}