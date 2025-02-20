#pragma once
#include <cstring>

#include "jsmn.h"

namespace infra
{

  int jsoneq(const char *json, jsmntok_t *tok, const char *s)
  {
    if (tok->type == JSMN_STRING && (int)strlen(s) == tok->end - tok->start &&
        strncmp(json + tok->start, s, tok->end - tok->start) == 0)
    {
      return 0;
    }
    return -1;
  }
}