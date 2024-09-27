#pragma once

#include <stdarg.h>
#include <stdint.h>
#include <stdio.h>
#include "../i_timestamp.h"
#include "i_dumpable.h"
#include "i_logger_system.h"
#include "config.h"

namespace core
{
  namespace logger
  {
    struct DumpableNone : public IDumpable
    {
      virtual void dump(const char *tag, const uint8_t *data, size_t size) const {}
      virtual void dump_ascii(const char *tag, const uint8_t *data, size_t size) const {}
      virtual void enable() {}
      virtual void disable() {}
    };

    class Dumpable : public IDumpable
    {
    public:
      Dumpable(
          const Configuration &config,
          const ILoggerSystem &system,
          const ITimestamp &timestamp,
          bool tag = false,
          bool enabled = true)
          : _config(config), _sys(system), _timestamp(timestamp), _tag(tag), _enabled(enabled) {}
      virtual void dump(const char *tag, const uint8_t *data, size_t size) const override
      {
        _dump(tag, data, size, "%02X");
      }
      virtual void dump_ascii(const char *tag, const uint8_t *data, size_t size) const override
      {
        _dump(tag, data, size, "%c");
      }

      virtual void enable() override { _enabled = true; };
      virtual void disable() override { _enabled = false; }

    private:
      void _dump(const char *tag, const uint8_t *data, size_t size, const char *format) const
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

          shift += snprintf((char *)&_config.buffer[shift],
                            _config.size - shift, "%d|", static_cast<int>(size));

          size_t text_limit = _config.size - 4; // one for \0 and three for dots
          size_t actual_size = size > _config.size - shift ? _config.size - shift : size;
          for (size_t i = 0; (shift < _config.size) && (i < size); i++)
          {
            if (shift < text_limit)
            {
              shift += snprintf((char *)&_config.buffer[shift],
                                _config.size - shift, format, data[i]);
            }
            else
            {
              shift += snprintf((char *)&_config.buffer[shift],
                                _config.size - shift, ".");
            }
          }

          _sys.output(_config.verbosity, tag, _config.buffer, shift);
        }
      }

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