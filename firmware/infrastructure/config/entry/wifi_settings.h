#pragma once
#include <stdint.h>
#include <framework/util/util.h>

namespace infra
{
  template <size_t N>
  struct WifiSettings
  {
    WifiSettings()
    {
      _ssid.clear();
      _pswd.clear();
    }

    WifiSettings(const char *ssid)
    {
      _ssid.replace(ssid);
      _pswd.clear();
    }

    WifiSettings(const char *ssid, const char *pswd)
    {
      _ssid.replace(ssid);
      _pswd.replace(pswd);
    }

    bool is_ssid_vaid() const
    {
      return !_ssid.empty();
      // return strlen(_ssid) > 0;
    }

    const CharContainer<N> &ssid() const
    {
      return _ssid;
    }

    const CharContainer<N> &pswd() const
    {
      return _pswd;
    }

  private:
    CharContainer<N> _ssid;
    CharContainer<N> _pswd;
    // char _ssid[N];
    // char _pswd[N];

    /* data */
  };
  using WifiSettingsN = WifiSettings<32>;

} // namespace infra
