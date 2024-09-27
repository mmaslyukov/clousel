#pragma once

#include <stdarg.h>
#include <stdint.h>
#include <stdio.h>
#include "i_logger.h"

namespace core
{
  namespace logger
  {

    class Logger : public ILogger
    {
      Logger() = delete;

    public:
      Logger(
          const IPrintable &err,
          const IPrintable &wrn,
          const IPrintable &inf,
          const IPrintable &dbg,
          const IDumpable &raw) : _err(err), _wrn(wrn), _inf(inf), _dbg(dbg), _raw(raw) {}

      virtual const IPrintable &dbg() const override
      {
        return _dbg;
      }
      virtual const IPrintable &err() const override
      {
        return _err;
      }
      virtual const IPrintable &inf() const override
      {
        return _inf;
      }
      virtual const IPrintable &wrn() const override
      {
        return _wrn;
      }
      virtual const IDumpable &raw() const override
      {
        return _raw;
      }

    private:
      const IPrintable &_err;
      const IPrintable &_wrn;
      const IPrintable &_inf;
      const IPrintable &_dbg;
      const IDumpable &_raw;
    };
  }
}