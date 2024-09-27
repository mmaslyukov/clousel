#pragma once
#include <stdint.h>
namespace core
{
  namespace logger
  {
    class Verbosity
    {
    public:
      constexpr Verbosity(const char *verbosty_name)
          : _id(0), _name(verbosty_name) {}
      constexpr Verbosity(const uint32_t verbosity_id)
          : _id(verbosity_id), _name(nullptr) {}
      virtual uint32_t id() const { return _id; }
      virtual const char *name() const { return _name; }

    private:
      const uint32_t _id;
      const char *_name;
    };
  }
}