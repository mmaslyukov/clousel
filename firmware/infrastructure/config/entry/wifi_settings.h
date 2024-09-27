#pragma once
#include <stdint.h>

namespace infra
{
  template <size_t N>
  struct WifiSettings
  {
    WifiSettings()
    {
      memset(_ssid, 0, sizeof(_ssid));
      memset(_pswd, 0, sizeof(_pswd));
    }
    WifiSettings(const char *ssid, const char *pswd)
    {
      if (strlen(ssid) < N)
      {
        strcpy_s(_ssid, N, ssid);
      }
      else
      {
        memset(_ssid, 0, sizeof(_ssid));
      }

      if (strlen(pswd) < N)
      {
        strcpy_s(_pswd, N, pswd);
      }
      else
      {
        memset(_pswd, 0, sizeof(_pswd));
      }
    }
    const char *ssid() const
    {
      return _ssid;
    }
    const char *pswd() const
    {
      return _pswd;
    }

  private:
    char _ssid[N];
    char _pswd[N];

    /* data */
  };
  using WifiSettingsN = WifiSettings<16>;

} // namespace infra
