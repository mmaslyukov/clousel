#pragma once
namespace wifi
{

  struct WifiSoftApConfiguration
  {
    constexpr WifiSoftApConfiguration(const char *password, const char *ssid)
        : password(password), ssid(ssid) {}
    const char *password;
    const char *ssid;
  };

  struct IWifiSoftAp
  {
    virtual bool enable() = 0;
    virtual void disable() = 0;
  };
}