#pragma once
#include <stdint.h>
#include <cstring>
#include <framework/util/util.h>

namespace infra
{
  template <size_t N>
  class CharContainer
  {
  public:
    constexpr CharContainer() : _len(0), _value{0} {} // set the first el as '\0'
    CharContainer(const char *value) : _len(0)
    {
      if (value)
      {
        if (!core::util::strcpy_s(_value, capacity(), value))
        {
          _len = strlen(value);
        }
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
    size_t len() const
    {
      return _len;
    }
    bool eq(const char *str)
    {
      return strcmp(_value, str) == 0 ? true : false;
    }

    bool empty() const
    {
      return _len == 0;
    }
    void clear()
    {
      memset(_value, 0, capacity());
      _len = 0;
    }

    bool append(const char *data)
    {
      if (strlen(_len) < capacity())
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
      if (strlen(data) < capacity())
      {
        clear();
        core::util::strcpy_s(_value, capacity(), data);
        _len = strlen(data);
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
        clear();
        core::util::mmemcpy_s(_value, capacity(), data, len);
        _len = len;
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
    size_t _len;
    char _value[N];
  };
}