#pragma once
#include "i_error.h"

namespace core
{
  namespace error
  {
    struct ErrorConst : public IError
    {
      constexpr ErrorConst() : _str(nullptr) {}
      constexpr ErrorConst(const char *str) : _str(str) {}
      
      virtual IError& clear() override
      {
        _str = nullptr;
        return *this;
      }

      virtual const char *error() const override
      {
        return _str;
      }

    private:
      const char *_str;
    };

  } // namespace error
} // namespace core
