#pragma once
#include <stdint.h>
#include <cstddef>
#include "i_message.h"

namespace core
{
  namespace observer
  {

    using EventId = size_t;

    struct EventBase : public IMessage
    {
      EventBase(): _name("") {}
      EventBase(const char* name): _name(name) {}
      virtual const char *name() const override
      {
        return _name;
      }

    private:
      const char *_name;
    };

    struct Event : public EventBase
    {
      virtual EventId id() const
      {
        return _id;
      };
      Event(const char *name) : EventBase(name)
      {
        static EventId id = 0;
        _id = ++id;
      }

    private:
      EventId _id;
    };
  }
}
