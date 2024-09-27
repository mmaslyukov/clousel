#pragma once
#include "../core/i_runnable.h"
#include "i_step.h"
#include <stdint.h>

namespace strategy
{
  class StepRunner : public core::IRunnable
  {
  public:
    StepRunner(IStep **steps, size_t size)
        : _steps(steps), _size(size), _index(_size) {}
    virtual void run() override
    {
      if (_steps && (_index < _size))
      {
        if (_steps[_index]->execute())
        {
          _index++;
        }
      }
    }
    void reset()
    {
      _index = 0;
    }
    bool finished()
    {
      return _index == _size;
    }

  private:
    IStep **_steps;
    size_t _size;
    size_t _index;
  };
}