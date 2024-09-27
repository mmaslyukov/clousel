#pragma once
#include <framework/core/i_runnable.h>
#include <service/mode/mode_service.h>
#include <service/coin/coin_service.h>

namespace infra
{
  class Status
      : public service::mode::IPortAdapterStatus,
        public service::coin::IPortAdapterStatus,
        public core::IRunnable
  {
  public:
    Status(
        const core::logger::ILogger &logger,
        const core::ITimestamp &ts,
        core::io::IActuator<bool> &led_coin,
        core::io::IActuator<bool> &led_wifi_softap,
        core::io::IActuator<bool> &led_wifi_station,
        core::io::IActuator<bool> &led_wifi_station_connected)
        : _logger(logger),
          _ts(ts),
          _led_coin(led_coin),
          _led_wifi_softap(led_wifi_softap),
          _led_wifi_station(led_wifi_station),
          _led_wifi_station_connected(led_wifi_station_connected),
          _led_coin_ts(ts.get())
    {
    }

    virtual void led_wifi_station(bool lit) override
    {
      // _logger.dbg().log(TAG, "led wifi-station %s", lit ? "on" : "off");
      _led_wifi_station.set(lit);
    }

    virtual void led_wifi_softap(bool lit) override
    {
      // _logger.dbg().log(TAG, "led wifi-softap %s", lit ? "on" : "off");
      _led_wifi_softap.set(lit);
    }

    virtual void led_wifi_connected(bool lit) override
    {
      // _logger.dbg().log(TAG, "led wifi-connected %s", lit ? "on" : "off");
      _led_wifi_station_connected.set(lit);
    }

    virtual void led_coin_blink() override
    {
      // _logger.dbg().log(TAG, "led coin-blink");
      _led_coin_ts = _ts.get() + LED_COIN_BLINK_MS;
    }

    virtual void run() override
    {
      if (_led_coin_ts > _ts.get())
      {
        _led_coin.set(true);
      }
      else
      {
        _led_coin.set(false);
      }
    }

  private:
    const core::logger::ILogger &_logger;
    const core::ITimestamp &_ts;
    core::io::IActuator<bool> &_led_coin;
    core::io::IActuator<bool> &_led_wifi_softap;
    core::io::IActuator<bool> &_led_wifi_station;
    core::io::IActuator<bool> &_led_wifi_station_connected;
    size_t _led_coin_ts;
    static constexpr const char *TAG = "status";
    static constexpr size_t LED_COIN_BLINK_MS = 1500;
  };

} // namespace infra
