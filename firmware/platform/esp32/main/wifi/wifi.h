#pragma once
#include <framework/wifi.h>
#include <framework/core/logger.h>
#include <service/mode/mode_service.h>
// #include <service/mode/mode_service.h>

#include "../event/esp_event_provider.h"
#include "wifi_mode.h"
#include "wifi mode_station.h"
#include "wifi_mode_softap.h"

#include "esp_mac.h"
#include "esp_wifi.h"

class Wifi : public wifi::IWifiManager, public service::mode::IPortAdapterWifi, public EspEvetProvider<2>
{
public:
  Wifi(WifiMode &mode, const core::logger::ILogger &logger) : _mode(mode), _logger(logger)
  {
    ESP_ERROR_CHECK(esp_netif_init());

    esp_netif_create_default_wifi_sta();
    esp_netif_create_default_wifi_ap();

    wifi_init_config_t cfg = WIFI_INIT_CONFIG_DEFAULT();
    ESP_ERROR_CHECK(esp_wifi_init(&cfg));
    ESP_ERROR_CHECK(esp_event_handler_instance_register(WIFI_EVENT,
                                                        ESP_EVENT_ANY_ID,
                                                        &wifi_event_handler,
                                                        this,
                                                        &_instance_any_id));
    ESP_ERROR_CHECK(esp_event_handler_instance_register(IP_EVENT,
                                                        // ESP_EVENT_ANY_ID,
                                                        IP_EVENT_STA_GOT_IP,
                                                        &wifi_event_handler,
                                                        this,
                                                        &_instance_got_ip));
    ESP_ERROR_CHECK(esp_wifi_set_mode(WIFI_MODE_NULL));
    ESP_ERROR_CHECK(esp_wifi_set_storage(WIFI_STORAGE_RAM));
    add_handler(&mode.sta());
    add_handler(&mode.sap());
  }

  virtual wifi::IWifiMode &mode()
  {
    return _mode;
  }

  virtual bool swith_to_softap() override
  {
    bool res = true;
    res = res && _mode.station().disable();
    res = res && _mode.soft_ap().enable();
    return res;
  }

  virtual bool swith_to_station() override
  {
    bool res = true;
    res = res && _mode.soft_ap().disable();
    res = res && _mode.station().enable();
    return res;
  }

  virtual bool is_softap() const override
  {
    return _mode.soft_ap().is_enabled();
    return false;
  }

  virtual bool is_station() const override
  {
    return _mode.station().is_enabled();
  }

  virtual bool is_station_connected() const override
  {
    return _mode.station().is_connected();
  }

private:
  static void wifi_event_handler(void *arg, esp_event_base_t event_base, int32_t event_id, void *event_data)
  {
    Wifi *self = reinterpret_cast<Wifi *>(arg);
    bool res = self->handle(event_base, event_id, event_data);
    self->_logger.inf().log(TAG, "eb:%p, eid:%d, res:%d, wifi_event:%p, ip_event:%p", event_base, event_id, res, WIFI_EVENT, IP_EVENT);
  }


private:
  wifi::IWifiMode &_mode;
  const core::logger::ILogger &_logger;
  esp_event_handler_instance_t _instance_any_id;
  esp_event_handler_instance_t _instance_got_ip;

  static constexpr const char *TAG = "wifi.mode";
};
