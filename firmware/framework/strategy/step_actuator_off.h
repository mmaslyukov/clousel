#pragma once
#include "../core/io/i_actuator.h"
#include "i_step.h"

namespace strategy
{
  class StepActuatorOff : public IStep
  {
  public:
    StepActuatorOff(core::io::IActuator<bool>& out) : _out(out) {}
    virtual bool execute() override
    {
      _out.set(false);
      return true;
    }

  private:
    core::io::IActuator<bool>& _out;
  };
}