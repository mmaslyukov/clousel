#pragma once

// #include "esp_https_server.h"
#include "esp_http_server.h"

struct IHttpEvetHandler
{
  virtual esp_err_t handle(httpd_req_t *req) = 0;
};