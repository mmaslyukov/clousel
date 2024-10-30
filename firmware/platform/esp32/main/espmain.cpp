#include <stdio.h>
#include <sys/time.h>

#include <framework/core/logger.h>
#include <framework/core/io.h>
#include <infrastructure/config/config.h>
#include <infrastructure/status/status.h>

#include <service/mode/mode_service.h>
#include <service/coin/coin_service.h>
#include <service/web/web_service.h>

#include "wifi/wifi.h"
#include "server/router.h"
#include "mqtt/mqtt.h"

#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "driver/gpio.h"
#include "esp_log.h"
#include "sdkconfig.h"
#include "nvs_flash.h"

class LoggerSystem : public core::logger::ILoggerSystem, public core::ITimestamp
{
public:
  virtual size_t get() const override { return 0; };
  virtual void output(const core::logger::Verbosity &verbosity, size_t tsms, const char *tag, const char *data, size_t size) const override
  {
    ESP_LOG_LEVEL_LOCAL(verbosity.id(), tag, "%.*s", size, data);
  };
};

class Timestamp : public core::ITimestamp
{
public:
  virtual size_t get() const override
  {
    return esp_log_timestamp();
  }
};

struct Flash : public persistency::IPersistencyFlash
{
  Flash(const core::logger::ILogger &logger)
      : _logger(logger)
  {
    esp_err_t ret = nvs_flash_init();
    if (ret == ESP_ERR_NVS_NO_FREE_PAGES || ret == ESP_ERR_NVS_NEW_VERSION_FOUND)
    {
      ESP_ERROR_CHECK(nvs_flash_erase());
      ret = nvs_flash_init();
    }
    ESP_ERROR_CHECK(ret);
  }
  virtual bool load(uint8_t *memory, size_t size) const
  {
    bool res = true;
    bool to_close = false;
    nvs_handle_t my_handle;
    do
    {
      res = res && nvs_open(NAMESPACE, NVS_READWRITE, &my_handle) == ESP_OK;
      if (!res)
      {
        _logger.err().log(TAG, "Fail to open nvs");
        break;
      }
      to_close = true;
      size_t actual_size = size;
      res = res && nvs_get_blob(my_handle, NAMESPACE, memory, &actual_size) == ESP_OK;
      if (!res)
      {
        _logger.err().log(TAG, "Fail to get data from nvs");
        break;
      }
      res = res && actual_size == size;
      if (!res)
      {
        res = false;
        _logger.wrn().log(TAG, "Storage sizes are mismatch, probably fresh nvs");
        break;
      }
      _logger.dbg().log(TAG, "Loaded successfully");
    } while (false);

    if (to_close)
    {
      nvs_close(my_handle);
    }
    return res;
  };
  virtual bool save(const uint8_t *memory, size_t size) const
  {
    bool res = true;
    bool to_close = false;
    nvs_handle_t my_handle;
    do
    {
      res = res && nvs_open(NAMESPACE, NVS_READWRITE, &my_handle) == ESP_OK;
      if (!res)
      {
        _logger.err().log(TAG, "Fail to open nvs");
        break;
      }
      to_close = true;
      res = res && nvs_set_blob(my_handle, NAMESPACE, memory, size) == ESP_OK;
      if (!res)
      {
        _logger.err().log(TAG, "Fail to set data to nvs");
        break;
      }
      res = res && nvs_commit(my_handle) == ESP_OK;
      if (!res)
      {
        _logger.err().log(TAG, "Fail to commit data to nvs");
        break;
      }
      _logger.dbg().log(TAG, "Saved successfully");
    } while (false);

    if (to_close)
    {
      nvs_close(my_handle);
    }
    return res;
  };

private:
  const core::logger::ILogger &_logger;
  static constexpr const char *TAG = "flash";
  static constexpr const char *NAMESPACE = "persistency";
};

class Actuator : public core::io::IActuator<bool>
{
public:
  Actuator(uint32_t gpio, bool def) : _id(static_cast<gpio_num_t>(gpio)), _value(def)
  {
    init(_id, _value);
  }
  Actuator(gpio_num_t gpio, bool def) : _id(gpio), _value(def)
  {
    init(_id, _value);
  }
  virtual uint32_t id() const override
  {
    return static_cast<uint32_t>(_id);
  }
  virtual void set(const bool &value) override
  {
    if (_value != value)
    {
      _value = value;
      gpio_set_level(_id, _value);
    }
  }
  virtual bool get() const override
  {
    return _value;
  }

private:
  static void init(gpio_num_t gpio, bool value)
  {
    gpio_reset_pin(gpio);
    gpio_set_direction(gpio, GPIO_MODE_OUTPUT);
    gpio_set_level(gpio, value);
  }

private:
  gpio_num_t _id;
  bool _value;
};

class ButtonInput : public core::io::ISensor<bool>
{
public:
  ButtonInput(gpio_num_t gpio) : _id(gpio)
  {
    init(_id);
  }
  ButtonInput(uint32_t gpio) : _id(static_cast<gpio_num_t>(gpio))
  {
    init(_id);
  }

  virtual uint32_t id() const override
  {
    return static_cast<uint32_t>(_id);
  }
  virtual bool get() const override
  {
    // inverting signal
    return !gpio_get_level(_id);
  }

private:
  static void init(gpio_num_t gpio)
  {
    gpio_reset_pin(gpio);
    gpio_set_direction(gpio, GPIO_MODE_INPUT);
  }

private:
  gpio_num_t _id;
};

class MyGpio : public core::io::IActuator<bool>
{
public:
  MyGpio(const char *name, uint32_t id, const core::logger::ILogger &logger)
      : _name(name), _value(false), _id(id), _logger(logger)
  {
    log();
  }
  virtual uint32_t id() const override
  {
    return _id;
  }
  virtual void set(const bool &value) override
  {
    if (_value != value)
    {
      _value = value;
      log();
    }
  }
  virtual bool get() const override
  {
    return _value;
  }

private:
  void log()
  {
    _logger.dbg().log("gpio", "%c GPIO%d[%s]", _value ? 'I' : 'O', _id, _name);
  }

private:
  const char *_name;
  bool _value;
  uint32_t _id;
  const core::logger::ILogger &_logger;
};

static const char *TAG = "main";

const core::logger::ILogger &init_logger(const Timestamp &ts)
{
  using namespace core::logger;
  static char buff[256];
  static LoggerSystem ls;
  static Verbosity verr(ESP_LOG_ERROR);
  static Configuration cerr(buff, sizeof(buff), verr);
  static Printable perr(cerr, ls, ts, true);

  static Verbosity vwrn(ESP_LOG_WARN);
  static Configuration cwrn(buff, sizeof(buff), vwrn);
  static Printable pwrn(cwrn, ls, ts, true);

  static Verbosity vinf(ESP_LOG_INFO);
  static Configuration cinf(buff, sizeof(buff), vinf);
  static Printable pinf(cinf, ls, ts, true);

  static Verbosity vdbg(ESP_LOG_DEBUG);
  static Configuration cdbg(buff, sizeof(buff), vdbg);
  static Printable pdbg(cdbg, ls, ts, true);
  static Dumpable pd(cdbg, ls, ts, true);

  static Verbosity vvrb(ESP_LOG_VERBOSE);
  static Configuration cvrb(buff, sizeof(buff), vvrb);
  static Printable pvrb(cvrb, ls, ts, true);

  static const Logger logger(perr, pwrn, pinf, pdbg, pvrb, pd);
  return logger;
}

infra::Status &init_status(const Timestamp &ts, const core::logger::ILogger &logger)
{
  static Actuator led_coin(CONFIG_COIN_LED, false);
  static MyGpio led_wifi_softap("wifi-softap-led", 3, logger);
  static MyGpio led_wifi_station("wifi-station-led", 4, logger);
  static MyGpio led_wifi_station_connected("wifi-softap-connected-led", 5, logger);

  static infra::Status status(logger, ts,
                              led_coin,
                              led_wifi_softap,
                              led_wifi_station,
                              led_wifi_station_connected);
  return status;
}

strategy::StepRunner &init_step_runner(const Timestamp &ts, const core::logger::ILogger &logger, infra::Config &config)
{
  // static MyGpio coin_pin("coin-ctrl", 1, logger);
  const infra::CoinPulseProps *props = config.coin_pulse_props();
  if (!props->is_valid())
  {
    config.persistency().reset_default(infra::Config::PersistencyId::COIN_PULSE_PROPS);
  }

  logger.inf().log(TAG, "count:%d, duration:%d", props->count, props->duration);

  static Actuator coin_pin(CONFIG_COIN_OUT, false);

  static strategy::StepWait swt(props->duration, ts);
  static strategy::StepActuatorOff saoff(coin_pin);
  static strategy::StepActuatorOn saon(coin_pin);
  static strategy::IStep *st[17];
  memset(st, 0, sizeof(st));

  for (uint8_t count = 0, index = 0; count < props->count; count++)
  {
    if (index < sizeof(st) / sizeof(st[0]))
    {
      st[index++] = &saon;
    }
    if (index < sizeof(st) / sizeof(st[0]))
    {
      st[index++] = &swt;
    }
    if (index < sizeof(st) / sizeof(st[0]))
    {
      st[index++] = &saoff;
    }
    if (index < sizeof(st) / sizeof(st[0]))
    {
      st[index++] = &swt;
    }
  }

  static strategy::StepRunner step_runner(st, sizeof(st) / sizeof(st[0]));
  return step_runner;
}

struct ScenarioMaker : public strategy::IScenarioMaker
{
  ScenarioMaker(const Timestamp &ts, const core::logger::ILogger &logger, infra::Config &config): _ts(ts), _logger(logger), _config(config) {}
  virtual strategy::StepRunner &make() const override
  {
    // static MyGpio coin_pin("coin-ctrl", 1, logger);
    const infra::CoinPulseProps *props = _config.coin_pulse_props();
    if (!props->is_valid())
    {
      _config.persistency().reset_default(infra::Config::PersistencyId::COIN_PULSE_PROPS);
    }

    _logger.inf().log(TAG, "count:%d, duration:%d", props->count, props->duration);

    static Actuator coin_pin(CONFIG_COIN_OUT, false);

    static strategy::StepWait swt(props->duration, _ts);
    static strategy::StepActuatorOff saoff(coin_pin);
    static strategy::StepActuatorOn saon(coin_pin);
    static strategy::IStep *st[17];
    memset(st, 0, sizeof(st));

    for (uint8_t count = 0, index = 0; count < props->count; count++)
    {
      if (index < sizeof(st) / sizeof(st[0]))
      {
        st[index++] = &saon;
      }
      if (index < sizeof(st) / sizeof(st[0]))
      {
        st[index++] = &swt;
      }
      if (index < sizeof(st) / sizeof(st[0]))
      {
        st[index++] = &saoff;
      }
      if (index < sizeof(st) / sizeof(st[0]))
      {
        st[index++] = &swt;
      }
    }

    static strategy::StepRunner step_runner(st, sizeof(st) / sizeof(st[0]));
    return step_runner;
  }

private:
  const Timestamp &_ts;
  const core::logger::ILogger &_logger;
  infra::Config &_config;
};

persistency::Persistency<infra::Config::PersistencyId> &
init_persistency(const core::logger::ILogger &logger)
{

  const size_t PERSISTENCE_SIZE = 256;
  static uint8_t memory[PERSISTENCE_SIZE];
  static constexpr persistency::Persistency<infra::Config::PersistencyId>::Row persistency_table[] = {
      persistency::Persistency<infra::Config::PersistencyId>::make_persistency_row<infra::WifiSettingsN>(infra::Config::PersistencyId::WIFI_CONFIG_SOFTAP),
      persistency::Persistency<infra::Config::PersistencyId>::make_persistency_row<infra::WifiSettingsN>(infra::Config::PersistencyId::WIFI_CONFIG_STATION),
      persistency::Persistency<infra::Config::PersistencyId>::make_persistency_row<infra::BrokerUrl>(infra::Config::PersistencyId::MQTT_BROKER_URL),
      persistency::Persistency<infra::Config::PersistencyId>::make_persistency_row<infra::CharContainer<MQTT_USERNAME_CAP>>(infra::Config::PersistencyId::MQTT_BROKER_USERNAME),
      persistency::Persistency<infra::Config::PersistencyId>::make_persistency_row<infra::CharContainer<MQTT_PASSWORD_CAP>>(infra::Config::PersistencyId::MQTT_BROKER_PASSWORD),
      persistency::Persistency<infra::Config::PersistencyId>::make_persistency_row<infra::CoinPulseProps>(infra::Config::PersistencyId::COIN_PULSE_PROPS),
      {infra::Config::PersistencyId::_LAST, 0, nullptr},
  };
  static_assert(
      persistency::Persistency<infra::Config::PersistencyId>::check_persistency_table_size(
          persistency_table,
          sizeof(memory)),
      "Error, insufficient memory for persistency table mapping");
  static Flash flash(logger);
  static persistency::Persistency<infra::Config::PersistencyId> storage(persistency_table, memory, sizeof(memory), flash);
  return storage;
}

void set_defaults(infra::Config &config)
{
  infra::CharContainer<20> mac;
  if (WifiModeSoftAp::dump_mac(mac.data_mut(), mac.capacity()))
  {
    config.set_wifi_config_softap(infra::WifiSettingsN(mac.data()));
  }
  else
  {
    config.set_wifi_config_softap(infra::WifiSettingsN("-=MAC-read-failed=-"));
  }
  config.set_broker_url("mqtt://192.168.0.150:1883");
  config.set_wifi_config_station(infra::WifiSettingsN("TP-LINK_52F1_2.4G", "29132423"));
}

extern "C"
{
  void app_main(void)
  {
    ESP_ERROR_CHECK(esp_event_loop_create_default());
    esp_log_level_set(service::coin::CoinService<0>::TAG, ESP_LOG_VERBOSE); // set all components to ERROR level
    Timestamp ts;
    auto &logger = init_logger(ts);
    auto &status = init_status(ts, logger);
    auto &storage = init_persistency(logger);
    infra::Config config(storage);
    // auto &steps = init_step_runner(ts, logger, config);
    ButtonInput btn_input(0);
    core::io::Button<2> btn(btn_input, ts);

    if (!storage.load())
    {
      storage.reset_default_all();
      set_defaults(config);
      storage.save();
    }

    const infra::WifiSettingsN *sap_cfg = config.wifi_config_softap();
    const infra::WifiSettingsN *sta_cfg = config.wifi_config_station();
    WifiModeStation sta(wifi::WifiStationConfiguration(sta_cfg->ssid().data(), sta_cfg->pswd().data()), logger);
    WifiModeSoftAp sap(wifi::WifiSoftApConfiguration(sap_cfg->ssid().data(), sap_cfg->pswd().data()), logger);
    WifiMode mode(sap, sta);
    Wifi wifi(mode, logger);
    service::mode::ModeService modsvc(logger, wifi, status);
    service::web::WebService ws(logger, config);
    Router r(logger, ws);
    Mqtt mqtt(sta, config, logger);
    service::coin::CoinService<384> coinsvc(mqtt, status, ScenarioMaker(ts, logger, config), logger, ts, config);

    mqtt.add_subscriber(&coinsvc, coinsvc.sub_topic());
    ws.add_subscriber(&modsvc, service::web::event::EventWifiCredSaved::event_name());
    modsvc.add_subscriber(&ws, service::mode::event::EventWifiModeChanged::event_name());
    modsvc.add_subscriber(&coinsvc, service::mode::event::EventWifiModeChanged::event_name());
    btn.add_subscriber(&modsvc, core::io::event::EventButtonClicked::event_name());
    btn.add_subscriber(&modsvc, core::io::event::EventButtonHeld::event_name());
    if (sta_cfg->is_ssid_vaid())
    {
      modsvc.switch_to_station();
    }
    else
    {
      modsvc.switch_to_softap();
    }

    size_t ts_log = ts.get();
    while (1)
    {
      if (ts_log < ts.get())
      {
        ts_log = ts.get() + 1000;
        logger.inf().log(TAG, "sap:%d, sta:%d, ip:%d, broker:%d", wifi.is_softap(), wifi.is_station(), wifi.mode().station().is_connected(), mqtt.is_connected());
      }
      modsvc.run();
      coinsvc.run();
      status.run();
      btn.run();
      vTaskDelay(50.0 / portTICK_PERIOD_MS);
    }
  }
}