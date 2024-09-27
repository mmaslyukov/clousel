#pragma once
namespace infra
{
  struct IJsonParser
  {
    virtual bool parse(const char *json_str) = 0;
  };
}