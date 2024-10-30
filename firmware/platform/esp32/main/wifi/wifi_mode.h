#pragma once
#include <service/mode/mode_service.h>
#include <framework/wifi.h>
#include "wifi_mode_softap.h"
#include "wifi mode_station.h"

class WifiMode : public wifi::IWifiMode
{
public:
  WifiMode(WifiModeSoftAp &sap, WifiModeStation &sta) : _sap(sap), _sta(sta)
  {
  }

  virtual wifi::IWifiSoftAp &soft_ap()
  {
    return _sap;
  }
  virtual wifi::IWifiStation &station()
  {
    return _sta;
  }
  WifiModeSoftAp &sap()
  {
    return _sap;
  }
  WifiModeStation &sta()
  {
    return _sta;
  }

private:
  WifiModeSoftAp &_sap;
  WifiModeStation &_sta;
};