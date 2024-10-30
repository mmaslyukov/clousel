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

        bool is_station() const
        {
          return _station;
        }
        bool is_sofap() const
        {
          return _softap;
        }
        static inline EventWifiModeChanged to_station()
        {
          return EventWifiModeChanged(true, false);
        }
        static inline EventWifiModeChanged to_softap()
        {
          return EventWifiModeChanged(false, true);
        }
        static inline EventWifiModeChanged to_none()
        {
          return EventWifiModeChanged(false, false);
        }

      private:
        bool _station;
        bool _softap;
      };

    } // namespace event

  } // namespace mode

} // namespace service
