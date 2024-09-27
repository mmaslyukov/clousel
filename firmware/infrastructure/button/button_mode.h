#pragma once
#include <framework/core/i_runnable.h>
#include <framework/core/i_io.h>
#include <service/mode/port_mode_controller_button.h>

namespace infra
{
  class ModeButton : public core::IRunnable, public core::io::IButton
  {
    enum State
    {
      UNKNOWN,
      PRESSED,
      HELD,
    };

  public:
    constexpr ModeButton(
        service::mode::IPortButtonController &handler,
                         const size_t hold_ms)
        : _handler(handler), _hold_ms(hold_ms), _ms(0) {}

    virtual void run()
    {
      // if (_btn.held())
      // {
      //   _handler.held();
      // }
      // else if (_btn.pressed())
      // {
      //   _handler.pressed();
      // }
      // if (_btn.get())
      // {
      //   switch (_state)
      //   {
      //   case UNKNOWN:
      //     if (_ts.get() - _ms > _debounce_ms)
      //     {
      //       _handler.pressed();
      //       _state = PRESSED;
      //     }
      //   case PRESSED:
      //     if (_ts.get() - _ms > _hold_ms)
      //     {
      //       _handler.held();
      //       _state = HELD;
      //     }
      //     break;
      //   case HELD:
      //   default:
      //     break;
      //   }
      // }
      // else
      // {
      //   _ms = _ts.get();
      // }
    }
      virtual uint32_t id() const = 0;
      virtual bool clicked() const = 0;
      virtual bool held() const = 0;
  private:
    service::mode::IPortButtonController &_handler;
    const size_t _hold_ms;
    size_t _ms;
    const size_t _debounce_ms = 300;
  };
} // namespace infra
