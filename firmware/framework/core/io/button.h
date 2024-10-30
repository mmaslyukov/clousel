#pragma once
#include <stdint.h>
#include <framework/core/observer.h>
#include <framework/core/i_runnable.h>
#include <framework/core/i_timestamp.h>
#include "i_sensor.h"
#include "event/event_button_clicked.h"
#include "event/event_button_held.h"

namespace core
{
  namespace io
  {
    template <size_t N>
    class Button : public core::observer::Publisher<N>, public core::IRunnable
    {
    private:
      enum ButtonState
      {
        RELEASED,
        PRESSED,
        CLICKED,
        HELD,
        HELD_CONTINUE
      };

    public:
      Button(const ISensor<bool> &input, const core::ITimestamp &ts, size_t click_ms = 300, size_t hold_ms = 3000)
          : _input(input), _ts(ts),
            _button_click_ms(click_ms),
            _button_hold_ms(hold_ms),
            _state(ButtonState::RELEASED) {}
      virtual uint32_t id() const
      {
        return _input.id();
      }
      virtual void run() override
      {
        switch (_state)
        {
        case ButtonState::RELEASED:
          if (_input.get())
          {
            _button_pressed_ts = _ts.get() + _button_click_ms;
            _button_held_ts = _ts.get() + _button_hold_ms;
            _state = ButtonState::PRESSED;
          }
          break;
        case ButtonState::PRESSED:
          if (_input.get())
          {
            if (_button_held_ts < _ts.get())
            {
              _state = ButtonState::HELD;
            }
          }
          else if (_button_pressed_ts < _ts.get())
          {
            _state = ButtonState::CLICKED;
          }
          else
          {
            _state = ButtonState::RELEASED;
          }
          break;
        case ButtonState::CLICKED:
          this->publish(event::EventButtonClicked(_input.id()));
          _state = ButtonState::RELEASED;
          break;
        case ButtonState::HELD:
          this->publish(event::EventButtonHeld(_input.id()));
          _state = ButtonState::HELD_CONTINUE;

          break;
        case ButtonState::HELD_CONTINUE:
          if (!_input.get())
          {
            _state = ButtonState::RELEASED;
          }
          break;
        default:
          _state = ButtonState::RELEASED;
          break;
        }
      }

    private:
      const ISensor<bool> &_input;
      const core::ITimestamp &_ts;
      const size_t _button_click_ms;
      const size_t _button_hold_ms;
      ButtonState _state;
      size_t _button_pressed_ts;
      size_t _button_held_ts;
    };
  }
}