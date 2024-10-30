#pragma once
#include <stdint.h>

namespace broker
{
  namespace topic
  {

    bool contains(const char *a, size_t cap, const char *b, size_t len)
    {
      bool res = true;
      for (size_t i = 0; i < (cap - 1) && (i < len); i++)
      {
        if (a[i] != b[i])
        {
          res = false;
          break;
        }
      }
      return res;
    }
    static size_t copy(char *to, size_t n, const char *from, size_t len)
    {
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
  } // namespace topic

  struct ITopic
  {
    virtual bool contains(const ITopic &topic) const = 0;
    virtual bool contains(const char *topic) const = 0;
    virtual bool append(const ITopic &topic) = 0;
    virtual bool append(const char *topic) = 0;
    virtual bool part_of(const char* topic) const = 0;
    virtual bool part_of(const ITopic &topic) const = 0;
    virtual const char *get() const = 0;
    virtual size_t len() const = 0;
  };

  class TopicRef : public ITopic
  {
  public:
    TopicRef(const char *name) : _len(strlen(name)), _topic(name) {}
    TopicRef(const char *name, size_t len) : _len(len), _topic(name) {}
    virtual bool contains(const ITopic &topic) const override
    {
      return topic::contains(_topic, _len, topic.get(), topic.len());
    }
    virtual bool contains(const char *topic) const override
    {
      return topic::contains(_topic, _len, topic, strlen(topic));
    }
    virtual bool append(const ITopic &topic) override
    {
      return false;
    }
    virtual bool append(const char *topic) override
    {
      return false;
    }
    virtual bool part_of(const ITopic &topic) const override
    {
      return topic::contains(topic.get(), topic.len(), _topic, _len);
    }
    virtual bool part_of(const char* topic) const override
    {
      return topic::contains(topic, strlen(topic), _topic, _len);
    }
    virtual const char *get() const override
    {
      return _topic;
    }
    virtual size_t len() const override
    {
      return _len;
    }

  private:
    uint8_t _len;
    const char *_topic;
  };

  template <size_t N>
  class TopicContainer : public ITopic
  {
  public:
    TopicContainer()
    {
      memset(_topic, 0, N);
    }

    TopicContainer(const char *name)
    {
      memset(_topic, 0, N);
      _len = topic::copy(_topic, N, name, strlen(name));
    }

    TopicContainer(const char *name, size_t len)
    {
      memset(_topic, 0, N);
      _len = topic::copy(_topic, N, name, len);
    }
    virtual bool append(const char *topic) override
    {
      return append(topic, strlen(topic));
    }
    virtual bool append(const ITopic &topic) override
    {
      return append(topic.get(), topic.len());
    }

    virtual bool contains(const ITopic &topic) const override
    {
      return topic::contains(_topic, N, topic.get(), topic.len());
    }
    virtual bool contains(const char *topic) const override
    {
      return topic::contains(_topic, N, topic, strlen(topic));
    }

    virtual bool part_of(const ITopic &topic) const override
    {
      return topic::contains(topic.get(), topic.len(), _topic, N);
    }
    virtual bool part_of(const char *topic) const override
    {
      return topic::contains(topic, strlen(topic), _topic, N);
    }

    virtual const char *get() const override
    {
      return _topic;
    }

    virtual size_t len() const override
    {
      return _len;
    }

  private:
    bool append(const char *topic, size_t len)
    {
      if (_len < N - 1)
      {
        _topic[_len++] = '/';
        _topic[_len] = '\0';
      }
      _len += topic::copy(&_topic[_len], N - _len, topic, len);
      return _len > sizeof(_topic) ? false : true;
    }
    // static size_t copy(char *to, size_t n, const char *from)
    // {
    //   return topic::copy(to, n, from, strlen(from));
    // }

  private:
    uint8_t _len;
    char _topic[N];
  };

} // namespace broker
