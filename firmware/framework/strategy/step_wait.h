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
        : _ts(timeapi), _tm(timeout_ms)
    {
    }
    virtual void prepare() override
    {
      _target_timestamp = _ts.get() + _tm;
    }
    virtual void complete() override
    {
    }
    virtual bool execute() override
    {
      return _target_timestamp <= _ts.get();
    }

  private:
    // const uint32_t _tm;
    const core::ITimestamp &_ts;
    size_t _tm;
    size_t _target_timestamp;
  };
}