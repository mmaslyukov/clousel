#include "framework/core/logger.h"
#include "framework/core/i_timestamp.h"
#include "framework/core/observer.h"
#include <stdint.h>
#include <chrono>

class LoggerSystem : public core::logger::ILoggerSystem, public core::ITimestamp
{
public:
  virtual size_t get() const override { return 0; };
  virtual void output(const core::logger::Verbosity &verbosity, size_t tsms, const char *tag, const char *data, size_t size) const override
  {
    printf("%s (%zu) <%s> %s\n", verbosity.name(), tsms, tag, data);
  };
};

class Timestamp : public core::ITimestamp
{
public:
  virtual size_t get() const override
  {
    return std::chrono::duration_cast<std::chrono::milliseconds>(std::chrono::system_clock::now().time_since_epoch()).count();
  }
};

struct TestListenerOne : public core::observer::IListener
{
  TestListenerOne(const core::logger::ILogger &logger) : _logger(logger) {}
  virtual void notify(const core::observer::Event &event) override
  {
    _logger.dbg().log("TestListenerOne", "Get an event %s, id %d", event.name(), event.id());
  }

private:
  const core::logger::ILogger &_logger;
};

struct TestListenerTwo : public core::observer::IListener
{
  TestListenerTwo(const core::logger::ILogger &logger) : _logger(logger) {}
  virtual void notify(const core::observer::Event &event) override
  {
    _logger.dbg().log("TestListenerTwo", "Get an event %s, id %d", event.name(), event.id());
  }

private:
  const core::logger::ILogger &_logger;
};

struct MyCustomEvent : public core::observer::Event
{
  static inline const char *event_name() { return "my.custom.event"; }
  MyCustomEvent() : core::observer::Event(event_name()) {}
};

// class MyCustomEvent : public core::observer::Event
// {
//   static const char *_name() { return "my.custom.event"; }

// public:
//   static core::observer::EventBase base() { return EventBase(_name()); }
//   MyCustomEvent() : Event(_name()) {}
// };

struct TestSubscriber : public core::observer::Publisher<2>
{
  TestSubscriber(const core::logger::ILogger &logger) : _logger(logger) {}
  void action()
  {
    MyCustomEvent ev;
    _logger.dbg().log("TestSubscriber", "About to send event %s", ev.name());
    publish(ev);
  }

private:
  const core::logger::ILogger &_logger;
};

int main()
{
  using namespace core::logger;

  char buff[128];
  PrintableNone none;
  const Logger logger(
      none,
      none,
      none,
      Printable(Configuration(buff, sizeof(buff), Verbosity("D")), LoggerSystem(), Timestamp(), true),
      none,
      DumpableNone());
  TestListenerOne tl1(logger);
  TestListenerTwo tl2(logger);
  TestSubscriber tsb(logger);
  tsb.add_subscriber(&tl1, MyCustomEvent::event_name());
  tsb.add_subscriber(&tl2, MyCustomEvent::event_name());
  tsb.action();
}