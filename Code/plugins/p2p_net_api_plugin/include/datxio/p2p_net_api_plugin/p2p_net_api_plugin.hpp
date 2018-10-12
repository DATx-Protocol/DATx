/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#pragma once

#include <datxio/p2p_net_plugin/p2p_net_plugin.hpp>
#include <datxio/http_plugin/http_plugin.hpp>

#include <appbase/application.hpp>

namespace datxio {

using namespace appbase;

class p2p_net_api_plugin : public plugin<p2p_net_api_plugin> {
public:
   APPBASE_PLUGIN_REQUIRES((p2p_net_plugin) (http_plugin))

   p2p_net_api_plugin() = default;
   p2p_net_api_plugin(const p2p_net_api_plugin&) = delete;
   p2p_net_api_plugin(p2p_net_api_plugin&&) = delete;
   p2p_net_api_plugin& operator=(const p2p_net_api_plugin&) = delete;
   p2p_net_api_plugin& operator=(p2p_net_api_plugin&&) = delete;
   virtual ~p2p_net_api_plugin() override = default;

   virtual void set_program_options(options_description& cli, options_description& cfg) override {}
   void plugin_initialize(const variables_map& vm);
   void plugin_startup();
   void plugin_shutdown() {}

private:
};

}
