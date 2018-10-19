/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#pragma once

#include <datxos/producer_accessory/producer_accessory.hpp>
#include <datxos/http_accessory/http_accessory.hpp>

#include <appbase/application.hpp>

namespace datxos {

using namespace appbase;

class producer_api_accessory : public accessory<producer_api_accessory> {
   public:
      APPBASE_accessory_REQUIRES((producer_accessory) (http_accessory))

      producer_api_accessory() = default;
      producer_api_accessory(const producer_api_accessory&) = delete;
      producer_api_accessory(producer_api_accessory&&) = delete;
      producer_api_accessory& operator=(const producer_api_accessory&) = delete;
      producer_api_accessory& operator=(producer_api_accessory&&) = delete;
      virtual ~producer_api_accessory() override = default;

      virtual void set_program_options(options_description& cli, options_description& cfg) override {}
      void accessory_initialize(const variables_map& vm);
      void accessory_startup();
      void accessory_shutdown() {}

   private:
};

}
