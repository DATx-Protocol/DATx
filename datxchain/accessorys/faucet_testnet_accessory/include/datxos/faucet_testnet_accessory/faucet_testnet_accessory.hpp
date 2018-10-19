/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#pragma once
#include <appbase/application.hpp>
#include <datxos/http_accessory/http_accessory.hpp>

namespace datxos {

using namespace appbase;

class faucet_testp2p_net_accessory : public appbase::accessory<faucet_testp2p_net_accessory> {
public:
   faucet_testp2p_net_accessory();
   ~faucet_testp2p_net_accessory();

   APPBASE_accessory_REQUIRES((http_accessory))
   virtual void set_program_options(options_description&, options_description& cfg) override;
 
   void accessory_initialize(const variables_map& options);
   void accessory_startup();
   void accessory_shutdown();

private:
   std::unique_ptr<struct faucet_testp2p_net_accessory_impl> my;
};

}
