#pragma once
#include "i_wifi_mode.h"

namespace wifi
{
  struct IWifiManager
  {
    // virtual bool enable() = 0;
    // virtual bool disable() = 0;
    virtual IWifiMode &mode() = 0;
  };

} // namespace wifi
