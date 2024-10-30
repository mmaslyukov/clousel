#pragma once
#include <framework/core/observer.h>
namespace core
{
  namespace io
  {
    namespace event
    {
      struct EventButtonClicked : public core::observer::Event
      {
        EventButtonClicked(uint32_t button_id)
            : core::observer::Event(event_name()),
              _button_id(button_id) {}
        static const char *event_name()
        {
          return "botton.pressed";
        }
        uint32_t button_id() const 
        {
          return _button_id;
        }
      private:
        uint32_t _button_id;
      };

    } // namespace event

  } // namespace io

} // namespace core
