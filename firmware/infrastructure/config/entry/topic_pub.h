#pragma once
#include <framework/persistency.h>

#include "char_container.h"

#include <stdint.h>
#include <cstring>

namespace infra
{

#define ROOT_TOPIC_PUB_CAP 30
#define DEFAULT_ROOT_TOPIC_PUB "/clousel/carousel"
constexpr const size_t topic_pub_capacity = persistency::max(sizeof(DEFAULT_ROOT_TOPIC_PUB), ROOT_TOPIC_PUB_CAP);

  struct TopicPub : public CharContainer<topic_pub_capacity>
  {
    TopicPub(const char *value) : CharContainer<topic_pub_capacity>(value) {}
    TopicPub() : CharContainer<topic_pub_capacity>(default_topic()) {}

    static constexpr const char *default_topic()
    {
      return DEFAULT_ROOT_TOPIC_PUB;
    }

  };
} // namespace infra
