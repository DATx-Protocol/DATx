/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#pragma once
#include <appbase/application.hpp>

#include <datxio/core_plugin/core_plugin.hpp>

namespace fc { class variant; }

namespace datxio {
   using chain::transaction_id_type;
   using std::shared_ptr;
   using namespace appbase;
   using chain::name;
   using fc::optional;
   using chain::uint128_t;

   typedef shared_ptr<class bp2p_plugin_impl> bnet_ptr;
   typedef shared_ptr<const class bp2p_plugin_impl> bnet_const_ptr;



/**
 *  This plugin tracks all actions and keys associated with a set of configured accounts. It enables
 *  wallets to paginate queries for bnet.  
 *
 *  An action will be included in the account's bnet if any of the following:
 *     - receiver
 *     - any account named in auth list
 *
 *  A key will be linked to an account if the key is referneced in authorities of updateauth or newaccount 
 */
class bp2p_plugin : public plugin<bp2p_plugin> {
   public:
      APPBASE_PLUGIN_REQUIRES((core_plugin))

      bp2p_plugin();
      virtual ~bp2p_plugin();

      virtual void set_program_options(options_description& cli, options_description& cfg) override;

      void plugin_initialize(const variables_map& options);
      void plugin_startup();
      void plugin_shutdown();

   private:
      bnet_ptr my;
};

} /// namespace datxio


