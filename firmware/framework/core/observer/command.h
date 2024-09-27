#pragma once
#include <stdint.h>
#include "i_message.h"

namespace core
{
  namespace observer
  {
    using CommandId = size_t;
    struct Command : public IMessage
    {
      Command(const char *name) : _name(name)
      {
        static CommandId id = 0;
        _id = id;
      }
      virtual const char *name() const override
      {
        return _name;
      }

      virtual CommandId id() const
      {
        return _id;
      };

    private:
      const char *_name;
      CommandId _id;
    };
  }
}
