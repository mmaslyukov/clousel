#pragma once
namespace infra
{
  template <typename T>
  struct Named
  {    
    constexpr Named(const char *name)
        : value(), name(name) {}
    constexpr Named(const T &value, const char *name)
        : value(value), name(name) {}
    T value;
    const char *name;
  };
}