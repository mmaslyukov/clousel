#pragma once

#include <stdarg.h>
#include <stdint.h>
#include <stdio.h>
#include "../i_timestamp.h"
#include "verbosity.h"
#include "config.h"

namespace core
{
  namespace logger
  {
    struct PrintableNone : public IPrintable
    {
      virtual void log(const char *tag, const char *format, ...) const {}
      virtual void enable() {}
      virtual void disable() {}
    };

    class Printable : public IPrintable
    {
    public:
      Printable(
          const Configuration &config,
          const ILoggerSystem &system,
          const ITimestamp &timestamp,
          bool tag = false,
          bool enabled = true)
          : _config(config), _sys(system), _timestamp(timestamp), _tag(tag), _enabled(enabled) {}
      virtual void log(const char *tag, const char *format, ...) const
      {
        if (_enabled)
        {
          size_t shift = 0;
          const char *verbosity = _config.verbosity.name();
          if (verbosity)
          {
            shift += snprintf((char *)&_config.buffer[shift],
                              _config.size - shift, "%s ", verbosity);
          }

          size_t ts = _timestamp.get();
          if (ts > 0)
          {
            shift += snprintf((char *)&_config.buffer[shift],
                              _config.size - shift, "%zu ", ts);
          }

          if (_tag)
          {
            shift += snprintf((char *)&_config.buffer[shift],
                              _config.size - shift, "<%s> ", tag);
          }

          va_list args;
          va_start(args, format);
          shift += vsnprintf((char *)&_config.buffer[shift],
                             _config.size - shift, format, args);
          va_end(args);
          _sys.output(_config.verbosity, tag, _config.buffer, shift);
        }
      }
      virtual void enable() { _enabled = true; };
      virtual void disable() { _enabled = false; }

    private:
      const Configuration &_config;
      const ILoggerSystem &_sys;
      const ITimestamp &_timestamp;
      // const Verbosity &_verbosity;
      bool _tag;
      bool _enabled;
    };

  }
}