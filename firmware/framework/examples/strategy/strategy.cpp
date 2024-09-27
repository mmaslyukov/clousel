#include "framework/core/logger.h"
#include "framework/strategy.h"
#include "framework/core/i_timestamp.h"
#include <stdint.h>
#include <chrono>

class LoggerSystem : public core::logger::ILoggerSystem, public core::ITimestamp
{
public:
  virtual size_t get() const override { return 0; };
  virtual void output(const core::logger::Verbosity &verbosity, const char *tag, const char *data, size_t size) const override
  {
    printf("%d %s\n", verbosity.id(), data);
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
class MyGpio : public core::io::IActuator<bool>
{
public:
  MyGpio(uint32_t id, const core::logger::ILogger &logger) : _id(id), _logger(logger) {}
  virtual uint32_t id() const override
  {
    return _id;
  }
  virtual void set(const bool &value) override
  {
    _logger.dbg().log("gpio", "GPIO%d %s", _id, value ? "on" : "off");
  }
  virtual bool get() const override 
  {
    return true;
  }

private:
  uint32_t _id;
  const core::logger::ILogger &_logger;
};

int main()
{
  using namespace core::logger;
  using namespace strategy;

  Timestamp ts;
  char buff[128];
  PrintableNone none;
  const Logger logger(
      none,
      none,
      none,
      Printable(Configuration(buff, sizeof(buff), Verbosity("D")), LoggerSystem(), ts, true),
      DumpableNone());
  MyGpio gpio(1, logger);
  StepWait sw100(500, ts);
  StepActuatorOff saof(gpio);
  StepActuatorOn saon(gpio);
  IStep *st[] = {
      &saon,
      &sw100,
      &saof,
  };
  StepRunner sr(st, sizeof(st) / sizeof(st[0]));

  while (!sr.finished())
  {
    sr.run();
  }
  logger.dbg().log("main", "Done");
  return 0;
}