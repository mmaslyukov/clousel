#pragma once
#include <infrastructure/config/entry/carousel_id.h>

namespace service
{
  namespace coin
  {
    struct IPortAdapterConfig
    {
      virtual const char *root_sub_topic() const = 0;
      virtual const char *root_pub_topic() const = 0;
      virtual const infra::CarouselId *carousel_id() const = 0;
      virtual const size_t heartbeat_tm_ms() const = 0;
    };
  }
}