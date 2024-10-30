#pragma once
#include <infrastructure/json/i_json_dumper.h>
#include <infrastructure/json/jsmn.h>

#include "defined.h"
#include "config.h"

#include <stdio.h>

namespace service
{
  namespace coin
  {
    namespace msg
    {
      struct ResponseConfig : public infra::IJsonDumper
      {
        ResponseConfig(const char *carousel_id,
                       const char *correlation_id,
                       const char *error,
                       Config config,
                       uint32_t sequence_num)
            : type(TypeContainer(BMTR_CONFIG), "Type"),
              carousel_id(CarouselIdContainer(carousel_id), "CarouselId"),
              correlation_id(EventIdContainer(correlation_id), "CorrelationId"),
              sequence_num(sequence_num, "SequenceNum"),
              config(config, "Config"),
              error_str(error, "Error") {}
        virtual size_t dump(char *json_str, size_t cap) const override
        {
          int shift = 0;
          if (json_str)
          {
            if (error_str.value.empty())
            {

              shift += snprintf(&json_str[shift], cap - shift, R"({)");
              shift += snprintf(&json_str[shift], cap - shift, R"("%s":"%s","%s":"%s","%s":"%s","%s":%d,"%s":)",
                       type.name, type.value.data(),
                       carousel_id.name, carousel_id.value.data(),
                       correlation_id.name, correlation_id.value.data(),
                       sequence_num.name, (int)sequence_num.value,
                       config.name);
              shift += config.value.dump(&json_str[shift], cap - shift);
              shift += snprintf(&json_str[shift], cap - shift, R"(})");
            }
            else
            {
              shift = snprintf(json_str, cap, R"({"%s":"%s","%s":"%s","%s":"%s","%s":%d,"%s":"%s"})",
                       type.name, type.value.data(),
                       carousel_id.name, carousel_id.value.data(),
                       correlation_id.name, correlation_id.value.data(),
                       sequence_num.name, (int)sequence_num.value,
                       error_str.name, error_str.value.data());
            }
          }
          return shift;
        }
        infra::Named<TypeContainer> type;
        infra::Named<CarouselIdContainer> carousel_id;
        infra::Named<EventIdContainer> correlation_id;
        infra::Named<uint32_t> sequence_num;
        infra::Named<Config> config;
        infra::Named<ErrorContainer> error_str;
      };
    }
  }
}