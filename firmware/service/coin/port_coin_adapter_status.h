#pragma once
namespace service
{
  namespace coin
  {
    struct IPortAdapterStatus
    {
      virtual void led_coin_blink() = 0;
    };
  } // namespace mode
}
