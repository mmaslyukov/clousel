#pragma once

namespace core
{
  namespace error
  {
    struct IError
    {
      virtual IError& clear() = 0;
      virtual const char *error() const = 0;
    };

  } // namespace error
} // namespace core
