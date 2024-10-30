#pragma once
#include <stdint.h>
#include <string.h>

namespace broker
{
  // using Token = int;
  struct Token
  {
    constexpr Token() : _result(false), _id(-1) {}
    constexpr Token(bool result, int id) : _result(result), _id(id) {}
    Token &set_result(bool result)
    {
      _result = result;
      return *this;
    }
    bool is_ok() const
    {
      return _result;
    }
    const int &id() const
    {
      return _id;
    }
    int *id_ptr()
    {
      return &_id;
    }
    bool is_valid()
    {
      return _id != -1;
    }
    Token& invalidate()
    {
      _result = false;
      _id = -1;
      return *this;
    }

  private:
    bool _result;
    int _id;
  };

  struct Message
  {
  public:
    Message(const char *data)
        : data((const uint8_t *)data), size(strlen(data)) {}

    constexpr Message()
        : data(Message::_empty), size(sizeof(Message::_empty)) {}
    constexpr Message(const uint8_t *data, const size_t size)
        : data(data), size(size) {}
    constexpr Message(const void *data, const size_t size)
        : data((const uint8_t *)data), size(size) {}

    const uint8_t *data;
    const size_t size;

  private:
    static constexpr const uint8_t _empty[] = {0xDE, 0xAD, 0xBE, 0xEF};
  };
}