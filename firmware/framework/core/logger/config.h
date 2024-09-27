#pragma once
#include <stdint.h>
#include "verbosity.h"
namespace core
{
  namespace logger
  {
    struct Configuration
    {
      constexpr Configuration(char *buffer, const size_t size, const Verbosity &verbosity)
          : buffer(buffer), size(size), verbosity(verbosity) {}
      char *buffer;
      const size_t size;
      const Verbosity &verbosity;
    };
  }
}
