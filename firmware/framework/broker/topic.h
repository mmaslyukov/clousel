#pragma once
#include <stdint.h>

namespace broker
{
  template <size_t N>
  class Topic
  {
  public:
    Topic(const char *name)
    {
      memset(_topic, 0, N);
      _len = Topic::copy(_topic, N, name);
    }

    static size_t copy(char *to, size_t n, const char *from)
    {
      size_t len = strlen(from);
      size_t i = 0;
      for (; i < len; i++)
      {
        if (i < n - 1)
        {
          to[i] = from[i];
        }
      }
      return i;
    }

    static bool contain(const char *a, const char *b)
    {
      bool res = true;
      size_t len = strlen(b);
      for (size_t i = 0; i < (N - 1) && (i < len); i++)
      {
        if (a[i] != b[i])
        {
          res = false;
          break;
        }
      }
      return res;
    }

    bool append(const char *name)
    {
      if (_len < N - 1)
      {
        _topic[_len++] = '/';
        _topic[_len] = '\0';
      }
      _len += Topic::copy(&_topic[_len], N - _len, name);
      return _len > sizeof(_topic) ? false : true;
    }

    bool contain(const char *name) const
    {
      return Topic::contain(_topic, name);
    }
    
    bool part_of(const char *name) const
    {
      return Topic::contain(name, _topic);
    }

    const char *get() const
    {
      return _topic;
    }

    const size_t len() const
    {
      return _len;
    }

  private:
    uint8_t _len;
    char _topic[N];
  };

} // namespace broker
