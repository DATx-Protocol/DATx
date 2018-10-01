/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#pragma once

#include <datxio/chain/types.hpp>
#include <datxio/chain/contract_types.hpp>

namespace datxio { namespace chain {

   class apply_context;

   /**
    * @defgroup native_action_handlers Native Action Handlers
    */
   ///@{
   void apply_datxio_newaccount(apply_context&);
   void apply_datxio_updateauth(apply_context&);
   void apply_datxio_deleteauth(apply_context&);
   void apply_datxio_linkauth(apply_context&);
   void apply_datxio_unlinkauth(apply_context&);

   /*
   void apply_datxio_postrecovery(apply_context&);
   void apply_datxio_passrecovery(apply_context&);
   void apply_datxio_vetorecovery(apply_context&);
   */

   void apply_datxio_setcode(apply_context&);
   void apply_datxio_setabi(apply_context&);

   void apply_datxio_canceldelay(apply_context&);
   ///@}  end action handlers

} } /// namespace datxio::chain
