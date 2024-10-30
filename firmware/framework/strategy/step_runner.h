#pragma once
#include "../core/i_runnable.h"
#include "i_step.h"
#include <stdint.h>

namespace strategy
{
  class StepRunner : public core::IRunnable
  {
  public:
    StepRunner() : _steps(nullptr), _size(0), _index(0), _exec(nullptr) {}
    StepRunner(IStep **steps, size_t size)
        : _steps(steps), _size(size), _index(_size), _exec(nullptr) {}

    virtual void run() override
    {
      do
      {
        if (!_steps)
        {
          break;
        }

        if (_exec)
        {
          if (_exec->execute())
          {
            _exec = nullptr;
            _steps[_index]->complete();
            _index++;
          }
          else
          {
            // step is still being executed
            break;
          }
        }

        if (_index >= _size || !_steps[_index])
        {
          break;
        }
        _steps[_index]->prepare();
        _exec = _steps[_index];

      } while (true);
    }

    void apply(IStep **steps, size_t size)
    {
      _index = 0;
      _steps = steps;
      _size = size;
    }

    void reset()
    {
      _index = 0;
    }

    bool finished()
    {
      return _index >= _size || !_steps[_index];
    }

  private:
    IStep **_steps;
    size_t _size;
    size_t _index;
    IExecute *_exec;
  };
}