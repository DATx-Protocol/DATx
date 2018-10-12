/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#pragma once
#include <appbase/application.hpp>
#include <datxio/core_plugin/core_plugin.hpp>
#include <datxio/p2p_net_plugin/protocol.hpp>

namespace datxio {
   using namespace appbase;

   struct connection_status {
      string            peer;
      bool              connecting = false;
      bool              syncing    = false;
      handshake_message last_handshake;
   };

   class p2p_net_plugin : public appbase::plugin<p2p_net_plugin>
   {
      public:
        p2p_net_plugin();
        virtual ~p2p_net_plugin();

        APPBASE_PLUGIN_REQUIRES((core_plugin))
        virtual void set_program_options(options_description& cli, options_description& cfg) override;

        void plugin_initialize(const variables_map& options);
        void plugin_startup();
        void plugin_shutdown();

        void   broadcast_block(const chain::signed_block &sb);

        string                       connect( const string& endpoint );
        string                       disconnect( const string& endpoint );
        optional<connection_status>  status( const string& endpoint )const;
        vector<connection_status>    connections()const;

        size_t num_peers() const;
      private:
        std::unique_ptr<class p2p_net_plugin_impl> my;
   };

}

FC_REFLECT( datxio::connection_status, (peer)(connecting)(syncing)(last_handshake) )
