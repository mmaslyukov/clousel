#pragma once
#include "../core/io/i_actuator.h"
#include "i_step.h"

namespace strategy
{
  class StepActuatorOn : public IStep
  {
  public:
    StepActuatorOn(core::io::IActuator<bool>& out) : _out(out) {}
    virtual bool execute() override
    {
      _out.set(true);
      return true;
    }

  private:
    core::io::IActuator<bool>& _out;
  };
}