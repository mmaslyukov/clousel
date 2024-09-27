#pragma once
#include <stdint.h>
#include <cstring>

#include "..\gen.h"

namespace infra
{
  struct CarouselId
  {
    // Uuid(const char *value)
    // {
    //   strcpy_s(_value, sizeof(_value), value);
    // }
    const char *value() const
    {
      return _value;
    }

  private:
    // char _value[37]; // 36 + 1
    static constexpr const char* _value = __CAROUSEL_ID;
  };
}