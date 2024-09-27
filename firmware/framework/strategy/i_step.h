#pragma once

namespace strategy
{
  struct IStep
  {
    virtual bool execute() = 0;
  };
}
