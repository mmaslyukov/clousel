#pragma once
#include <service/mode/mode_service.h>
#include <framework/wifi.h>
#include <framework/core/logger.h>

#include "../event/i_esp_event_handler.h"

#include "esp_mac.h"
#include "esp_wifi.h"

class WifiModeSoftAp : public wifi::IWifiSoftAp, public IEspEvetHandler
{
public:
  WifiModeSoftAp(const wifi::WifiSoftApConfiguration &config, const core::logger::ILogger &logger) : _config(config), _logger(logger), _enabled(false)
  {
    assert(sizeof(_wifi_config.ap.ssid) > strlen(_config.ssid) + 1);
    assert(sizeof(_wifi_config.ap.password) > strlen(_config.password) + 1);

    memset((void *)&_wifi_config, 0, sizeof(_wifi_config));
    _wifi_config.ap.ssid_len = strlen(_config.ssid);
    strcpy(reinterpret_cast<char *>(_wifi_config.ap.ssid), _config.ssid);
    strcpy(reinterpret_cast<char *>(_wifi_config.ap.password), _config.password);

    _wifi_config.ap.authmode = WIFI_AUTH_WPA2_PSK;
    _wifi_config.ap.sae_pwe_h2e = WPA3_SAE_PWE_BOTH;
    _wifi_config.ap.pmf_cfg.required = true;
    _wifi_config.ap.channel = 1;
    _wifi_config.ap.max_connection = 4;

    if (strlen(_config.password) == 0)
    {
      _wifi_config.ap.authmode = WIFI_AUTH_OPEN;
    }
  }

  virtual bool enable()
  {
    bool res = true;
    if (!_enabled)
    {
      _logger.inf().log(TAG, "Enabling SoftAp Mode");
      res = res && esp_wifi_set_mode(WIFI_MODE_AP) == ESP_OK;
      res = res && esp_wifi_set_config(WIFI_IF_AP, &_wifi_config) == ESP_OK;
      res = res && esp_wifi_start() == ESP_OK;
      if (res)
      {
        _enabled = true;
      }
    }
    return res;
  }

  virtual bool disable()
  {
    bool res = true;
    if (_enabled)
    {
      _logger.inf().log(TAG, "Disabling Station Mode");
      _enabled = false;
      res = res && esp_wifi_stop() == ESP_OK;
      res = res && esp_wifi_set_mode(WIFI_MODE_NULL) == ESP_OK;
    }
    return res;
  }

  virtual bool is_enabled() const
  {
    return _enabled;
  }
  static bool dump_mac(char *str, size_t cap)
  {
    uint8_t mac[6];
    bool res = true;
    res = res && esp_read_mac(mac, ESP_MAC_WIFI_SOFTAP) == ESP_OK;
    res = res && cap > 17;
    if (res)
    {
      snprintf(str, cap, "%02X-%02X-%02X-%02X-%02X-%02X", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5]);
    }
    return res;
  }
private:
  virtual bool handle(esp_event_base_t event_base, int32_t event_id, void *event_data) override
  {
    bool res = true;
    if (_enabled && event_base == WIFI_EVENT)
    {
      switch (event_id)
      {
      case WIFI_EVENT_AP_STACONNECTED:
      {
        wifi_event_ap_staconnected_t *event = reinterpret_cast<wifi_event_ap_staconnected_t *>(event_data);
        _logger.inf().log(TAG, "station " MACSTR " join, AID=%d", MAC2STR(event->mac), event->aid);
      }
      break;
      case WIFI_EVENT_AP_STADISCONNECTED:
      {
        wifi_event_ap_stadisconnected_t *event = reinterpret_cast<wifi_event_ap_stadisconnected_t *>(event_data);
        _logger.inf().log(TAG, "station " MACSTR " leave, AID=%d", MAC2STR(event->mac), event->aid);
      }
      break;
      default:
        res = false;
        break;
      }
    }
    else
    {
      res = false;
    }
    return res;
  }



private:
  const wifi::WifiSoftApConfiguration &_config;
  const core::logger::ILogger &_logger;
  bool _enabled;
  esp_netif_t *_netif;
  wifi_config_t _wifi_config;
  static constexpr const char *TAG = "wifi.sap";
};
