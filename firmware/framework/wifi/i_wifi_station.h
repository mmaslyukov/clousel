#pragma once

namespace wifi
{

  struct WifiStationConfiguration
  {
    constexpr WifiStationConfiguration(const char *ssid, const char *password)
        : ssid(ssid), password(password) {}
    const char *ssid;
    const char *password;
  };

  struct IWifiStation
  {
    virtual bool enable() = 0;
    virtual bool disable() = 0;
    virtual bool is_enabled() const = 0;
    
    virtual bool connect() = 0;
    virtual bool disconnect() = 0;
    virtual bool is_connected() const = 0;
  };

} // namespace wifi
