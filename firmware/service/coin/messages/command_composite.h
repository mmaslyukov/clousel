#pragma once

#include "command.h"
#include "config.h"

#include <cstdlib>
#include <stdio.h>

namespace service
{
  namespace coin
  {
    namespace msg
    {
      struct CommandComposite : public infra::IJsonParser //, public infra::IJsonDumper
      {
        constexpr CommandComposite() : config("Config") {}
        virtual bool parse(const char *json_str, size_t len) override
        {
          bool res = true;
          res = res && general.parse(json_str, len);
          if (general.command.value.eq(BMTC_CONFIG_WRITE))
          {
            res = res && parse_config(json_str, len);
          }
          return res;
        }
        void clear()
        {
          general.clear();
          config.value.clear();
        }
      private:
        bool parse_config(const char *json_str, size_t len)
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
              res = false;
              break;
            }
            if (tokens_count < 1 || t[0].type != JSMN_OBJECT)
            {
              res = false;
              break;
            }
            for (int32_t i = 1; i < tokens_count; i++)
            {
              if (infra::jsoneq(json_str, &t[i], config.name) == 0)
              {
                res = res && config.value.parse(json_str + t[i + 1].start, t[i + 1].end - t[i + 1].start);
                i++;
              }
            }
          } while (false);
          return res;
        }


      public:
        Command general;
        infra::Named<Config> config;
      };
    }
  }
}
