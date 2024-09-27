#pragma once
#include <framework/core/error.h>
namespace infra
{
  struct IJsonDumper
  {
    virtual bool dump(char *json_str, size_t cap) const = 0;
  };
}
