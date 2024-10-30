#pragma once
#include <service/mode/mode_service.h>
// #include <service/web/event/event_wifi_cred_saved.h>
#include <framework/wifi.h>
#include <framework/core/logger.h>
// #include <framework/core/observer.h>

#include "../event/i_esp_event_handler.h"

#include <string.h>

#include "esp_mac.h"
#include "esp_wifi.h"
#include "esp_system.h"
#include "esp_event.h"
#include "esp_log.h"

// #include "freertos/FreeRTOS.h"
// #include "freertos/task.h"
// #include "freertos/event_groups.h"
// #include "esp_system.h"
// #include "esp_wifi.h"
// #include "esp_event.h"
// #include "esp_log.h"
// #include "nvs_flash.h"

// #include "lwip/err.h"
// #include "lwip/sys.h"
class WifiModeStation : public wifi::IWifiStation, public IEspEvetHandler //, public core::observer::IListener
{
public:
  constexpr WifiModeStation(const wifi::WifiStationConfiguration &config, const core::logger::ILogger &logger) : _config(config), _logger(logger), _connected(false), _enabled(false)
  {
    assert(sizeof(_wifi_config.sta.ssid) > strlen(_config.ssid) + 1);
    assert(sizeof(_wifi_config.sta.password) > strlen(_config.password) + 1);
    memset((void *)&_wifi_config, 0, sizeof(_wifi_config));
    // strcpy(reinterpret_cast<char *>(_wifi_config.sta.ssid), _config.ssid);
    // strcpy(reinterpret_cast<char *>(_wifi_config.sta.password), _config.password);
    _wifi_config.sta.threshold.authmode = WIFI_AUTH_WPA2_PSK;
    _wifi_config.sta.sae_pwe_h2e = WPA3_SAE_PWE_BOTH;
    _wifi_config.sta.sae_h2e_identifier[0] = 0;
  }

  virtual bool enable()
  {

    bool res = true;
    if (!_enabled)
    {
      strcpy(reinterpret_cast<char *>(_wifi_config.sta.ssid), _config.ssid);
      strcpy(reinterpret_cast<char *>(_wifi_config.sta.password), _config.password);
      _logger.inf().log(TAG, "Enabling Station Mode ssid:'%s', pswd:'%s'", _wifi_config.sta.ssid, _wifi_config.sta.password);

      res = res && esp_wifi_set_mode(WIFI_MODE_STA) == ESP_OK;
      res = res && esp_wifi_set_config(WIFI_IF_STA, &_wifi_config) == ESP_OK;
      res = res && esp_wifi_start() == ESP_OK;
      if (res)
      {
        _enabled = true;
        _logger.inf().log(TAG, "Station Mode has been enabled");
      }
      else
      {
        _logger.err().log(TAG, "Fail to enable station mode");
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
      res = res && disconnect();
      res = res && esp_wifi_stop() == ESP_OK;
      res = res && esp_wifi_set_mode(WIFI_MODE_NULL) == ESP_OK;
      if (res)
      {
        _enabled = false;
        _logger.inf().log(TAG, "Station Mode has been disabled");
      }
      else
      {
        _logger.err().log(TAG, "Fail to disable station mode");
      }
    }
    return res;
  }
  virtual bool is_enabled() const
  {
    return _enabled;
  }

  virtual bool connect()
  {
    if (_enabled)
    {
      _logger.inf().log(TAG, "Connecting to %s", _config.ssid);
      esp_err_t err = esp_wifi_connect();
      if (err != ESP_OK)
      {
        _logger.err().log(TAG, "Connection to wifi is failed with error:%d", err);
      }
      return err == ESP_OK;
    }
    return false;
  }
  virtual bool disconnect()
  {
    if (_enabled)
    {
      _logger.inf().log(TAG, "Disconnecting from %s", _config.ssid);
      esp_err_t err = esp_wifi_disconnect();
      if (err != ESP_OK)
      {
        _logger.err().log(TAG, "Disconnection from wifi is failed with error:%d", err);
      }
      _connected = false;
      return err == ESP_OK;
    }
    return false;
  }
  virtual bool is_connected() const
  {
    return _connected;
  }

private:
  virtual bool handle(esp_event_base_t event_base, int32_t event_id, void *event_data) override
  {
    bool res = true;
    if (_enabled && event_base == WIFI_EVENT)
    {
      res = handle_wifi_event(event_id, event_data);
    }
    else if (_enabled && event_base == IP_EVENT)
    {
      res = handle_ip_event(event_id, event_data);
    }
    else
    {
      res = false;
    }
    return res;
  }

  bool handle_wifi_event(int32_t event_id, void *event_data)
  {
    bool res = true;
    switch (event_id)
    {
    case WIFI_EVENT_STA_START:
      connect();
      break;
    case WIFI_EVENT_STA_CONNECTED:
      // uint8_t mac[6];
      // esp_base_mac_addr_get(mac);
      // _logger.inf().log(TAG, "Connected to %s, my MAC:" MACSTR, _config.ssid, MAC2STR(mac));
      _logger.inf().log(TAG, "Connected to %s", _config.ssid);
      break;
    case WIFI_EVENT_STA_DISCONNECTED:
      _logger.inf().log(TAG, "Disconnected from %s", _config.ssid);
      if (_connected)
      {
        _connected = false;
      }
      connect();
      break;
    default:
      res = false;
      break;
    }
    return res;
  }
  bool handle_ip_event(int32_t event_id, void *event_data)
  {
    bool res = true;
    switch (event_id)
    {
    case IP_EVENT_STA_GOT_IP:
      _logger.inf().log(TAG, "got ip:" IPSTR, IP2STR(&reinterpret_cast<ip_event_got_ip_t *>(event_data)->ip_info.ip));
      _connected = true;
      break;
    default:
      res = false;
      break;
    }
    return res;
  }
  //   virtual void notify(const core::observer::Event &event) override
  // {
  //   if (event.name() == service::web::event::EventWifiCredSaved::event_name())
  //   {
  //     _config.password
  //   }

  // }

private:
  const wifi::WifiStationConfiguration &_config;
  const core::logger::ILogger &_logger;
  bool _connected;
  bool _enabled;
  esp_netif_t *_netif;
  wifi_config_t _wifi_config;
  static constexpr const char *TAG = "wifi.sta";
  // EventGroupHandle_t _wifi_event_group;
  // static constexpr const uint32_t WIFI_CONNECTED_BIT = BIT0;
  // static constexpr const uint32_t WIFI_FAIL_BIT = BIT1;
};