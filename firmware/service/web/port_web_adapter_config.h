#pragma once
#include <framework/broker.h>
#include <infrastructure/config/entry/wifi_settings.h>

namespace service
{
  namespace web
  {
    struct IPortAdapterConfig
    {
      virtual bool set_wifi_config_station(const infra::WifiSettingsN &settings) = 0;
      virtual bool set_wifi_config_softap(const infra::WifiSettingsN &settings) = 0;
      virtual bool load() const = 0;
      virtual bool save() const = 0;
    };
  }
}