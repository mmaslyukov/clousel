#pragma once
// #include <infrastructure/config/entry/char_container.h>
// #include <infrastructure/json/named.h>
#include <infrastructure/json/i_json_dumper.h>
#include <infrastructure/json/jsmn.h>

#include "defined.h"

#include <stdio.h>
// type EventMinimal struct {
// 	Type        string `json:"Type"`
// 	CarouselId  string `json:"CarouselId"`
// 	SequenceNum int    `json:"SequenceNum"`
// }

// type EventAck struct {
// 	EventMinimal
// 	CorrelationId string `json:"CorrelationId"`

// 	// optional field for showing status of the last command
// 	Error string `json:"Error"`
// }

namespace service
{
  namespace coin
  {
    namespace msg
    {
      struct ResponseAck : public infra::IJsonDumper
      {
        ResponseAck(const char *carousel_id, const char *correlation_id, const char* error = nullptr, uint32_t sequence_num = 0)
            : type(TypeContainer(BMT_ACK), "Type"),
              carousel_id(CarouselIdContainer(carousel_id), "CarouselId"),
              correlation_id(EventIdContainer(correlation_id), "CorrelationId"),
              sequence_num(sequence_num, "SequenceNum"),
              error_str(error, "Error") { }
        virtual bool dump(char *json_str, size_t cap) const override
        {
          if (json_str)
          {
            snprintf(json_str, cap, R"({"%s":"%s","%s":"%s","%s":"%s","%s":%d})",
                     type.name, type.value.data(),
                     carousel_id.name, carousel_id.value.data(),
                     correlation_id.name, correlation_id.value.data(),
                     sequence_num.name, sequence_num.value);
            return true;
          }
          return false;
        }
        infra::Named<TypeContainer> type;
        infra::Named<CarouselIdContainer> carousel_id;
        infra::Named<EventIdContainer> correlation_id;
        infra::Named<uint32_t> sequence_num;
        infra::Named<ErrorContainer> error_str;
        
      };
    }
  }
}