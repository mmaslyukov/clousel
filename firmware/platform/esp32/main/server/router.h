#pragma once
#include <framework/core/logger.h>
#include <service/web/port_web_controller_api.h>

#include "../event/i_http_event_handler.h"

class HttpReqestRoot : public IHttpEvetHandler
{
public:
  HttpReqestRoot(const core::logger::ILogger &logger) : _logger(logger) {}
  virtual esp_err_t handle(httpd_req_t *req) override
  {
    static const char *page = R"(<!DOCTYPE html><html><head> <meta name="viewport" content="width=1280, initial-scale=1"> <link rel="icon" href="data:,"><style> h1 { text-align: center; font-family: "Impact"; margin-top: 4rem; margin-bottom: 4rem; font-size: 60px; } .button-div { margin-top: 3em; text-align: center; } .button { font-family: "Impact"; font-size: 30px; min-width: 7rem; padding: 10px 30px; } .button:active { background-color: rgb(152, 152, 152); } input[type=text] { border: none; border-bottom: 2px solid; font-size: 25px; } .entry { font-family: "Impact"; text-align: center; margin-top: 1em; margin-right: 1em; } .lable { font-size: 25px; display: inline-block; text-align: right; min-width: 5em; } .field { display: inline-block; text-align: left; max-width: 8em; } .border { border-style: solid; } </style></head><script> var siteWidth = 1280; var scale = 1; console.log("screen.width:", screen.width, " scale:", scale); document.querySelector('meta[name="viewport"]').setAttribute('content', 'width=' + screen.width + ', initial-scale=' + scale + '');</script><body> <h1>CLOUSEL</h1> <form action="/submit"> <div class="entry"> <div class="lable"> <label class="label" for="ssid">SSID:</label> </div> <input class="field" type="text" id="ssid" name="ssid"><br><br> </div> <div class="entry"> <div class="lable"> <label for="pswd">Password:</label> </div> <input class="field" type="text" id="pswd" name="pswd"><br><br> </div> <div class="button-div"> <input class="button" type="submit" value="Save"> </div> </form></body></html>)";
    httpd_resp_set_type(req, "text/html");
    httpd_resp_send(req, page, HTTPD_RESP_USE_STRLEN);
    return ESP_OK;
  }

private:
  const core::logger::ILogger &_logger;
  static constexpr const char *TAG = "router.root";
};

class HttpReqestSubmit : public IHttpEvetHandler
{
public:
  HttpReqestSubmit(const core::logger::ILogger &logger, service::web::IPortWebControllerApi &api) : _logger(logger), _api(api) {}
  virtual esp_err_t handle(httpd_req_t *req) override
  {
    do
    {
      if (httpd_req_get_url_query_len(req) == 0)
      {
        _logger.err().log(TAG, "Fail to get query length");
        break;
      }
      char query_buffer[100];
      if (httpd_req_get_url_query_str(req, query_buffer, sizeof(query_buffer)) != ESP_OK)
      {
        _logger.err().log(TAG, "Fail copy query to a buffer");
        break;
      }
      char query_ssid[32];// = {0}; // 0th byte is 0 (\0), meaning zero size string
      memset(query_ssid, 0, sizeof(query_ssid));
      if (httpd_query_key_value(query_buffer, "ssid", query_ssid, sizeof(query_ssid)) != ESP_OK)
      {
        _logger.err().log(TAG, "Fail find 'ssid' in query:'%s'", query_buffer);
        break;
      }
      char query_pswd[32];// = {0}; // 0th byte is 0 (\0), meaning zero size string
      memset(query_pswd, 0, sizeof(query_pswd));
      if (httpd_query_key_value(query_buffer, "pswd", query_pswd, sizeof(query_pswd)) != ESP_OK)
      {
        _logger.err().log(TAG, "Fail find 'pswd' in query:'%s'", query_buffer);
        break;
      }
      bool ok = false;
      if (strlen(query_ssid) > 0 && strlen(query_pswd) > 0)
      {
        if (_api.submit(infra::WifiSettingsN(query_ssid, query_pswd)))
        {
          ok = true;
        }
        else
        {
          _logger.err().log(TAG, "Fail to submit new Wifi credentials");
        }
      }
      else
      {
        _logger.err().log(TAG, "Fail to apply query data");
      }
      if (ok)
      {
        const char *page_ok = R"(<!DOCTYPE html><html><head><meta name="viewport" content="width=1280,initial-scale=1"><style>h1{text-align:center;font-family:Impact;margin-top:4rem;margin-bottom:4rem;font-size:60px}p{text-align:center;font-size:20px}</style></head><script>var siteWidth=1280,scale=1;console.log("screen.width:",screen.width," scale:",scale),document.querySelector('meta[name="viewport"]').setAttribute("content","width="+screen.width+", initial-scale="+scale)</script><body><h1>CLOUSEL</h1><p>Settings has been saved successfully</p></body></html>)";
        httpd_resp_set_type(req, "text/html");
        httpd_resp_send(req, page_ok, HTTPD_RESP_USE_STRLEN);
      }
      else
      {
        const char *page_fail = R"(<!DOCTYPE html><html><head><meta name="viewport" content="width=1280,initial-scale=1"><style>h1{text-align:center;font-family:Impact;margin-top:4rem;margin-bottom:4rem;font-size:60px}p{text-align:center;font-size:20px}</style></head><script>var siteWidth=1280,scale=1;console.log("screen.width:",screen.width," scale:",scale),document.querySelector('meta[name="viewport"]').setAttribute("content","width="+screen.width+", initial-scale="+scale),window.setTimeout(function(){window.location="/"},3e3)</script><body><h1>CLOUSEL</h1><p>Fail to save ssid and password</p><p>Please try again</p></body></html>)";
        httpd_resp_set_type(req, "text/html");
        httpd_resp_send(req, page_fail, HTTPD_RESP_USE_STRLEN);
      }

    } while (false);

    // for (auto &p : req.params)
    // {
    //   if (!strcmp(p.first.c_str(), "pswd"))
    //   {
    //     pswd = p.second.c_str();
    //   }
    //   else if (!strcmp(p.first.c_str(), "ssid"))
    //   {
    //     ssid = p.second.c_str();
    //   }

    //   _logger.dbg().log(TAG, "query %s=%s", p.first.c_str(), p.second.c_str());
    // }

    // if (ssid && pswd)
    // {
    //   if (_api.submit(infra::WifiSettingsN(ssid, pswd)))
    //   {
    //     ok = true;
    //   }
    //   else
    //   {
    //     _logger.err().log(TAG, "Fail to submit new Wifi credentials");
    //   }
    // }
    // else
    // {
    //   _logger.err().log(TAG, "Fail to apply query data");
    // }
    // if (ok)
    // {
    //   const char *page_ok = R"(<!DOCTYPE html><html><head><meta name="viewport" content="width=1280,initial-scale=1"><style>h1{text-align:center;font-family:Impact;margin-top:4rem;margin-bottom:4rem;font-size:60px}p{text-align:center;font-size:20px}</style></head><script>var siteWidth=1280,scale=1;console.log("screen.width:",screen.width," scale:",scale),document.querySelector('meta[name="viewport"]').setAttribute("content","width="+screen.width+", initial-scale="+scale)</script><body><h1>CLOUSEL</h1><p>Settings has been saved successfully</p></body></html>)";
    //   res.set_content(page_ok, "text/html");
    // }
    // else
    // {
    //   const char *page_fail = R"(<!DOCTYPE html><html><head><meta name="viewport" content="width=1280,initial-scale=1"><style>h1{text-align:center;font-family:Impact;margin-top:4rem;margin-bottom:4rem;font-size:60px}p{text-align:center;font-size:20px}</style></head><script>var siteWidth=1280,scale=1;console.log("screen.width:",screen.width," scale:",scale),document.querySelector('meta[name="viewport"]').setAttribute("content","width="+screen.width+", initial-scale="+scale),window.setTimeout(function(){window.location="/"},3e3)</script><body><h1>CLOUSEL</h1><p>Fail to save ssid and password</p><p>Please try again</p></body></html>)";
    //   // const char *page_fail = R"(<!DOCTYPE html><html><head><meta name="viewport" content="width=1280,initial-scale=1"><style>h1{text-align:center;font-family:Impact;margin-top:4rem;margin-bottom:4rem;font-size:60px}p{text-align:center;font-size:20px}</style></head><script>var siteWidth=1280,scale=1;console.log("screen.width:",screen.width," scale:",scale),document.querySelector('meta[name="viewport"]').setAttribute("content","width="+screen.width+", initial-scale="+scale)</script><body><h1>CLOUSEL</h1><p>Fail to save ssid and password</p><p>Please try again</p></body></html>)";
    //   res.set_content(page_fail, "text/html");
    // }
    return ESP_OK;
  }

private:
  const core::logger::ILogger &_logger;
  service::web::IPortWebControllerApi &_api;
  static constexpr const char *TAG = "router.submit";
};

class Router
{
public:
  Router(const core::logger::ILogger &logger, service::web::IPortWebControllerApi &api)
      : _logger(logger),
        _api(api),
        _handler_root(logger),
        _handler_submit(logger, api),
        _uris{{"/", HTTP_GET, http_route_handler, &_handler_root},
              {"/submit", HTTP_GET, http_route_handler, &_handler_submit}}
  {
    ESP_ERROR_CHECK(esp_event_handler_instance_register(WIFI_EVENT,
                                                        ESP_EVENT_ANY_ID,
                                                        &router_event_handler,
                                                        this,
                                                        &_instance_any_id));
  }

private:
  static esp_err_t http_route_handler(httpd_req_t *req)
  {
    if (req && req->user_ctx)
    {
      return reinterpret_cast<IHttpEvetHandler *>(req->user_ctx)->handle(req);
    }
    return false;
  }

  void start()
  {
    httpd_config_t conf = HTTPD_DEFAULT_CONFIG();

    esp_err_t err = httpd_start(&_server, &conf);
    if (err != ESP_OK)
    {
      _logger.err().log(TAG, "Error starting server!, error:%d", err);
      return;
    }

    for (uint8_t i = 0; i < sizeof(_uris) / sizeof(_uris[0]); i++)
    {
      esp_err_t err = httpd_register_uri_handler(_server, &_uris[i]);
      if (err != ESP_OK)
      {
        _logger.err().log(TAG, "Fail to register http route: '%s', error:%d", _uris[i].uri, err);
      }
      else
      {
        _logger.inf().log(TAG, "Registered http route: '%s'", _uris[i].uri);
      }
    }
  }

  void stop()
  {
    esp_err_t err = httpd_stop(_server);
    if (err != ESP_OK)
    {
      _logger.err().log(TAG, "Fail to stop http server, error:%d", err);
    }
  }

  static void router_event_handler(void *arg, esp_event_base_t event_base, int32_t event_id, void *event_data)
  {
    Router *self = reinterpret_cast<Router *>(arg);

    if (event_base == WIFI_EVENT)
    {
      switch (event_id)
      {
      case WIFI_EVENT_AP_START:
      {
        self->start();
        // for (uint8_t i = 0; i < sizeof(_uris) / sizeof(_uris[0]); i++)
        // {
        //   esp_err_t err = httpd_register_uri_handler(self->_server, &self->_uris[i]);
        //   if (err != ESP_OK)
        //   {
        //     self->_logger.err().log(self->TAG, "Fail to register http route: '%s', error:%d", self->_uris[i].uri, err);
        //   }
        //   else
        //   {
        //     self->_logger.inf().log(self->TAG, "Registered http route: '%s'", self->_uris[i].uri);
        //   }
        // }
      }
      break;
      case WIFI_EVENT_AP_STOP:
      {
        self->stop();
        // esp_err_t err = httpd_stop(self->_server);
        // if (err != ESP_OK)
        // {
        //   self->_logger.err().log(self->TAG, "Fail to stop http router, error:%d", err);
        // }
      }
      break;
      default:
        break;
      }
    }
  }

private:
  const core::logger::ILogger &_logger;
  service::web::IPortWebControllerApi &_api;
  HttpReqestRoot _handler_root;
  HttpReqestSubmit _handler_submit;
  const httpd_uri_t _uris[2];
  httpd_handle_t _server;
  esp_event_handler_instance_t _instance_any_id;
  static constexpr const char *TAG = "router";
};