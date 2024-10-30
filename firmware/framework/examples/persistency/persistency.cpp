/**
 * Build:
 * cmake --preset=default
 * cmake --build build
 */

#include "framework/core/logger.h"
#include "framework/persistency.h"
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

struct Flash : public persistency::IPersistencyFlash
{
  virtual bool load(uint8_t *memory, size_t size) const
  {
    return true;
  };
  virtual bool save(const uint8_t *memory, size_t size) const
  {
    return true;
  };
};

#pragma pack(push, 1)
struct CfgOne
{
  uint8_t u8;
  uint16_t u16;
  uint32_t u32;
};
struct CfgTwo
{
  uint8_t u8;
  uint16_t u16;
  uint32_t u32;
};
#pragma pack(pop)

enum PersistencyId : uint8_t
{
  DUMMY_CFG_ONE,
  DUMMY_CFG_TWO,
  // Add new item id here for addressing data in persistency
  _LAST,
};
int main()
{
  using namespace core::logger;
  Timestamp ts;
  LoggerSystem ls;
  char buff[128];
  PrintableNone none;
  const Logger logger(
      none,
      none,
      none,
      Printable(Configuration(buff, sizeof(buff), Verbosity("D")), ls, ts, true),
      none,
      Dumpable(Configuration(buff, sizeof(buff), Verbosity("D")), ls, ts, true));

  const size_t PERSISTENCE_SIZE = 256;
  static uint8_t memory[PERSISTENCE_SIZE];

  static constexpr persistency::Persistency<PersistencyId>::Row persistency_table[] = {
      persistency::Persistency<PersistencyId>::make_persistency_row<CfgOne>(PersistencyId::DUMMY_CFG_ONE),
      persistency::Persistency<PersistencyId>::make_persistency_row<CfgTwo>(PersistencyId::DUMMY_CFG_TWO),
      {PersistencyId::_LAST, 0, nullptr},
  };

  static_assert(
      persistency::Persistency<PersistencyId>::check_persistency_table_size(
          persistency_table,
          PERSISTENCE_SIZE),
      "Error, insufficient memory for persistency table mapping");

  persistency::Persistency<PersistencyId> storage(persistency_table, memory, sizeof(memory), Flash());
  CfgOne c1;
  storage.read(PersistencyId::DUMMY_CFG_ONE, &c1);
  c1.u8 = 0x12;
  c1.u16 = 0x1234;
  c1.u32 = 0x12345678;
  storage.write(PersistencyId::DUMMY_CFG_ONE, &c1);
  logger.raw().dump("CfgOne", (const uint8_t *)storage.get<CfgOne>(PersistencyId::DUMMY_CFG_ONE), sizeof(CfgOne));
}