#pragma once

#include <infrastructure/json/i_json_parser.h>
#include <infrastructure/json/i_json_dumper.h>
#include <infrastructure/json/my_jsmn_ext.h>
#include <infrastructure/json/jsmn.h>

#include "defined.h"

#include <cstdlib>
#include <stdio.h>

namespace service
{
  namespace coin
  {
    namespace msg
    {
      struct Command : public infra::IJsonParser, public infra::IJsonDumper
      {
        constexpr Command()
            : carousel_id(BMF_CAROUSEL_ID),
              event_id(BMF_EVENT_ID),
              type(BMF_TYPE),
              sequence_num(BMF_SEQUENCE_NUM) {}

        virtual bool parse(const char *json_str, size_t len) override
        {
          jsmn_parser p;
          jsmntok_t t[64];
          int32_t tokens_count = 0;
          jsmn_init(&p);
          bool res = true;
          do
          {
            if ((tokens_count = jsmn_parse(&p, json_str, len, t, sizeof(t) / sizeof(t[0]))) < 0)
            {
              break;
            }
            if (tokens_count < 1 || t[0].type != JSMN_OBJECT)
            {
              break;
            }

            for (int32_t i = 1; i < tokens_count; i++)
            {
              if (infra::jsoneq(json_str, &t[i], carousel_id.name) == 0)
              {
                res = res && carousel_id.value.replace(json_str + t[i + 1].start, t[i + 1].end - t[i + 1].start);
                i++;
              }
              else if (infra::jsoneq(json_str, &t[i], event_id.name) == 0)
              {
                res = res && event_id.value.replace(json_str + t[i + 1].start, t[i + 1].end - t[i + 1].start);
                i++;
              }
              else if (infra::jsoneq(json_str, &t[i], type.name) == 0)
              {
                res = res && type.value.replace(json_str + t[i + 1].start, t[i + 1].end - t[i + 1].start);
                i++;
              }
              else if (infra::jsoneq(json_str, &t[i], sequence_num.name) == 0)
              {
                sequence_num.value = strtol(json_str + t[i + 1].start, nullptr, 10);
                i++;
              }
            }
          } while (false);
          return res;
        }
        virtual size_t dump(char *json_str, size_t cap) const override
        {
          int shift = 0;
          if (json_str)
          {
            shift = snprintf(json_str, cap, R"({"%s":"%s","%s":"%s","%s":"%s","%s":%d})",
                                  carousel_id.name, carousel_id.value.data(),
                                  event_id.name, event_id.value.data(),
                                  type.name, type.value.data(),
                                  sequence_num.name, (int)sequence_num.value);
            //  printf("cap:%d, len:%d\n", cap, len);
          }
          return shift;
        }
        void clear()
        {
          carousel_id.value.clear();
          event_id.value.clear();
          type.value.clear();
          sequence_num.value = 0;
        }

        infra::Named<CarouselIdContainer> carousel_id;
        infra::Named<EventIdContainer> event_id;
        infra::Named<TypeContainer> type;
        infra::Named<uint32_t> sequence_num;
      };
    }
  }
}
