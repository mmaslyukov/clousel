#pragma once
#include "httplib.h"
#include <service/web/port_web_controller_api.h>
#include <framework/core/logger.h>
namespace platform
{
  namespace router
  {

    using namespace httplib;
    class Router
    {
    public:
      constexpr Router(Server &server,
                       service::web::IPortWebControllerApi &api,
                       const core::logger::ILogger &logger)
          : _server(server), _api(api), _logger(logger) {}
      Router &route()
      {
        _server.Get("/",
                    [&](const Request & /*req*/, Response &res)
                    {
                      const char *page = R"(<!DOCTYPE html><html><head> <meta name="viewport" content="width=1280, initial-scale=1"> <link rel="icon" href="data:,"><style> h1 { text-align: center; font-family: "Impact"; margin-top: 4rem; margin-bottom: 4rem; font-size: 60px; } .button-div { margin-top: 3em; text-align: center; } .button { font-family: "Impact"; font-size: 30px; min-width: 7rem; padding: 10px 30px; } .button:active { background-color: rgb(152, 152, 152); } input[type=text] { border: none; border-bottom: 2px solid; font-size: 25px; } .entry { font-family: "Impact"; text-align: center; margin-top: 1em; margin-right: 1em; } .lable { font-size: 25px; display: inline-block; text-align: right; min-width: 5em; } .field { display: inline-block; text-align: left; max-width: 8em; } .border { border-style: solid; } </style></head><script> var siteWidth = 1280; var scale = 1; console.log("screen.width:", screen.width, " scale:", scale); document.querySelector('meta[name="viewport"]').setAttribute('content', 'width=' + screen.width + ', initial-scale=' + scale + '');</script><body> <h1>CLOUSEL</h1> <form action="/submit"> <div class="entry"> <div class="lable"> <label class="label" for="ssid">SSID:</label> </div> <input class="field" type="text" id="ssid" name="ssid"><br><br> </div> <div class="entry"> <div class="lable"> <label for="pswd">Password:</label> </div> <input class="field" type="text" id="pswd" name="pswd"><br><br> </div> <div class="button-div"> <input class="button" type="submit" value="Save"> </div> </form></body></html>)";
                      res.set_content(page, "text/html");
                    });

        _server.Get("/submit",
                    [&](const Request &req, Response &res)
                    {
                      bool ok = false;
                      const char *ssid = nullptr;
                      const char *pswd = nullptr;
                      for (auto &p : req.params)
                      {
                        if (!strcmp(p.first.c_str(), "pswd"))
                        {
                          pswd = p.second.c_str();
                        }
                        else if (!strcmp(p.first.c_str(), "ssid"))
                        {
                          ssid = p.second.c_str();
                        }

                        _logger.dbg().log(TAG, "query %s=%s", p.first.c_str(), p.second.c_str());
                      }

                      if (ssid && pswd)
                      {
                        if (_api.submit(infra::WifiSettingsN(ssid, pswd)))
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
                        res.set_content(page_ok, "text/html");
                      }
                      else
                      {
                        const char *page_fail = R"(<!DOCTYPE html><html><head><meta name="viewport" content="width=1280,initial-scale=1"><style>h1{text-align:center;font-family:Impact;margin-top:4rem;margin-bottom:4rem;font-size:60px}p{text-align:center;font-size:20px}</style></head><script>var siteWidth=1280,scale=1;console.log("screen.width:",screen.width," scale:",scale),document.querySelector('meta[name="viewport"]').setAttribute("content","width="+screen.width+", initial-scale="+scale),window.setTimeout(function(){window.location="/"},3e3)</script><body><h1>CLOUSEL</h1><p>Fail to save ssid and password</p><p>Please try again</p></body></html>)";
                        // const char *page_fail = R"(<!DOCTYPE html><html><head><meta name="viewport" content="width=1280,initial-scale=1"><style>h1{text-align:center;font-family:Impact;margin-top:4rem;margin-bottom:4rem;font-size:60px}p{text-align:center;font-size:20px}</style></head><script>var siteWidth=1280,scale=1;console.log("screen.width:",screen.width," scale:",scale),document.querySelector('meta[name="viewport"]').setAttribute("content","width="+screen.width+", initial-scale="+scale)</script><body><h1>CLOUSEL</h1><p>Fail to save ssid and password</p><p>Please try again</p></body></html>)";
                        res.set_content(page_fail, "text/html");
                      }

                      // validate
                    });
        return *this;
      }
      void listen(const char *host, const size_t port)
      {
        _server.listen(host, port);
      }

    private:
      httplib::Server &_server;
      service::web::IPortWebControllerApi &_api;
      const core::logger::ILogger &_logger;
      static constexpr const char *TAG = "router";
    };

  } // namespace server

} // namespace platform
