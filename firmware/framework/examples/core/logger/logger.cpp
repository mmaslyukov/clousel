/**
 * Build:
 * cmake --preset=default
 * cmake --build build
 */


#include "framework/core/logger.h"
#include "framework/core/i_timestamp.h"
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



int main()
{
  using namespace core::logger;

  Timestamp ts;
  LoggerSystem ls;
  char buff[128];
  const Logger logger(
      Printable(Configuration(buff, sizeof(buff), Verbosity("E")), ls, ts, true),
      Printable(Configuration(buff, sizeof(buff), Verbosity("W")), ls, ts, true),
      Printable(Configuration(buff, sizeof(buff), Verbosity("I")), ls, ts, true),
      Printable(Configuration(buff, sizeof(buff), Verbosity("D")), ls, ts, true),
      Printable(Configuration(buff, sizeof(buff), Verbosity("V")), ls, ts, true),
      Dumpable(Configuration(buff, sizeof(buff), Verbosity("D")), ls, ts , true));
  logger.raw().dump("buf", (uint8_t*)buff, 210);
  logger.err().log("test", "Hello Johny %d times", 5);
  logger.dbg().log("test", "Hello Johny %d times", 5);
}