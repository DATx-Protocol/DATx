/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#pragma once

#include <datxos/http_accessory/http_accessory.hpp>
#include <datxos/core_accessory/core_accessory.hpp>

#include <appbase/application.hpp>

namespace datxos {

using namespace appbase;

struct db_size_index_count {
   string   index;
   uint64_t row_count;
};

struct db_size_stats {
   uint64_t                    free_bytes;
   uint64_t                    used_bytes;
   uint64_t                    size;
   vector<db_size_index_count> indices;
};

class db_size_api_accessory : public accessory<db_size_api_accessory> {
public:
   APPBASE_accessory_REQUIRES((http_accessory) (core_accessory))

   db_size_api_accessory() = default;
   db_size_api_accessory(const db_size_api_accessory&) = delete;
   db_size_api_accessory(db_size_api_accessory&&) = delete;
   db_size_api_accessory& operator=(const db_size_api_accessory&) = delete;
   db_size_api_accessory& operator=(db_size_api_accessory&&) = delete;
   virtual ~db_size_api_accessory() override = default;

   virtual void set_program_options(options_description& cli, options_description& cfg) override {}
   void accessory_initialize(const variables_map& vm) {}
   void accessory_startup();
   void accessory_shutdown() {}

   db_size_stats get();

private:
};

}

FC_REFLECT( datxos::db_size_index_count, (index)(row_count) )
FC_REFLECT( datxos::db_size_stats, (free_bytes)(used_bytes)(size)(indices) )