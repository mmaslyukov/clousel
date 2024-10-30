#pragma once
namespace wifi
{
  struct WifiSoftApConfiguration
  {
    constexpr WifiSoftApConfiguration(const char *ssid, const char *password)
        : ssid(ssid), password(password) {}
    constexpr WifiSoftApConfiguration(const char *ssid)
        : ssid(ssid), password("") {}
    const char *ssid;
    const char *password;
  };

  struct IWifiSoftAp
  {
    virtual bool enable() = 0;
    virtual bool disable() = 0;
    virtual bool is_enabled() const = 0;
    // virtual bool available() const = 0;
  };
}