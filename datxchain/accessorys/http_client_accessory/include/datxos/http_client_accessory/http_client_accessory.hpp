/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#pragma once
#include <appbase/application.hpp>
#include <fc/network/http/http_client.hpp>

namespace datxos {
   using namespace appbase;
   using fc::http_client;

   class http_client_accessory : public appbase::accessory<http_client_accessory>
   {
      public:
        http_client_accessory();
        virtual ~http_client_accessory();

        APPBASE_accessory_REQUIRES()
        virtual void set_program_options(options_description&, options_description& cfg) override;

        void accessory_initialize(const variables_map& options);
        void accessory_startup();
        void accessory_shutdown();

        http_client& get_client() {
           return *my;
        }

      private:
        std::unique_ptr<http_client> my;
   };

}
