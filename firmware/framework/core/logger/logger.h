#pragma once

#include <stdarg.h>
#include <stdint.h>
#include <stdio.h>
#include "i_logger.h"
#include "printable.h"
#include "dumpable.h"

namespace core
{
  namespace logger
  {

    class Logger : public ILogger
    {
      Logger() = delete;

    public:
      Logger(
          const IPrintable &err, // = PrintableNone(),
          const IPrintable &wrn, // = PrintableNone(),
          const IPrintable &inf, // = PrintableNone(),
          const IPrintable &dbg, // = PrintableNone(),
          const IPrintable &vrb, // = PrintableNone(),
          const IDumpable &raw/*  = DumpableNone() */) : _err(err), _wrn(wrn), _inf(inf), _dbg(dbg), _vrb(vrb), _raw(raw) {}

      virtual const IPrintable &err() const override
      {
        return _err;
      }
      virtual const IPrintable &wrn() const override
      {
        return _wrn;
      }
      virtual const IPrintable &inf() const override
      {
        return _inf;
      }
      virtual const IPrintable &dbg() const override
      {
        return _dbg;
      }
      virtual const IPrintable &vrb() const override
      {
        return _vrb;
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
      const IPrintable &_vrb;
      const IDumpable &_raw;
    };
  }
}