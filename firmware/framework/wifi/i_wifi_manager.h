#pragma once
#include "i_wifi_softap.h"
#include "i_wifi_station.h"

namespace wifi
{
  struct IWifiManager
  {
    virtual IWifiManager &mode() = 0;
  };

  struct IWifiMode
  {
    /* switches to SoftAP (if not switched before) and return reference to the class */
    virtual IWifiSoftAp &soft_ap() = 0;
    /* switches to Station (if not switched before) and return reference to the class */
    virtual IWifiStation &station() = 0;
  };

} // namespace wifi
