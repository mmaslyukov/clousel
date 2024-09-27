#pragma once
#include <framework/broker.h>
#include <infrastructure/config/entry/wifi_settings.h>

namespace service
{
  namespace mode
  {
    struct IPortAdapterConfig
    {
      virtual const infra::WifiSettingsN* wifi_config_station() = 0;
      virtual const infra::WifiSettingsN* wifi_config_softap() = 0;
    };

  }
}