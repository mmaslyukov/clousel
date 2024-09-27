#pragma once
#include <stdint.h>
#include <cstring>
#include <string>
#include "i_persistency_flash.h"

namespace persistency
{

  
  constexpr const size_t max(size_t a, size_t b)
  {
    return a > b ? a : b;
  }

  template <typename Id>
  class Persistency
  {
  public:
    struct Row
    {
      Id id;
      uint16_t size;
      void (*const reset_default)(uint8_t *, size_t);
    };

    constexpr Persistency(
        const Row *const table,
        uint8_t *memory,
        const size_t size,
        const IPersistencyFlash &flash)
        : _table(table), _memory(memory), _size(size), _flash(flash)
    {
    }

    template <typename T>
    static constexpr Row make_persistency_row(Id id)
    {
      Row row{id, sizeof(T), apply_default<T>};
      return row;
    }

    static constexpr bool check_persistency_table_size(
        const Persistency::Row *table_data,
        const size_t table_size_bytes)
    {
      size_t size = 0;
      for (const Persistency::Row *row_ptr = table_data;
           row_ptr->id != Id::_LAST;
           row_ptr++)
      {
        size += row_ptr->size;
      }
      return size < table_size_bytes;
    }

    bool load() const
    {
      return _flash.load(_memory, _size);
    }

    bool save() const
    {
      return _flash.save(_memory, _size);
    }

    template <typename T>
    bool read(const Id id, T *value) const
    {
      auto row_entry = find_row(id);
      if (value &&
          row_entry.row &&
          row_entry.row->size == sizeof(T) &&
          ((row_entry.index + sizeof(T)) < _size))
      {
        memcpy(
            reinterpret_cast<T *>(value),
            &_memory[row_entry.index],
            sizeof(T));
        return true;
      }
      return false;
    }

    template <typename T>
    bool write(const Id id, const T *value)
    {
      auto row_entry = find_row(id);
      if (value &&
          row_entry.row &&
          row_entry.row->size == sizeof(T) &&
          ((row_entry.index + sizeof(T)) < _size))
      {
        memcpy(
            &_memory[row_entry.index],
            reinterpret_cast<const T *>(value),
            sizeof(T));
        return true;
      }
      return false;
    }

    template <typename T>
    const T *const get(Id id) const
    {
      auto row_entry = find_row(id);
      if (row_entry.row &&
          row_entry.row->size == sizeof(T) &&
          ((row_entry.index + sizeof(T)) < _size))
      {
        return reinterpret_cast<T *>(&_memory[row_entry.index]);
      }
      else
      {
        return nullptr;
      }
    }

    bool reset_default(Id id)
    {
      auto row_entry = find_row(id);
      if (row_entry.row && row_entry.row->reset_default)
      {
        row_entry.row->reset_default(
            &_memory[row_entry.index],
            row_entry.row->size);
        return true;
      }
      return false;
    }

    void reset_default_all()
    {
      size_t index = 0;
      for (const Row *row_ptr = _table;
           row_ptr->id != Id::_LAST;
           index += row_ptr->size, row_ptr++)
      {
        if (row_ptr->reset_default)
        {
          row_ptr->reset_default(&_memory[index], row_ptr->size);
        }
      }
    }

  private:
    template <typename T>
    static constexpr void apply_default(uint8_t *mem, size_t size)
    {
      if (mem)
      {
        T t;
        memcpy(mem, reinterpret_cast<uint8_t *>(&t), size);
      }
    }
    struct Entry
    {
      size_t index;
      const Row *row;
    };

    const Entry find_row(Id id) const
    {
      Entry entry{0, nullptr};
      size_t index = 0;
      for (const Row *row_ptr = _table;
           row_ptr->id != Id::_LAST;
           index += row_ptr->size, row_ptr++)
      {
        if (id == row_ptr->id)
        {
          entry = Entry{index, row_ptr};
        }
      }
      return entry;
    }

    uint8_t *_memory;
    const size_t _size;
    const Row *const _table;
    const IPersistencyFlash &_flash;
  };
}
