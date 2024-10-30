#pragma once
#include <stdint.h>
#include <assert.h>
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
      constexpr Publisher() : _index_listener(0) {}
      bool add_subscriber(IListener *listener, const EventBase &event) override
      {
        assert(_index_listener < N);

        if (_index_listener < N)
        {
          _listeners[_index_listener].event = event;
          _listeners[_index_listener].listener = listener;
          _index_listener++;
          return true;
        }
        return false;
      }

    protected:
      void publish(const Event &event) const
      {
        for (auto i = 0; i < _index_listener; i++)
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
      size_t _index_listener;
    };
    using Publisher_1 = Publisher<1>;
    using Publisher_2 = Publisher<2>;
    using Publisher_3 = Publisher<3>;
  }
}
