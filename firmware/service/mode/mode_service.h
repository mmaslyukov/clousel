#pragma once
#include <framework/core/observer.h>
#include <framework/core/logger.h>
#include <framework/core/i_runnable.h>

#include <service/web/event/event_wifi_cred_saved.h>

#include "port_mode_controller_button.h"
#include "port_mode_adapter_wifi.h"
#include "port_mode_adapter_status.h"
#include "port_mode_adapter_broker.h"

#include "event/event_wifi_mode_changed.h"

namespace service
{
  namespace mode
  {
    class ModeService
        // subscribe on web service event to switch back to wifi station mode
        : public core::IRunnable,
          public IPortButtonController,
          public core::observer::IListener,
          public core::observer::Publisher_1
    {
    public:
      constexpr ModeService(
          const core::logger::ILogger &logger,
          IPortAdapterWifi &wifi,
          IPortAdapterStatus &status,
          IPortAdapterBroker &broker)
          // IPortAdapterConfig &config)
          : _logger(logger), _wifi(wifi), _status(status), _broker(broker) /*,  _config(config) */
      {
      }

      virtual void notify(const core::observer::Event &event) override
      {
        if (event.name() == web::event::EventWifiCredSaved::event_name())
        {
          _wifi.swith_to_station();
        }
      }

      virtual void run() override
      {
        if (_wifi.is_softap())
        {
          _status.led_wifi_softap(true);
          _status.led_wifi_station(false);
          if (_broker.is_connected())
          {
            _broker.disconnect();
          }
          publish(event::EventWifiModeChanged(false, true));
        }
        else if (_wifi.is_station())
        {
          _status.led_wifi_softap(false);
          _status.led_wifi_station(true);
          _status.led_wifi_connected(_wifi.is_station_connected());
          if (!_broker.is_connected())
          {
            _broker.connect();
          }
          publish(event::EventWifiModeChanged(true, false));
        }
      }

      // private:
      virtual void clicked() override
      {
        // do nothing
      }

      virtual void pressed() override
      {
        if (!_wifi.is_softap())
        {
          _logger.inf().log(TAG, "Button has been pressed, and start wifi swithing to SoftAp mode");
          _broker.disconnect();
          _wifi.swith_to_softap();
        }
        else
        {
          _logger.inf().log(TAG, "Button has been pressed, but wifi is in SoftAp mode");
        }
      }

    private:
      static constexpr const char *TAG = "modsvc";
      const core::logger::ILogger &_logger;
      IPortAdapterWifi &_wifi;
      IPortAdapterStatus &_status;
      IPortAdapterBroker &_broker;
      // IPortAdapterConfig &_config;
    };
  }
}