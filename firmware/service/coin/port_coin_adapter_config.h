#pragma once
#include <infrastructure/config/entry/carousel_id.h>
#include <infrastructure/config/entry/coin_pulse.h>
#include <infrastructure/config/entry/broker_url.h>

namespace service
{
  namespace coin
  {
    struct IPortAdapterConfig
    {
      virtual bool save() const = 0;
      virtual const char *root_sub_topic() const = 0;
      virtual const char *root_pub_topic() const = 0;
      virtual const infra::CarouselId *carousel_id() const = 0;
      virtual size_t heartbeat_tm_ms() const = 0;

      virtual const infra::CoinPulseProps *coin_pulse_props() const = 0;
      virtual bool set_coin_pulse_props(const infra::CoinPulseProps &props) = 0;
      virtual bool set_broker_username(const infra::BrokerUsername &username) = 0;
      virtual bool set_broker_password(const infra::BrokerPassword &password) = 0;
      virtual bool set_broker_url(const infra::BrokerUrl &url) const = 0;
      virtual const infra::BrokerUrl *broker_url() const = 0;
    };
  }
}