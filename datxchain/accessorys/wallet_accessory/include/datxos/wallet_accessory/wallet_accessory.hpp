/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#pragma once
#include <appbase/application.hpp>
#include <fc/variant.hpp>
#include <datxos/chain/contract_types.hpp>
#include <datxos/chain/transaction.hpp>

namespace fc { class variant; }

namespace datxos {
   using namespace appbase;

   namespace wallet {
      class wallet_manager;
   }
   using namespace wallet;

class wallet_accessory : public accessory<wallet_accessory> {
public:
   APPBASE_accessory_REQUIRES()

   wallet_accessory();
   wallet_accessory(const wallet_accessory&) = delete;
   wallet_accessory(wallet_accessory&&) = delete;
   wallet_accessory& operator=(const wallet_accessory&) = delete;
   wallet_accessory& operator=(wallet_accessory&&) = delete;
   virtual ~wallet_accessory() override = default;

   virtual void set_program_options(options_description& cli, options_description& cfg) override;
   void accessory_initialize(const variables_map& options);
   void accessory_startup() {}
   void accessory_shutdown() {}

   // api interface provider
   wallet_manager& get_wallet_manager();

private:
   std::unique_ptr<wallet_manager> wallet_manager_ptr;
};

}
