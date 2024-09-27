#pragma once
#include <framework/core/observer.h>
namespace service
{
  namespace web
  {
    namespace event
    {
      struct EventWifiCredSaved : public core::observer::Event {
        EventWifiCredSaved() : core::observer::Event(event_name()) {}
        static const char* event_name() 
        {
          return "wifi.cred.saved";
        }
      };
      
    } // namespace event
    
  } // namespace mode
  
  
} // namespace service
