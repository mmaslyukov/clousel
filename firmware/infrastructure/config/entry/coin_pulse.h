#pragma once
// #include <stddef.h>
#include <stdint.h>
namespace infra
{
#pragma pack(push, 1)
  struct CoinPulseProps
  {
    constexpr CoinPulseProps() : count(2), duration(100) {}
    constexpr CoinPulseProps(uint8_t count, uint32_t duration) : count(count), duration(duration) {}
    bool is_valid() const
    {
      bool res = true;
      if (duration < 50 || duration > 150)
      {
        res = false;
      }
      if (count < 1 || count > 4)
      {
        res = false;
      }
      return res;
    }
    uint8_t count;
    uint32_t duration;
  };
#pragma pack(pop)
}