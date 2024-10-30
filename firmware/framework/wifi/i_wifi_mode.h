#pragma once
#include "i_wifi_softap.h"
#include "i_wifi_station.h"

namespace wifi
{
  struct IWifiMode
  {
    virtual IWifiSoftAp &soft_ap() = 0;
    virtual IWifiStation &station() = 0;
  };

} // namespace wifi
