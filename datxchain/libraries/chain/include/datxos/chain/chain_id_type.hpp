/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#pragma once

#include <fc/crypto/sha256.hpp>

struct hello;

namespace datxos {

   class p2p_net_accessory_impl;
   struct handshake_message;

   namespace chain_apis {
      class read_only;
   }

namespace chain {

   struct chain_id_type : public fc::sha256 {
      using fc::sha256::sha256;

      template<typename T>
      inline friend T& operator<<( T& ds, const chain_id_type& cid ) {
        ds.write( cid.data(), cid.data_size() );
        return ds;
      }

      template<typename T>
      inline friend T& operator>>( T& ds, chain_id_type& cid ) {
        ds.read( cid.data(), cid.data_size() );
        return ds;
      }

      void reflector_verify()const;

      private:
         chain_id_type() = default;

         // Some exceptions are unfortunately necessary:
         template<typename T>
         friend T fc::variant::as()const;

         friend class datxos::chain_apis::read_only;

         friend class datxos::p2p_net_accessory_impl;
         friend struct datxos::handshake_message;

         friend struct ::hello; // TODO: Rushed hack to support bp2p_accessory. Need a better solution.
   };

} }  // namespace datxos::chain

namespace fc {
  class variant;
  void to_variant(const datxos::chain::chain_id_type& cid, fc::variant& v);
  void from_variant(const fc::variant& v, datxos::chain::chain_id_type& cid);
} // fc
