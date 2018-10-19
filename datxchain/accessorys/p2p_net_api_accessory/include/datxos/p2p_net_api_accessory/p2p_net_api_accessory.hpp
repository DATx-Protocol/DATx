/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#pragma once

#include <datxos/p2p_net_accessory/p2p_net_accessory.hpp>
#include <datxos/http_accessory/http_accessory.hpp>

#include <appbase/application.hpp>

namespace datxos {

using namespace appbase;

class p2p_net_api_accessory : public accessory<p2p_net_api_accessory> {
public:
   APPBASE_accessory_REQUIRES((p2p_net_accessory) (http_accessory))

   p2p_net_api_accessory() = default;
   p2p_net_api_accessory(const p2p_net_api_accessory&) = delete;
   p2p_net_api_accessory(p2p_net_api_accessory&&) = delete;
   p2p_net_api_accessory& operator=(const p2p_net_api_accessory&) = delete;
   p2p_net_api_accessory& operator=(p2p_net_api_accessory&&) = delete;
   virtual ~p2p_net_api_accessory() override = default;

   virtual void set_program_options(options_description& cli, options_description& cfg) override {}
   void accessory_initialize(const variables_map& vm);
   void accessory_startup();
   void accessory_shutdown() {}

private:
};

}
