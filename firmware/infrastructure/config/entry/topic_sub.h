#pragma once
#include <stdint.h>
#include <cstring>
#include "char_container.h"
#include <framework/persistency.h>
namespace infra
{ 

#define ROOT_TOPIC_SUB_CAP 30
#define DEFAULT_ROOT_TOPIC_SUB "/clousel/cloud"
constexpr const size_t topic_sub_capacity = persistency::max(sizeof(DEFAULT_ROOT_TOPIC_SUB), ROOT_TOPIC_SUB_CAP);

  struct TopicSub : public CharContainer<topic_sub_capacity>
  {
    TopicSub(const char *value) : CharContainer<topic_sub_capacity>(value) {}
    TopicSub() : CharContainer<topic_sub_capacity>(default_topic()) {}

    static constexpr const char *default_topic()
    {
      return DEFAULT_ROOT_TOPIC_SUB;//"/clousel/cloud";
    }
  };
} // namespace infra
