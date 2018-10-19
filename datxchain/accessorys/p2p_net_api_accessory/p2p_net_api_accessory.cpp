/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#include <datxos/p2p_net_api_accessory/p2p_net_api_accessory.hpp>
#include <datxos/chain/exceptions.hpp>
#include <datxos/chain/transaction.hpp>

#include <fc/variant.hpp>
#include <fc/io/json.hpp>

#include <chrono>

namespace datxos { namespace detail {
  struct p2p_net_api_accessory_empty {};
}}

FC_REFLECT(datxos::detail::p2p_net_api_accessory_empty, );

namespace datxos {

static appbase::abstract_accessory& _p2p_net_api_accessory = app().register_accessory<p2p_net_api_accessory>();

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

#define INVOKE_R_R_R_R(api_handle, call_name, in_param0, in_param1, in_param2) \
     const auto& vs = fc::json::json::from_string(body).as<fc::variants>(); \
     auto result = api_handle.call_name(vs.at(0).as<in_param0>(), vs.at(1).as<in_param1>(), vs.at(2).as<in_param2>());

#define INVOKE_R_V(api_handle, call_name) \
     auto result = api_handle.call_name();

#define INVOKE_V_R(api_handle, call_name, in_param) \
     api_handle.call_name(fc::json::from_string(body).as<in_param>()); \
     datxos::detail::p2p_net_api_accessory_empty result;

#define INVOKE_V_R_R(api_handle, call_name, in_param0, in_param1) \
     const auto& vs = fc::json::json::from_string(body).as<fc::variants>(); \
     api_handle.call_name(vs.at(0).as<in_param0>(), vs.at(1).as<in_param1>()); \
     datxos::detail::p2p_net_api_accessory_empty result;

#define INVOKE_V_V(api_handle, call_name) \
     api_handle.call_name(); \
     datxos::detail::p2p_net_api_accessory_empty result;


void p2p_net_api_accessory::accessory_startup() {
   ilog("starting p2p_net_api_accessory");
   // lifetime of accessory is lifetime of application
   auto& net_mgr = app().get_accessory<p2p_net_accessory>();

   app().get_accessory<http_accessory>().add_api({
    //   CALL(net, net_mgr, set_timeout,
    //        INVOKE_V_R(net_mgr, set_timeout, int64_t), 200),
    //   CALL(net, net_mgr, sign_transaction,
    //        INVOKE_R_R_R_R(net_mgr, sign_transaction, chain::signed_transaction, flat_set<public_key_type>, chain::chain_id_type), 201),
       CALL(net, net_mgr, connect,
            INVOKE_R_R(net_mgr, connect, std::string), 201),
       CALL(net, net_mgr, disconnect,
            INVOKE_R_R(net_mgr, disconnect, std::string), 201),
       CALL(net, net_mgr, status,
            INVOKE_R_R(net_mgr, status, std::string), 201),
       CALL(net, net_mgr, connections,
            INVOKE_R_V(net_mgr, connections), 201),
    //   CALL(net, net_mgr, open,
    //        INVOKE_V_R(net_mgr, open, std::string), 200),
   });
}

void p2p_net_api_accessory::accessory_initialize(const variables_map& options) {
   try {
      const auto& _http_accessory = app().get_accessory<http_accessory>();
      if( !_http_accessory.is_on_loopback()) {
         wlog( "\n"
               "**********SECURITY WARNING**********\n"
               "*                                  *\n"
               "* --         Net API            -- *\n"
               "* - EXPOSED to the LOCAL NETWORK - *\n"
               "* - USE ONLY ON SECURE NETWORKS! - *\n"
               "*                                  *\n"
               "************************************\n" );
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
