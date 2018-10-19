/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#pragma once

#include <datxos/chain/types.hpp>
#include <datxos/chain/contract_types.hpp>

namespace datxos { namespace chain {

   class apply_context;

   /**
    * @defgroup native_action_handlers Native Action Handlers
    */
   ///@{
   void apply_datxos_newaccount(apply_context&);
   void apply_datxos_updateauth(apply_context&);
   void apply_datxos_deleteauth(apply_context&);
   void apply_datxos_linkauth(apply_context&);
   void apply_datxos_unlinkauth(apply_context&);

   /*
   void apply_datxos_postrecovery(apply_context&);
   void apply_datxos_passrecovery(apply_context&);
   void apply_datxos_vetorecovery(apply_context&);
   */

   void apply_datxos_setcode(apply_context&);
   void apply_datxos_setabi(apply_context&);

   void apply_datxos_canceldelay(apply_context&);
   ///@}  end action handlers

} } /// namespace datxos::chain
