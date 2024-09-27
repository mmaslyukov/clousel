#pragma once
namespace service
{
  namespace mode
  {
    struct IPortButtonController
    {
      virtual void clicked() = 0;
      virtual void pressed() = 0;
    };
  }
}