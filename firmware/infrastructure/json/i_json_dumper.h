#pragma once
#include <framework/core/error.h>
namespace infra
{
  struct IJsonDumper
  {
    virtual size_t dump(char *json_str, size_t cap) const = 0;
  };
}
