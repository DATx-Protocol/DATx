/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#pragma once
#include <appbase/application.hpp>
#include <datxio/http_plugin/http_plugin.hpp>

namespace datxio {

using namespace appbase;

class faucet_testp2p_net_plugin : public appbase::plugin<faucet_testp2p_net_plugin> {
public:
   faucet_testp2p_net_plugin();
   ~faucet_testp2p_net_plugin();

   APPBASE_PLUGIN_REQUIRES((http_plugin))
   virtual void set_program_options(options_description&, options_description& cfg) override;
 
   void plugin_initialize(const variables_map& options);
   void plugin_startup();
   void plugin_shutdown();

private:
   std::unique_ptr<struct faucet_testp2p_net_plugin_impl> my;
};

}
