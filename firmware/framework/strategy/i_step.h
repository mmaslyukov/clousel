#pragma once

namespace strategy
{

  struct IExecute
  {
    virtual bool execute() = 0;
  };

  struct IStep : public IExecute
  {
    virtual void prepare() = 0;
    virtual void complete() = 0;
  };
}
