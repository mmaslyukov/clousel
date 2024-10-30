#pragma once

#include <infrastructure/json/i_json_parser.h>
#include <infrastructure/json/i_json_dumper.h>
#include <infrastructure/json/my_jsmn_ext.h>
#include <infrastructure/json/jsmn.h>
// #include <infrastructure/config/config.h>

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
      struct Config : public infra::IJsonParser, public infra::IJsonDumper
      {
        constexpr Config()
            : coin_pulse_count("CoinPulseCnt"),
              coin_pulse_duration("CoinPulseDur"),
              broker_url("BrokerUrl"),
              broker_username("BrokerUsername"),
              broker_password("BrokerPassword")
        {
        }
        constexpr Config(uint8_t pulse_count, uint32_t pulse_duration, const infra::BrokerUrl* url) : Config()
        {
          coin_pulse_count.value = pulse_count;
          coin_pulse_duration.value = pulse_duration;
          if (url)
          {
            broker_url.value.replace(url->data());
          }
        }

        virtual bool parse(const char *json_str, size_t len) override
        {
          jsmn_parser p;
          jsmntok_t t[16];
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
            // infra::CharContainer<10> str;
            for (int32_t i = 1; i < tokens_count; i++)
            {
              if (infra::jsoneq(json_str, &t[i], coin_pulse_count.name) == 0)
              {
                // res = res && str.replace(json_str + t[i + 1].start, t[i + 1].end - t[i + 1].start);
                // coin_pulse_count.value = atoi(str.data());
                coin_pulse_count.value = strtol(json_str + t[i + 1].start, nullptr, 10);
                i++;
              }
              else if (infra::jsoneq(json_str, &t[i], coin_pulse_duration.name) == 0)
              {
                // res = res && str.replace(json_str + t[i + 1].start, t[i + 1].end - t[i + 1].start);
                // coin_pulse_duration.value = atoi(str.data());
                coin_pulse_duration.value = strtol(json_str + t[i + 1].start, nullptr, 10);
                i++;
              }
              else if (infra::jsoneq(json_str, &t[i], broker_url.name) == 0)
              {
                res = res && broker_url.value.replace(json_str + t[i + 1].start, t[i + 1].end - t[i + 1].start);
                i++;
              }
              else if (infra::jsoneq(json_str, &t[i], broker_username.name) == 0)
              {
                res = res && broker_username.value.replace(json_str + t[i + 1].start, t[i + 1].end - t[i + 1].start);
                i++;
              }
              else if (infra::jsoneq(json_str, &t[i], broker_password.name) == 0)
              {
                res = res && broker_password.value.replace(json_str + t[i + 1].start, t[i + 1].end - t[i + 1].start);
                i++;
              }
            }
          } while (false);
          return res;
        }
        virtual size_t dump(char *json_str, size_t cap) const override
        {
          size_t shift = 0;
          if (json_str)
          {
            shift = snprintf(json_str, cap, R"({"%s":%d,"%s":%lu,"%s":"%s"})",
                             coin_pulse_count.name, coin_pulse_count.value,
                             coin_pulse_duration.name, coin_pulse_duration.value,
                             broker_url.name, broker_url.value.data());
          }
          return shift;
        }
        void clear()
        {
          coin_pulse_count.value = 0;
          coin_pulse_duration.value = 0;
          broker_url.value.clear();
          broker_username.value.clear();
          broker_password.value.clear();
        }
        infra::Named<uint8_t> coin_pulse_count;
        infra::Named<uint32_t> coin_pulse_duration;
        infra::Named<infra::BrokerUrl> broker_url;
        infra::Named<infra::BrokerUsername> broker_username;
        infra::Named<infra::BrokerPassword> broker_password;
      };
    }
  }
}
