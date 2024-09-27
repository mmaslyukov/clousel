#pragma once

#include <infrastructure/json/i_json_parser.h>
#include <infrastructure/json/i_json_dumper.h>
#include <infrastructure/json/my_jsmn_ext.h>
#include <infrastructure/json/jsmn.h>

#include "defined.h"

#include <cstdlib>
#include <stdio.h>

// type MessageGeneral struct {
// 	CarouselId  string `json:CarouselId`
// 	SequenceNum int    `json:SequenceNum`
// 	EventId     string `json:EventId`
// 	Type        string `json:Type`
// }

// type MessageCommand struct {
// 	MessageGeneral
// 	Command string `json:Command`
// }

namespace service
{
  namespace coin
  {
    namespace msg
    {
      struct Command : public infra::IJsonParser, public infra::IJsonDumper
      {
        Command()
            : carousel_id("CarouselId"),
              event_id("EventId"),
              type("Type"),
              command("Command"),
              sequence_num("SequenceNum") {}

        virtual bool parse(const char *json_str) override
        {
          jsmn_parser p;
          jsmntok_t t[64];
          int32_t tokens_count = 0;
          jsmn_init(&p);
          do
          {
            if ((tokens_count = jsmn_parse(&p, json_str, strlen(json_str), t, sizeof(t) / sizeof(t[0]))) < 0)
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
                carousel_id.value.replace(json_str + t[i + 1].start, t[i + 1].end - t[i + 1].start);
                i++;
              }
              else if (infra::jsoneq(json_str, &t[i], event_id.name) == 0)
              {
                event_id.value.replace(json_str + t[i + 1].start, t[i + 1].end - t[i + 1].start);
                i++;
              }
              else if (infra::jsoneq(json_str, &t[i], type.name) == 0)
              {
                type.value.replace(json_str + t[i + 1].start, t[i + 1].end - t[i + 1].start);
                i++;
              }
              else if (infra::jsoneq(json_str, &t[i], command.name) == 0)
              {
                command.value.replace(json_str + t[i + 1].start, t[i + 1].end - t[i + 1].start);
                i++;
              }
              else if (infra::jsoneq(json_str, &t[i], sequence_num.name) == 0)
              {
                sequence_num.value = strtol(json_str + t[i + 1].start, nullptr, 10);
                i++;
              }
            }
            return true;
          } while (false);
          return false;
        }
        virtual bool dump(char *json_str, size_t cap) const override
        {
          if (json_str)
          {
            size_t len = snprintf(json_str, cap, R"({"%s":"%s","%s":"%s","%s":"%s","%s":"%s","%s":%d})",
                                  carousel_id.name, carousel_id.value.data(),
                                  event_id.name, event_id.value.data(),
                                  type.name, type.value.data(),
                                  command.name, command.value.data(),
                                  sequence_num.name, sequence_num.value);
            //  printf("cap:%d, len:%d\n", cap, len);
            return len > cap - 1 ? false : true;
          }
          return false;
        }

        infra::Named<CarouselIdContainer> carousel_id;
        infra::Named<EventIdContainer> event_id;
        infra::Named<TypeContainer> type;
        infra::Named<CommandContainer> command;
        infra::Named<uint32_t> sequence_num;
      };
    }
  }
}
