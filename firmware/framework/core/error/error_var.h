#pragma once
#include "i_error.h"
#include <stdint.h>
#include <cstring>

namespace core
{
  namespace error
  {
    template <size_t N>
    struct ErrorVar : public IError
    {
      ErrorVar(const char *str) : _str(_buf)
      {
        strcpy_s(_buf, sizeof(_buf), str);
      }
      virtual IError& clear() override
      {
        _str = nullptr;
        return *this;
      }
      virtual const char *error() const override
      {
        return _str;
      }

      ErrorVar &format(const char *tag, const char *format, ...)
      {
        _str = _buf;
        va_list args;
        va_start(args, format);
        vsnprintf(_buf, sizeof(_buf), format, args);
        va_end(args);
        return *this;
      }

    private:
      const char* _str;
      char _buf[N];
    };

  } // namespace error
} // namespace core
