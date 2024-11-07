#pragma once
// #include <infrastructure/config/entry/char_container.h>
// #include <infrastructure/json/named.h>
#include <infrastructure/json/i_json_dumper.h>
#include <infrastructure/json/jsmn.h>
#include "defined.h"

// {
// 	"Type": "EventHeartbeat",
// 	"CarouselId": "550e8400-e29b-41d4-a716-446655440000",
// 	"SequenceNum": 3
// }
namespace service
{
  namespace coin
  {
    namespace msg
    {
      struct Heartbeat : public infra::IJsonDumper
      {
        Heartbeat(const char *carousel_id, uint32_t sequence_num = 0)
            : type(TypeContainer(BMTE_HEARTBEAT), BMF_TYPE),
              carousel_id(CarouselIdContainer(carousel_id), BMF_CAROUSEL_ID),
              sequence_num(sequence_num, BMF_SEQUENCE_NUM) {}
        virtual size_t dump(char *json_str, size_t cap) const override 
        {
          int shift = 0;
          if (json_str)
          {
            shift = snprintf(json_str, cap, R"({"%s":"%s","%s":"%s","%s":%d})",
                     type.name, type.value.data(),
                     carousel_id.name, carousel_id.value.data(),
                     sequence_num.name, (int)sequence_num.value);
          }
          return shift;
        }
        infra::Named<TypeContainer> type;
        infra::Named<CarouselIdContainer> carousel_id;
        infra::Named<uint32_t> sequence_num;
      };
    }
  }
}