#pragma once
#include "i_esp_event_handler.h"
#include "i_http_event_handler.h"

template <typename T, size_t N>
struct EvetProvider
{
  EvetProvider()
  {
    for (uint32_t i = 0; i < N; i++)
    {
      _handlers[i] = nullptr;
    }
  }

  void add_handler(T *handler)
  {
    for (uint32_t i = 0; i < N; i++)
    {
      if (_handlers[i] == nullptr)
      {
        _handlers[i] = handler;
        break;
      }
    }
  }

protected:
  T *_handlers[N];
};

template <size_t N>
struct EspEvetProvider : public EvetProvider<IEspEvetHandler, N>, public IEspEvetHandler
{
public:
  EspEvetProvider() : EvetProvider<IEspEvetHandler, N>() {}
  virtual bool handle(esp_event_base_t event_base, int32_t event_id, void *event_data) override
  {
    bool res = false;
    for (uint32_t i = 0; i < N; i++)
    {
      if (this->_handlers[i] != nullptr)
      {
        res = res || this->_handlers[i]->handle(event_base, event_id, event_data);
      }
    }
    return res;
  }
};

template <size_t N>
struct HttpEvetProvider : public EvetProvider<IHttpEvetHandler, N>, public IHttpEvetHandler
{
public:
  HttpEvetProvider() : EvetProvider<IHttpEvetHandler, N>() {}
  virtual bool handle(httpd_req_t *req) override
  {
    bool res = false;
    for (uint32_t i = 0; i < N; i++)
    {
      if (this->_handlers[i] != nullptr)
      {
        res = res || this->_handlers[i]->handle(req);
      }
    }
    return res;
  }
};
