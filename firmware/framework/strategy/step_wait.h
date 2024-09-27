#pragma once
#include "../core/logger/i_logger.h"
#include "../core/i_timestamp.h"
#include "i_step.h"

namespace strategy
{
  class StepWait : public IStep
  {
  public:
    StepWait(uint32_t timeout_ms, const core::ITimestamp &timeapi)
        : _target_timestamp(timeapi.get() + timeout_ms), _timeapi(timeapi)
    {
    }
    virtual bool execute() override
    {
      return _target_timestamp <= _timeapi.get();
    }

  private:
    // const uint32_t _tm;
    const size_t _target_timestamp;
    const core::ITimestamp &_timeapi;
  };
}