#pragma once
#include <framework/core/i_runnable.h>
#include <framework/core/logger.h>
#include <framework/core/observer.h>

#include <service/mode/event/event_wifi_mode_changed.h>

#include "port_web_controller_api.h"
#include "port_web_adapter_config.h"
#include "event/event_wifi_cred_saved.h"
// #include <platform/windows/httplib.h>

#include <cstring>

namespace service
{
  namespace web
  {
    class WebService
        : /* public core::IRunnable, */
          public IPortWebControllerApi,
          public core::observer::IListener,
          public core::observer::Publisher_1
    {
    public:
      WebService(
          const core::logger::ILogger &logger,
          IPortAdapterConfig &config)
          : _logger(logger),
            _config(config),
            _wifi_mode_changed(mode::event::EventWifiModeChanged::to_none()) {}
      // virtual void run() override
      // {
      // }

      virtual void notify(const core::observer::Event &event) override
      {
        if (event.name() == mode::event::EventWifiModeChanged::event_name())
        {
          _wifi_mode_changed = reinterpret_cast<const mode::event::EventWifiModeChanged&>(event);
          // const auto &mode_changed = reinterpret_cast<const mode::event::EventWifiModeChanged&>(&event);
          // _enabled = mode_changed->is_sofap();
        }
      }

      virtual bool submit(const infra::WifiSettingsN &settings) const override
      {
        bool result = false;
        do
        {
          // if (!_enabled)
          if (!_wifi_mode_changed.is_sofap())
          {
            _logger.wrn().log(TAG, "Can be handled only in SoftAp mode");
            break;
          }
          // if (!strlen(settings.pswd()) || !strlen(settings.ssid()))
          if (settings.ssid().empty())
          {
            _logger.err().log(TAG, "ssid is empty");
            break;
          }
          _logger.inf().log(TAG, "ssid:%s, pswd:%s", settings.ssid().data(), settings.pswd().data());
          bool saved = _config.set_wifi_config_station(settings);
          if (!saved)
          {
            _logger.err().log(TAG, "Failed to save new Wifi credentials");
            break;
          }
          if (!_config.save())
          {
            _logger.err().log(TAG, "Failed to save persistency to the flash");
            break;
          }
          publish(event::EventWifiCredSaved());
          result = true;

        } while (false);

        return result;
      }

    private:
      const core::logger::ILogger &_logger;
      IPortAdapterConfig &_config;
      mode::event::EventWifiModeChanged _wifi_mode_changed;
      // bool _enabled;
      static constexpr const char *TAG = "web";
    };
  }

}