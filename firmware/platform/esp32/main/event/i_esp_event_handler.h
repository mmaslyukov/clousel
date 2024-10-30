#pragma once
#include "esp_event_base.h"

struct IEspEvetHandler
{
  virtual bool handle(esp_event_base_t event_base, int32_t event_id, void *event_data) = 0;
};