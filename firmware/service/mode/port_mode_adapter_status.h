#pragma once
namespace service
{
  namespace mode
  {
    struct IPortAdapterStatus
    {
      virtual void led_wifi_station(bool lit) = 0;
      virtual void led_wifi_softap(bool lit) = 0;
      virtual void led_wifi_connected(bool lit) = 0;
    };
  } // namespace mode
}
