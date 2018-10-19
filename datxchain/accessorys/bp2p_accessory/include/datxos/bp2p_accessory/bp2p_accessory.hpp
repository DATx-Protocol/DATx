/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#pragma once
#include <appbase/application.hpp>

#include <datxos/core_accessory/core_accessory.hpp>

namespace fc { class variant; }

namespace datxos {
   using chain::transaction_id_type;
   using std::shared_ptr;
   using namespace appbase;
   using chain::name;
   using fc::optional;
   using chain::uint128_t;

   typedef shared_ptr<class bp2p_accessory_impl> bnet_ptr;
   typedef shared_ptr<const class bp2p_accessory_impl> bnet_const_ptr;



/**
 *  This accessory tracks all actions and keys associated with a set of configured accounts. It enables
 *  wallets to paginate queries for bnet.  
 *
 *  An action will be included in the account's bnet if any of the following:
 *     - receiver
 *     - any account named in auth list
 *
 *  A key will be linked to an account if the key is referneced in authorities of updateauth or newaccount 
 */
class bp2p_accessory : public accessory<bp2p_accessory> {
   public:
      APPBASE_accessory_REQUIRES((core_accessory))

      bp2p_accessory();
      virtual ~bp2p_accessory();

      virtual void set_program_options(options_description& cli, options_description& cfg) override;

      void accessory_initialize(const variables_map& options);
      void accessory_startup();
      void accessory_shutdown();

   private:
      bnet_ptr my;
};

} /// namespace datxos


