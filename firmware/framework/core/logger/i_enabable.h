#pragma once

namespace core
{
  namespace logger
  {
    
    struct IEnabable
    {
      virtual void enable() = 0;
      virtual void disable() = 0;
    };
    
  }
}