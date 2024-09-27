#pragma once
#include <framework/persistency.h>

#include "char_container.h"

#include <stdint.h>
#include <cstring>

namespace infra
{
  
#define BROKER_URL_CAP 30
#define DEFAULT_BROKER_URL "tcp://192.168.0.150:1883"
constexpr const size_t broker_url_capacity = persistency::max(sizeof(DEFAULT_BROKER_URL), BROKER_URL_CAP);

  struct BrokerUrl : public CharContainer<broker_url_capacity>
  {
    
    BrokerUrl(const char *value) : CharContainer<broker_url_capacity>(value) {}
    BrokerUrl() : CharContainer<broker_url_capacity>(default_url()) {}
    static const char *default_url()
    {
      return "tcp://192.168.0.150:1883";
    }
  };
} // namespace infra
