#pragma once
// #include <infrastructure/config/wifi_settings.h>
namespace service
{
  namespace mode
  {
    struct IPortAdapterWifi
    {
      virtual bool swith_to_softap() = 0;
      virtual bool swith_to_station() = 0;
      virtual bool is_softap() const = 0;
      virtual bool is_station() const = 0;
      virtual bool is_station_connected() const = 0;
    };

  }
}