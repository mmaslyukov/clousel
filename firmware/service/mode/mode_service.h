#pragma once
#include <framework/core/observer.h>
#include <framework/core/logger.h>
#include <framework/core/i_runnable.h>

#include <service/web/event/event_wifi_cred_saved.h>

#include "port_mode_controller_button.h"
#include "port_mode_adapter_wifi.h"
#include "port_mode_adapter_status.h"
// #include "port_mode_adapter_broker.h"

#include "event/event_wifi_mode_changed.h"
// #include "event/event_wifi_connected.h"

namespace service
{
  namespace mode
  {
    class ModeService
        : public core::IRunnable,
          public IPortButtonController,
          public core::observer::IListener,
          public core::observer::Publisher_2
    {
    public:
      ModeService(
          const core::logger::ILogger &logger,
          IPortAdapterWifi &wifi,
          IPortAdapterStatus &status)
          : _logger(logger), _wifi(wifi), _status(status), _wifi_mode_to_change(event::EventWifiModeChanged::to_none())
      {
      }

      void switch_to_softap()
      {
        _wifi_mode_to_change = event::EventWifiModeChanged::to_softap();
      }
      void switch_to_station()
      {
        _wifi_mode_to_change = event::EventWifiModeChanged::to_station();
      }

      virtual void notify(const core::observer::Event &event) override
      {
        if (event.name() == web::event::EventWifiCredSaved::event_name())
        {
          switch_to_station();
        }
        else if (event.name() == core::io::event::EventButtonClicked::event_name())
        {
          const core::io::event::EventButtonClicked &ebc = reinterpret_cast<const core::io::event::EventButtonClicked &>(event);
          _logger.inf().log(TAG, "Has got a new ButtonClicked event btn_id:%d, id:%d", ebc.button_id(), ebc.id());
        }
        else if (event.name() == core::io::event::EventButtonHeld::event_name())
        {
          const core::io::event::EventButtonHeld &ebh = reinterpret_cast<const core::io::event::EventButtonHeld &>(event);
          _logger.inf().log(TAG, "Has got a new ButtonHeld event btn_id:%d, id:%d", ebh.button_id(), ebh.id());
          switch_to_softap();
        }
      }

      virtual void run() override
      {
        if (!_wifi.is_station() && _wifi_mode_to_change.is_station())
        {
          _wifi.swith_to_station();
          publish(_wifi_mode_to_change);
        }

        if (!_wifi.is_softap() && _wifi_mode_to_change.is_sofap())
        {
          _wifi.swith_to_softap();
          publish(_wifi_mode_to_change);
        }

        if (_wifi.is_softap())
        {
          _status.led_wifi_softap(true);
          _status.led_wifi_station(false);
        }
        else if (_wifi.is_station())
        {
          _status.led_wifi_softap(false);
          _status.led_wifi_station(true);
          _status.led_wifi_connected(_wifi.is_station_connected());
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
          _wifi_mode_to_change = event::EventWifiModeChanged::to_softap();
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
      event::EventWifiModeChanged _wifi_mode_to_change;
    };
  }
}