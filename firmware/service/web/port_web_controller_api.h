#pragma once

#include <infrastructure/config/entry/wifi_settings.h>

#include "response.h"

namespace service
{
  namespace web
  {
    struct IPortWebControllerApi
    {
      virtual bool submit(const infra::WifiSettingsN& settings) const = 0;
    };
    
    
  } // namespace web
  
  
} // namespace service
