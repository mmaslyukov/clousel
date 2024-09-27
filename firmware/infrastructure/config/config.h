#pragma once
#include <framework/persistency.h>
#include <service/coin/port_coin_adapter_config.h>
#include <service/mode/port_mode_adapter_config.h>
#include <service/web/port_web_adapter_config.h>

#include "entry/carousel_id.h"
#include "entry/wifi_settings.h"
#include "entry/topic_sub.h"
#include "entry/topic_pub.h"
#include "entry/broker_url.h"

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
    virtual const infra::WifiSettingsN *wifi_config_softap() override
    {
      return _persistency.get<WifiSettingsN>(PersistencyId::WIFI_CONFIG_SOFTAP);
    }

    virtual const bool set_wifi_config_station(const infra::WifiSettingsN &settings) override
    {
      return _persistency.write(PersistencyId::WIFI_CONFIG_STATION, &settings);
    }
    virtual const bool set_wifi_config_softap(const infra::WifiSettingsN &settings) override
    {
      return _persistency.write(PersistencyId::WIFI_CONFIG_SOFTAP, &settings);
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

    const char *broker_url() const
    {
      return BrokerUrl::default_url();
    }

    const char *broker_client_id() const
    {
      return "ExampleClientSub-app";
    }

    virtual const size_t heartbeat_tm_ms() const override
    {
      return 60 * 1000;
    }

  private:
    persistency::Persistency<PersistencyId> &_persistency;
  };

} // namespace infra
