/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#pragma once
#include <datxos/core_accessory/core_accessory.hpp>
#include <datxos/http_accessory/http_accessory.hpp>

#include <appbase/application.hpp>
#include <datxos/chain/controller.hpp>

namespace datxos {
   using datxos::chain::controller;
   using std::unique_ptr;
   using namespace appbase;

   class core_api_accessory : public accessory<core_api_accessory> {
      public:
        APPBASE_accessory_REQUIRES((core_accessory)(http_accessory))

        core_api_accessory();
        virtual ~core_api_accessory();

        virtual void set_program_options(options_description&, options_description&) override;

        void accessory_initialize(const variables_map&);
        void accessory_startup();
        void accessory_shutdown();

      private:
        unique_ptr<class core_api_accessory_impl> my;
   };

}
