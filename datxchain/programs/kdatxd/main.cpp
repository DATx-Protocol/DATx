/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */
#include <appbase/application.hpp>

#include <datxos/http_accessory/http_accessory.hpp>
#include <datxos/wallet_accessory/wallet_accessory.hpp>
#include <datxos/wallet_api_accessory/wallet_api_accessory.hpp>

#include <fc/log/logger_config.hpp>
#include <fc/exception/exception.hpp>

#include <boost/exception/diagnostic_information.hpp>

#include <pwd.h>

using namespace appbase;
using namespace datxos;

bfs::path determine_home_directory()
{
   bfs::path home;
   struct passwd* pwd = getpwuid(getuid());
   if(pwd) {
      home = pwd->pw_dir;
   }
   else {
      home = getenv("HOME");
   }
   if(home.empty())
      home = "./";
   return home;
}

int main(int argc, char** argv)
{
   try {
      bfs::path home = determine_home_directory();
      app().set_default_data_dir(home / "datxos-wallet");
      app().set_default_config_dir(home / "datxos-wallet");
      app().register_accessory<wallet_api_accessory>();
      if(!app().initialize<wallet_accessory, wallet_api_accessory, http_accessory>(argc, argv))
         return -1;
      auto& http = app().get_accessory<http_accessory>();
      http.add_handler("/v1/kdatxd/stop", [](string, string, url_response_callback cb) { cb(200, "{}"); std::raise(SIGTERM); } );
      app().startup();
      app().exec();
   } catch (const fc::exception& e) {
      elog("${e}", ("e",e.to_detail_string()));
   } catch (const boost::exception& e) {
      elog("${e}", ("e",boost::diagnostic_information(e)));
   } catch (const std::exception& e) {
      elog("${e}", ("e",e.what()));
   } catch (...) {
      elog("unknown exception");
   }
   return 0;
}
