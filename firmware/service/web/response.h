#pragma once
#include <stdint.h>
namespace service
{
  namespace web
  {
    struct Response
    {
      uint32_t status;
      const char* data;
      /* data */
    };
    
  }
  }