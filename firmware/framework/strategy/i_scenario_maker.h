#pragma once
#include <framework/strategy/step_runner.h>
namespace strategy
{
  struct IScenarioMaker
  {
    virtual strategy::StepRunner &make() const = 0;
  };
}
