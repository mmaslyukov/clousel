#pragma once
#include <stdint.h>
#include "i_subscribable.h"

namespace core
{
  namespace observer
  {
    template <size_t N>
    class Publisher : public ISubscribable
    {
      struct EventListenerPair
      {
        EventBase event;
        IListener *listener;
      };

    public:
      constexpr Publisher() : _index(0) {}
      bool add_subscriber(IListener *listener, const EventBase &event) override
      {
        if (_index < N)
        {
          _listeners[_index].event = event;
          _listeners[_index].listener = listener;
          _index++;
          return true;
        }
        return false;
      }

      void publish(const Event &event) const
      {
        for (auto i = 0; i < _index; i++)
        {
          // Note, here it compares pointer addresses rather than actual text
          if (_listeners[i].event.name() == event.name())
          {
            _listeners[i].listener->notify(event);
          }
        }
      }

    private:
      EventListenerPair _listeners[N];
      size_t _index;
    };
    using Publisher_1 = Publisher<1>;
    using Publisher_2 = Publisher<2>;
    using Publisher_3 = Publisher<3>;
  }
}
