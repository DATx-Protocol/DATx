/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#pragma once

#include <datxos/core_accessory/core_accessory.hpp>
#include <appbase/application.hpp>
#include <memory>

namespace datxos {

using mongo_db_accessory_impl_ptr = std::shared_ptr<class mongo_db_accessory_impl>;

/**
 * Provides persistence to MongoDB for:
 * accounts
 * actions
 * block_states
 * blocks
 * transaction_traces
 * transactions
 * pub_keys
 * account_controls
 *
 *   See data dictionary (DB Schema Definition - DATX API) for description of MongoDB schema.
 *
 *   If cmake -DBUILD_MONGO_DB_accessory=true  not specified then this accessory not compiled/included.
 */
class mongo_db_accessory : public accessory<mongo_db_accessory> {
public:
   APPBASE_accessory_REQUIRES((core_accessory))

   mongo_db_accessory();
   virtual ~mongo_db_accessory();

   virtual void set_program_options(options_description& cli, options_description& cfg) override;

   void accessory_initialize(const variables_map& options);
   void accessory_startup();
   void accessory_shutdown();

private:
   mongo_db_accessory_impl_ptr my;
};

}

