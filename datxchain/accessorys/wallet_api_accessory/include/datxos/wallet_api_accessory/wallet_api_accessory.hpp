/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#pragma once

#include <datxos/wallet_accessory/wallet_accessory.hpp>
#include <datxos/http_accessory/http_accessory.hpp>

#include <appbase/application.hpp>

namespace datxos {

using namespace appbase;

class wallet_api_accessory : public accessory<wallet_api_accessory> {
public:
   APPBASE_accessory_REQUIRES((wallet_accessory) (http_accessory))

   wallet_api_accessory() = default;
   wallet_api_accessory(const wallet_api_accessory&) = delete;
   wallet_api_accessory(wallet_api_accessory&&) = delete;
   wallet_api_accessory& operator=(const wallet_api_accessory&) = delete;
   wallet_api_accessory& operator=(wallet_api_accessory&&) = delete;
   virtual ~wallet_api_accessory() override = default;

   virtual void set_program_options(options_description& cli, options_description& cfg) override {}
   void accessory_initialize(const variables_map& vm);
   void accessory_startup();
   void accessory_shutdown() {}

private:
};

}
