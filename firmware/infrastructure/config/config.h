#pragma once
#include <framework/persistency.h>
#include <service/coin/port_coin_adapter_config.h>
#include <service/mode/port_mode_adapter_config.h>
#include <service/web/port_web_adapter_config.h>
// #include <>

#include "entry/carousel_id.h"
#include "entry/wifi_settings.h"
#include "entry/topic_sub.h"
#include "entry/topic_pub.h"
#include "entry/broker_url.h"
#include "entry/coin_pulse.h"


// #define MQTT_USERNAME_CAP 17 
// #define MQTT_PASSWORD_CAP 17


namespace infra
{
  class Config
      : public service::coin::IPortAdapterConfig,
        public service::mode::IPortAdapterConfig,
        public service::web::IPortAdapterConfig
  {
  public:
    enum PersistencyId
    {
      WIFI_CONFIG_STATION,
      WIFI_CONFIG_SOFTAP,
      MQTT_BROKER_URL,
      MQTT_BROKER_USERNAME,
      MQTT_BROKER_PASSWORD,
      COIN_PULSE_PROPS,
      CAROUSEL_ID,
      _LAST
    };

    Config(persistency::Persistency<PersistencyId> &persistency)
        : _persistency(persistency) {}

    virtual bool load() const override
    {
      return _persistency.load();
    }

    virtual bool save() const override
    {
      return _persistency.save();
    }

    virtual const infra::WifiSettingsN *wifi_config_station() override
    {
      return _persistency.get<WifiSettingsN>(PersistencyId::WIFI_CONFIG_STATION);
    }

    virtual bool set_wifi_config_station(const infra::WifiSettingsN &settings) override
    {
      return _persistency.write(PersistencyId::WIFI_CONFIG_STATION, &settings);
    }

    virtual const infra::WifiSettingsN *wifi_config_softap() override
    {
      return _persistency.get<WifiSettingsN>(PersistencyId::WIFI_CONFIG_SOFTAP);
    }

    virtual bool set_wifi_config_softap(const infra::WifiSettingsN &settings) override
    {
      return _persistency.write(PersistencyId::WIFI_CONFIG_SOFTAP, &settings);
    }

    virtual const CoinPulseProps* coin_pulse_props() const override
    {
      return _persistency.get<CoinPulseProps>(PersistencyId::COIN_PULSE_PROPS);
    }

    virtual bool set_coin_pulse_props(const CoinPulseProps& props) override
    {
      return _persistency.write(PersistencyId::COIN_PULSE_PROPS, &props);
    }

    virtual const char *root_sub_topic() const override
    {
      return TopicSub::default_topic();
    }

    virtual const char *root_pub_topic() const override
    {
      return TopicPub::default_topic();
    }

    virtual const CarouselId *carousel_id() const override
    {
      return _persistency.get<CarouselId>(PersistencyId::CAROUSEL_ID);
    }

    const infra::BrokerUsername *broker_username() const
    {
      return _persistency.get<BrokerUsername>(PersistencyId::MQTT_BROKER_USERNAME);
    }

    virtual bool set_broker_username(const BrokerUsername &username) override
    {
      return _persistency.write(PersistencyId::MQTT_BROKER_USERNAME, &username);
    }

    const infra::BrokerPassword *broker_password() const
    {
      return _persistency.get<BrokerPassword>(PersistencyId::MQTT_BROKER_PASSWORD);
    }

    virtual bool set_broker_password(const BrokerPassword &password) override
    {
      return _persistency.write(PersistencyId::MQTT_BROKER_USERNAME, &password);
    }

    virtual const BrokerUrl *broker_url() const override
    {
      return _persistency.get<BrokerUrl>(PersistencyId::MQTT_BROKER_URL);
    }

    virtual bool set_broker_url(const BrokerUrl& url) const override
    {
      return _persistency.write(PersistencyId::MQTT_BROKER_URL, &url);
    }

    const char *broker_client_id() const
    {
      return carousel_id()->value();
    }

    virtual size_t heartbeat_tm_ms() const override
    {
      return 60 * 1000;
    }

    persistency::Persistency<PersistencyId> & persistency()
    {
      return _persistency;
    }

  private:
    persistency::Persistency<PersistencyId> &_persistency;
  };

} // namespace infra
