#pragma once
#include <framework/core/observer.h>
namespace service
{
  namespace mode
  {
    namespace event
    {
      struct EventWifiModeChanged : public core::observer::Event
      {
        EventWifiModeChanged(bool station, bool softap)
            : core::observer::Event(event_name()),
              _station(station),
              _softap(softap) {}
        static const char *event_name()
        {
          return "wifi.mode.changed";
        }

        bool station() const
        {
          return _station;
        }
        bool sofap() const
        {
          return _softap;
        }

      private:
        bool _station;
        bool _softap;
      };

    } // namespace event

  } // namespace mode

} // namespace service
