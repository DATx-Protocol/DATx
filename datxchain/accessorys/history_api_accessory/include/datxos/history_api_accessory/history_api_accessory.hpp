/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */

#pragma once
#include <datxos/history_accessory/history_accessory.hpp>
#include <datxos/core_accessory/core_accessory.hpp>
#include <datxos/http_accessory/http_accessory.hpp>

#include <appbase/application.hpp>

namespace datxos {

   using namespace appbase;

   class history_api_accessory : public accessory<history_api_accessory> {
      public:
        APPBASE_accessory_REQUIRES((history_accessory)(core_accessory)(http_accessory))

        history_api_accessory();
        virtual ~history_api_accessory();

        virtual void set_program_options(options_description&, options_description&) override;

        void accessory_initialize(const variables_map&);
        void accessory_startup();
        void accessory_shutdown();

      private:
   };

}
