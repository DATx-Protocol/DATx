/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#pragma once
#include <appbase/application.hpp>
#include <datxos/core_accessory/core_accessory.hpp>
#include <datxos/p2p_net_accessory/protocol.hpp>

namespace datxos {
   using namespace appbase;

   struct connection_status {
      string            peer;
      bool              connecting = false;
      bool              syncing    = false;
      handshake_message last_handshake;
   };

   class p2p_net_accessory : public appbase::accessory<p2p_net_accessory>
   {
      public:
        p2p_net_accessory();
        virtual ~p2p_net_accessory();

        APPBASE_accessory_REQUIRES((core_accessory))
        virtual void set_program_options(options_description& cli, options_description& cfg) override;

        void accessory_initialize(const variables_map& options);
        void accessory_startup();
        void accessory_shutdown();

        void   broadcast_block(const chain::signed_block &sb);

        string                       connect( const string& endpoint );
        string                       disconnect( const string& endpoint );
        optional<connection_status>  status( const string& endpoint )const;
        vector<connection_status>    connections()const;

        size_t num_peers() const;
      private:
        std::unique_ptr<class p2p_net_accessory_impl> my;
   };

}

FC_REFLECT( datxos::connection_status, (peer)(connecting)(syncing)(last_handshake) )
