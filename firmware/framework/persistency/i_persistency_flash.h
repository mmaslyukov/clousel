#pragma once
#include <stdint.h>

namespace persistency
{
  struct IPersistencyFlash
  {
    
    // virtual bool load()  = 0;
    // virtual bool save()  = 0;
    
    virtual bool load(uint8_t *memory, size_t size) const = 0;
    virtual bool save(const uint8_t *memory, size_t size) const = 0;
    /* data */
  };

} // namespace persistency