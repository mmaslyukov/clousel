#pragma once
namespace infra
{
  struct IJsonParser
  {
    virtual bool parse(const char *json_str, size_t len) = 0;
  };
}