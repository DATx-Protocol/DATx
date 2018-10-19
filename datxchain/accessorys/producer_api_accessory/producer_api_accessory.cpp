/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#include <datxos/producer_api_accessory/producer_api_accessory.hpp>
#include <datxos/chain/exceptions.hpp>

#include <fc/variant.hpp>
#include <fc/io/json.hpp>

#include <chrono>

namespace datxos { namespace detail {
  struct producer_api_accessory_response {
     std::string result;
  };
}}

FC_REFLECT(datxos::detail::producer_api_accessory_response, (result));

namespace datxos {

static appbase::abstract_accessory& _producer_api_accessory = app().register_accessory<producer_api_accessory>();

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
     datxos::detail::producer_api_accessory_response result{"ok"};

#define INVOKE_V_R_R(api_handle, call_name, in_param0, in_param1) \
     const auto& vs = fc::json::json::from_string(body).as<fc::variants>(); \
     api_handle.call_name(vs.at(0).as<in_param0>(), vs.at(1).as<in_param1>()); \
     datxos::detail::producer_api_accessory_response result{"ok"};

#define INVOKE_V_V(api_handle, call_name) \
     api_handle.call_name(); \
     datxos::detail::producer_api_accessory_response result{"ok"};


void producer_api_accessory::accessory_startup() {
   ilog("starting producer_api_accessory");
   // lifetime of accessory is lifetime of application
   auto& producer = app().get_accessory<producer_accessory>();

   app().get_accessory<http_accessory>().add_api({
       CALL(producer, producer, pause,
            INVOKE_V_V(producer, pause), 201),
       CALL(producer, producer, resume,
            INVOKE_V_V(producer, resume), 201),
       CALL(producer, producer, paused,
            INVOKE_R_V(producer, paused), 201),
       CALL(producer, producer, get_runtime_options,
            INVOKE_R_V(producer, get_runtime_options), 201),
       CALL(producer, producer, update_runtime_options,
            INVOKE_V_R(producer, update_runtime_options, producer_accessory::runtime_options), 201),
       CALL(producer, producer, add_greylist_accounts,
            INVOKE_V_R(producer, add_greylist_accounts, producer_accessory::greylist_params), 201),
       CALL(producer, producer, remove_greylist_accounts,
            INVOKE_V_R(producer, remove_greylist_accounts, producer_accessory::greylist_params), 201), 
       CALL(producer, producer, get_greylist,
            INVOKE_R_V(producer, get_greylist), 201),                 
       CALL(producer, producer, get_whitelist_blacklist,
            INVOKE_R_V(producer, get_whitelist_blacklist), 201),
       CALL(producer, producer, set_whitelist_blacklist, 
            INVOKE_V_R(producer, set_whitelist_blacklist, producer_accessory::whitelist_blacklist), 201),   
   });
}

void producer_api_accessory::accessory_initialize(const variables_map& options) {
   try {
      const auto& _http_accessory = app().get_accessory<http_accessory>();
      if( !_http_accessory.is_on_loopback()) {
         wlog( "\n"
               "**********SECURITY WARNING**********\n"
               "*                                  *\n"
               "* --        Producer API        -- *\n"
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
