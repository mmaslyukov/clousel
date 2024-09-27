#pragma once
#include <stdint.h>
#include <cstring>
namespace infra
{
  template <size_t N>
  class CharContainer
  {
  public:
    constexpr CharContainer() : _value{0}, _len(0) {} // set the first el as '\0'
    CharContainer(const char *value)
    {
      if (value)
      {
        _len = strcpy_s(_value, sizeof(_value), value);
      }
    }
    const char *data() const
    {
      return _value;
    }
    char *data_mut()
    {
      return _value;
    }
    bool eq(const char *str)
    {
      return strcmp(_value, str) == 0 ? true : false;
    }

    bool empty() const
    {
      return _len == 0;
    }

    bool append(const char *data)
    {
      if (strlen(len) < capacity())
      {
        _len += snprintf(&_value[_len], capacity(), "%s", data);
      }
      else
      {
        return false;
      }
      return false;
    }
    bool replace(const char *data)
    {
      if (strlen(len) < capacity())
      {
        _len = strcpy_s(_value, capacity(), data);
        return true;
      }
      else
      {
        return false;
      }
    }
    bool replace(const char *data, const size_t len)
    {
      if (len < capacity())
      {
        _len = memcpy_s(_value, capacity(), data, len);
        return true;
      }
      else
      {
        return false;
      }
    }
    inline size_t capacity() const
    {
      return sizeof(_value);
    }

  private:
    char _value[N];
    size_t _len;
  };
}