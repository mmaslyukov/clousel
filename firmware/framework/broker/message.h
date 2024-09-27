#pragma once
#include <stdint.h>
#include <string.h>

namespace broker
{
  using Token = int;

  struct Message
  {
    public:
    Message(const char *data)
        : data((const uint8_t*)data), size(strlen(data)) {}
    
    constexpr Message()
        : data(Message::_empty), size(sizeof(Message::_empty)) {}
    constexpr Message(const uint8_t *data, const size_t size)
        : data(data), size(size) {}
    constexpr Message(const void *data, const size_t size)
        : data((const uint8_t*)data), size(size) {}
  
    const size_t size;
    const uint8_t *data;
  private:
  static constexpr const uint8_t _empty[] = {0xDE, 0xAD, 0xBE, 0xEF};
  };
}