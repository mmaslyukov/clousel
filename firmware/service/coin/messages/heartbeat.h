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
            : type(TypeContainer(BMT_HEARTBEAT), "Type"),
              carousel_id(CarouselIdContainer(carousel_id), "CarouselId"),
              sequence_num(sequence_num, "SequenceNum") {}
              // virtual const core::error::IError& dump(char *json_str, size_t cap) {
              //   static core::error::ErrorConst ec;
              //   return ec;
              // };
        virtual bool dump(char *json_str, size_t cap) const override 
        {
          if (json_str)
          {
            snprintf(json_str, cap, R"({"%s":"%s","%s":"%s","%s":%d})",
                     type.name, type.value.data(),
                     carousel_id.name, carousel_id.value.data(),
                     sequence_num.name, sequence_num.value);
            return true;
          }
          return false;
        }
        infra::Named<TypeContainer> type;
        infra::Named<CarouselIdContainer> carousel_id;
        infra::Named<uint32_t> sequence_num;
      };
    }
  }
}