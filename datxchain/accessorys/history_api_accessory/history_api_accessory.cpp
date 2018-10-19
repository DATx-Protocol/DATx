/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#include <datxos/history_api_accessory/history_api_accessory.hpp>
#include <datxos/chain/exceptions.hpp>

#include <fc/io/json.hpp>

namespace datxos {

using namespace datxos;

static appbase::abstract_accessory& _history_api_accessory = app().register_accessory<history_api_accessory>();

history_api_accessory::history_api_accessory(){}
history_api_accessory::~history_api_accessory(){}

void history_api_accessory::set_program_options(options_description&, options_description&) {}
void history_api_accessory::accessory_initialize(const variables_map&) {}

#define CALL(api_name, api_handle, api_namespace, call_name) \
{std::string("/v1/" #api_name "/" #call_name), \
   [this, api_handle](string, string body, url_response_callback cb) mutable { \
          try { \
             if (body.empty()) body = "{}"; \
             auto result = api_handle.call_name(fc::json::from_string(body).as<api_namespace::call_name ## _params>()); \
             cb(200, fc::json::to_string(result)); \
          } catch (...) { \
             http_accessory::handle_exception(#api_name, #call_name, body, cb); \
          } \
       }}

#define CHAIN_RO_CALL(call_name) CALL(history, ro_api, history_apis::read_only, call_name)
//#define CHAIN_RW_CALL(call_name) CALL(history, rw_api, history_apis::read_write, call_name)

void history_api_accessory::accessory_startup() {
   ilog( "starting history_api_accessory" );
   auto ro_api = app().get_accessory<history_accessory>().get_read_only_api();
   //auto rw_api = app().get_accessory<history_accessory>().get_read_write_api();

   app().get_accessory<http_accessory>().add_api({
//      CHAIN_RO_CALL(get_transaction),
      CHAIN_RO_CALL(get_actions),
      CHAIN_RO_CALL(get_transaction),
      CHAIN_RO_CALL(get_key_accounts),
      CHAIN_RO_CALL(get_controlled_accounts)
   });
}

void history_api_accessory::accessory_shutdown() {}

}
