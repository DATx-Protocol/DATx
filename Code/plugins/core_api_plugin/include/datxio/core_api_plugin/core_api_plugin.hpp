/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#pragma once
#include <datxio/core_plugin/core_plugin.hpp>
#include <datxio/http_plugin/http_plugin.hpp>

#include <appbase/application.hpp>
#include <datxio/chain/controller.hpp>

namespace datxio {
   using datxio::chain::controller;
   using std::unique_ptr;
   using namespace appbase;

   class core_api_plugin : public plugin<core_api_plugin> {
      public:
        APPBASE_PLUGIN_REQUIRES((core_plugin)(http_plugin))

        core_api_plugin();
        virtual ~core_api_plugin();

        virtual void set_program_options(options_description&, options_description&) override;

        void plugin_initialize(const variables_map&);
        void plugin_startup();
        void plugin_shutdown();

      private:
        unique_ptr<class core_api_plugin_impl> my;
   };

}
