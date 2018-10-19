/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#include <datxos/wallet_api_accessory/wallet_api_accessory.hpp>
#include <datxos/wallet_accessory/wallet_manager.hpp>
#include <datxos/chain/exceptions.hpp>
#include <datxos/chain/transaction.hpp>

#include <fc/variant.hpp>
#include <fc/io/json.hpp>

#include <chrono>

namespace datxos { namespace detail {
  struct wallet_api_accessory_empty {};
}}

FC_REFLECT(datxos::detail::wallet_api_accessory_empty, );

namespace datxos {

static appbase::abstract_accessory& _wallet_api_accessory = app().register_accessory<wallet_api_accessory>();

using namespace datxos;

#define CALL(api_name, api_handle, call_name, INVOKE, http_response_code) \
{std::string("/v1/" #api_name "/" #call_name), \
   [&api_handle](string, string body, url_response_callback cb) mutable { \
          try { \
             if (body.empty()) body = "{}"; \
             INVOKE \
             cb(http_response_code, fc::json::to_string(result)); \
          } catch (...) { \
             http_accessory::handle_exception(#api_name, #call_name, body, cb); \
          } \
       }}

#define INVOKE_R_R(api_handle, call_name, in_param) \
     auto result = api_handle.call_name(fc::json::from_string(body).as<in_param>());

#define INVOKE_R_R_R(api_handle, call_name, in_param0, in_param1) \
     const auto& vs = fc::json::json::from_string(body).as<fc::variants>(); \
     auto result = api_handle.call_name(vs.at(0).as<in_param0>(), vs.at(1).as<in_param1>());

#define INVOKE_R_R_R_R(api_handle, call_name, in_param0, in_param1, in_param2) \
     const auto& vs = fc::json::json::from_string(body).as<fc::variants>(); \
     auto result = api_handle.call_name(vs.at(0).as<in_param0>(), vs.at(1).as<in_param1>(), vs.at(2).as<in_param2>());

#define INVOKE_R_V(api_handle, call_name) \
     auto result = api_handle.call_name();

#define INVOKE_V_R(api_handle, call_name, in_param) \
     api_handle.call_name(fc::json::from_string(body).as<in_param>()); \
     datxos::detail::wallet_api_accessory_empty result;

#define INVOKE_V_R_R(api_handle, call_name, in_param0, in_param1) \
     const auto& vs = fc::json::json::from_string(body).as<fc::variants>(); \
     api_handle.call_name(vs.at(0).as<in_param0>(), vs.at(1).as<in_param1>()); \
     datxos::detail::wallet_api_accessory_empty result;

#define INVOKE_V_R_R_R(api_handle, call_name, in_param0, in_param1, in_param2) \
     const auto& vs = fc::json::json::from_string(body).as<fc::variants>(); \
     api_handle.call_name(vs.at(0).as<in_param0>(), vs.at(1).as<in_param1>(), vs.at(2).as<in_param2>()); \
     datxos::detail::wallet_api_accessory_empty result;

#define INVOKE_V_V(api_handle, call_name) \
     api_handle.call_name(); \
     datxos::detail::wallet_api_accessory_empty result;


void wallet_api_accessory::accessory_startup() {
   ilog("starting wallet_api_accessory");
   // lifetime of accessory is lifetime of application
   auto& wallet_mgr = app().get_accessory<wallet_accessory>().get_wallet_manager();

   app().get_accessory<http_accessory>().add_api({
       CALL(wallet, wallet_mgr, set_timeout,
            INVOKE_V_R(wallet_mgr, set_timeout, int64_t), 200),
       CALL(wallet, wallet_mgr, sign_transaction,
            INVOKE_R_R_R_R(wallet_mgr, sign_transaction, chain::signed_transaction, flat_set<public_key_type>, chain::chain_id_type), 201),
       CALL(wallet, wallet_mgr, sign_digest,
            INVOKE_R_R_R(wallet_mgr, sign_digest, chain::digest_type, public_key_type), 201),
       CALL(wallet, wallet_mgr, create,
            INVOKE_R_R(wallet_mgr, create, std::string), 201),
       CALL(wallet, wallet_mgr, open,
            INVOKE_V_R(wallet_mgr, open, std::string), 200),
       CALL(wallet, wallet_mgr, lock_all,
            INVOKE_V_V(wallet_mgr, lock_all), 200),
       CALL(wallet, wallet_mgr, lock,
            INVOKE_V_R(wallet_mgr, lock, std::string), 200),
       CALL(wallet, wallet_mgr, unlock,
            INVOKE_V_R_R(wallet_mgr, unlock, std::string, std::string), 200),
       CALL(wallet, wallet_mgr, import_key,
            INVOKE_V_R_R(wallet_mgr, import_key, std::string, std::string), 201),
       CALL(wallet, wallet_mgr, remove_key,
            INVOKE_V_R_R_R(wallet_mgr, remove_key, std::string, std::string, std::string), 201),
       CALL(wallet, wallet_mgr, create_key,
            INVOKE_R_R_R(wallet_mgr, create_key, std::string, std::string), 201),
       CALL(wallet, wallet_mgr, list_wallets,
            INVOKE_R_V(wallet_mgr, list_wallets), 200),
       CALL(wallet, wallet_mgr, list_keys,
            INVOKE_R_R_R(wallet_mgr, list_keys, std::string, std::string), 200),
       CALL(wallet, wallet_mgr, get_public_keys,
            INVOKE_R_V(wallet_mgr, get_public_keys), 200)
   });
}

void wallet_api_accessory::accessory_initialize(const variables_map& options) {
   try {
      const auto& _http_accessory = app().get_accessory<http_accessory>();
      if( !_http_accessory.is_on_loopback()) {
         if( !_http_accessory.is_secure()) {
            elog( "\n"
                  "********!!!SECURITY ERROR!!!********\n"
                  "*                                  *\n"
                  "* --       Wallet API           -- *\n"
                  "* - EXPOSED to the LOCAL NETWORK - *\n"
                  "* -  HTTP RPC is NOT encrypted   - *\n"
                  "* - Password and/or Private Keys - *\n"
                  "* - are at HIGH risk of exposure - *\n"
                  "*                                  *\n"
                  "************************************\n" );
         } else {
            wlog( "\n"
                  "**********SECURITY WARNING**********\n"
                  "*                                  *\n"
                  "* --       Wallet API           -- *\n"
                  "* - EXPOSED to the LOCAL NETWORK - *\n"
                  "* - Password and/or Private Keys - *\n"
                  "* -   are at risk of exposure    - *\n"
                  "*                                  *\n"
                  "************************************\n" );
         }
      }
   } FC_LOG_AND_RETHROW()
}


#undef INVOKE_R_R
#undef INVOKE_R_R_R_R
#undef INVOKE_R_V
#undef INVOKE_V_R
#undef INVOKE_V_R_R
#undef INVOKE_V_V
#undef CALL

}
