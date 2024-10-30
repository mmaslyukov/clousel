#pragma once
#include <stdint.h>
#include <stddef.h>
#include <string.h>

namespace core
{
  namespace util
  {
    template<typename T>
    int32_t mmemcpy_s(T *dst, size_t cap, const T *src, size_t len)
    {
      size_t i = 0;
      for (; i < cap && i < len; i++)
      {
        dst[i] = src[i];
      }
      return i == len ? 0 : -1;
    }
    int32_t strcpy_s(char *dst, size_t cap, const char *src)
    {
      return mmemcpy_s(dst, cap, src, strlen(src) + 1);
    }

  } // namespace util

} // namespace care
