/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#include <fc/variant.hpp>
#include <fc/io/json.hpp>
#include <datxos/db_size_api_accessory/db_size_api_accessory.hpp>

namespace datxos {

static appbase::abstract_accessory& _db_size_api_accessory = app().register_accessory<db_size_api_accessory>();

using namespace datxos;

#define CALL(api_name, api_handle, call_name, INVOKE, http_response_code) \
{std::string("/v1/" #api_name "/" #call_name), \
   [api_handle](string, string body, url_response_callback cb) mutable { \
          try { \
             if (body.empty()) body = "{}"; \
             INVOKE \
             cb(http_response_code, fc::json::to_string(result)); \
          } catch (...) { \
             http_accessory::handle_exception(#api_name, #call_name, body, cb); \
          } \
       }}

#define INVOKE_R_V(api_handle, call_name) \
     auto result = api_handle->call_name();


void db_size_api_accessory::accessory_startup() {
   app().get_accessory<http_accessory>().add_api({
       CALL(db_size, this, get,
            INVOKE_R_V(this, get), 200),
   });
}

db_size_stats db_size_api_accessory::get() {
   chainbase::database& db = app().get_accessory<core_accessory>().chain().db();
   db_size_stats ret;

   ret.free_bytes = db.get_segment_manager()->get_free_memory();
   ret.size = db.get_segment_manager()->get_size();
   ret.used_bytes = ret.size - ret.free_bytes;

   chainbase::database::database_index_row_count_multiset indices = db.row_count_per_index();
   for(const auto& i : indices)
      ret.indices.emplace_back(db_size_index_count{i.second, i.first});

   return ret;
}

#undef INVOKE_R_V
#undef CALL

}
