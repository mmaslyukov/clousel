#pragma once

namespace wifi
{

  struct WifiStationConfiguration
  {
    constexpr WifiStationConfiguration(const char *password, const char *ssid)
        : password(password), ssid(ssid) {}
    const char *password;
    const char *ssid;
  };

  struct IWifiStation
  {
    virtual bool connect() = 0;
    virtual bool disconnect() = 0;
    virtual bool is_connected() = 0;
  };

} // namespace wifi
