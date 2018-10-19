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

class login_accessory : public accessory<login_accessory> {
 public:
   APPBASE_accessory_REQUIRES((core_accessory)(http_accessory))

   login_accessory();
   virtual ~login_accessory();

   virtual void set_program_options(options_description&, options_description&) override;
   void accessory_initialize(const variables_map&);
   void accessory_startup();
   void accessory_shutdown();

   struct start_login_request_params {
      chain::time_point_sec expiration_time;
   };

   struct start_login_request_results {
      chain::public_key_type server_ephemeral_pub_key;
   };

   struct finalize_login_request_params {
      chain::public_key_type server_ephemeral_pub_key;
      chain::public_key_type client_ephemeral_pub_key;
      chain::permission_level permission;
      std::string data;
      std::vector<chain::signature_type> signatures;
   };

   struct finalize_login_request_results {
      chain::sha256 digest{};
      flat_set<chain::public_key_type> recovered_keys{};
      bool permission_satisfied = false;
      std::string error{};
   };

   struct do_not_use_gen_r1_key_params {};

   struct do_not_use_gen_r1_key_results {
      chain::public_key_type pub_key;
      chain::private_key_type priv_key;
   };

   struct do_not_use_sign_params {
      chain::private_key_type priv_key;
      chain::bytes data;
   };

   struct do_not_use_sign_results {
      chain::signature_type sig;
   };

   struct do_not_use_get_secret_params {
      chain::public_key_type pub_key;
      chain::private_key_type priv_key;
   };

   struct do_not_use_get_secret_results {
      chain::sha512 secret;
   };

   start_login_request_results start_login_request(const start_login_request_params&);
   finalize_login_request_results finalize_login_request(const finalize_login_request_params&);

   do_not_use_gen_r1_key_results do_not_use_gen_r1_key(const do_not_use_gen_r1_key_params&);
   do_not_use_sign_results do_not_use_sign(const do_not_use_sign_params&);
   do_not_use_get_secret_results do_not_use_get_secret(const do_not_use_get_secret_params&);

 private:
   unique_ptr<class login_accessory_impl> my;
};

} // namespace datxos

FC_REFLECT(datxos::login_accessory::start_login_request_params, (expiration_time))
FC_REFLECT(datxos::login_accessory::start_login_request_results, (server_ephemeral_pub_key))
FC_REFLECT(datxos::login_accessory::finalize_login_request_params,
           (server_ephemeral_pub_key)(client_ephemeral_pub_key)(permission)(data)(signatures))
FC_REFLECT(datxos::login_accessory::finalize_login_request_results, (digest)(recovered_keys)(permission_satisfied)(error))

FC_REFLECT_EMPTY(datxos::login_accessory::do_not_use_gen_r1_key_params)
FC_REFLECT(datxos::login_accessory::do_not_use_gen_r1_key_results, (pub_key)(priv_key))
FC_REFLECT(datxos::login_accessory::do_not_use_sign_params, (priv_key)(data))
FC_REFLECT(datxos::login_accessory::do_not_use_sign_results, (sig))
FC_REFLECT(datxos::login_accessory::do_not_use_get_secret_params, (pub_key)(priv_key))
FC_REFLECT(datxos::login_accessory::do_not_use_get_secret_results, (secret))
