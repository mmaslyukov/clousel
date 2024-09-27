#pragma once
#include <stdint.h>
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
      // template<typename E> static EventBase event()  {
      //   return E.event();
      // }
      virtual EventId id() const
      {
        return _id;
      };
    // private:
      // Event() : EventBase(""), _id(0) {}
      Event(const char *name) : EventBase(name)
      {
        static EventId id = 0;
        _id = ++id;
      }

    private:
      EventId _id;
    };


    // template<typename T>
    // struct EventMaker {
    //   static inline T event()
    //   {
    //     return T();
    //   }
    //   static inline EventBase base()
    //   {
    //     return EventBase(T::event_name());
    //   }
    // };
  }
}
