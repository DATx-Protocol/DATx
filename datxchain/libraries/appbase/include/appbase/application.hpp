#pragma once
#include <appbase/accessory.hpp>
#include <appbase/channel.hpp>
#include <appbase/method.hpp>
#include <boost/filesystem/path.hpp>
#include <boost/core/demangle.hpp>
#include <typeindex>

namespace appbase {
   namespace bpo = boost::program_options;
   namespace bfs = boost::filesystem;

   class application
   {
      public:
         ~application();


         /** @brief Set version
          *
          * @param version Version output with -v/--version
          */
         void set_version(uint64_t version);
         /** @brief Get version
          *
          * @return Version output with -v/--version
          */
         uint64_t version() const;
         /** @brief Get version string; generated from git describe if available
          *
          * @return A string worthy of output with -v/--version, or "Unknown" if git not available
          */
         //string version_string() const;
         /** @brief Set default data directory
          *
          * @param data_dir Default data directory to use if not specified
          *                 on the command line.
          */
         void set_default_data_dir(const bfs::path& data_dir = "data-dir");
         /** @brief Get data directory
          *
          * @return Data directory, possibly from command line
          */
         bfs::path data_dir() const;
         /** @brief Set default config directory
          *
          * @param config_dir Default configuration directory to use if not
          *                   specified on the command line.
          */
         void set_default_config_dir(const bfs::path& config_dir = "etc");
         /** @brief Get config directory
          *
          * @return Config directory, possibly from command line
          */
         bfs::path config_dir() const;
         /** @brief Get logging configuration path.
          *
          * @return Logging configuration location from command line
          */
         bfs::path get_logging_conf() const;
         /**
          * @brief Looks for the --accessory commandline / config option and calls initialize on those accessorys
          *
          * @tparam Accessory List of accessorys to initalize even if not mentioned by configuration. For accessorys started by
          * configuration settings or dependency resolution, this template has no effect.
          * @return true if the application and accessorys were initialized, false or exception on error
          */
         template<typename... Accessory>
         bool                 initialize(int argc, char** argv) {
            return initialize_impl(argc, argv, {find_accessory<Accessory>()...});
         }

         void                  startup();
         void                  shutdown();

         /**
          *  Wait until quit(), SIGINT or SIGTERM and then shutdown
          */
         void                 exec();
         void                 quit();

         static application&  instance();

         abstract_accessory* find_accessory(const string& name)const;
         abstract_accessory& get_accessory(const string& name)const;

         template<typename Accessory>
         auto& register_accessory() {
            auto existing = find_accessory<Accessory>();
            if(existing)
               return *existing;

            auto plug = new Accessory();
            accessorys[plug->name()].reset(plug);
            plug->register_dependencies();
            return *plug;
         }

         template<typename Accessory>
         Accessory* find_accessory()const {
            string name = boost::core::demangle(typeid(Accessory).name());
            return dynamic_cast<Accessory*>(find_accessory(name));
         }

         template<typename Accessory>
         Accessory& get_accessory()const {
            auto ptr = find_accessory<Accessory>();
            return *ptr;
         }

         /**
          * Fetch a reference to the method declared by the passed in type.  This will construct the method
          * on first access.  This allows loose and deferred binding between accessorys
          *
          * @tparam MethodDecl - @ref appbase::method_decl
          * @return reference to the method described by the declaration
          */
         template<typename MethodDecl>
         auto get_method() -> std::enable_if_t<is_method_decl<MethodDecl>::value, typename MethodDecl::method_type&>
         {
            using method_type = typename MethodDecl::method_type;
            auto key = std::type_index(typeid(MethodDecl));
            auto itr = methods.find(key);
            if(itr != methods.end()) {
               return *method_type::get_method(itr->second);
            } else {
               methods.emplace(std::make_pair(key, method_type::make_unique()));
               return  *method_type::get_method(methods.at(key));
            }
         }

         /**
          * Fetch a reference to the channel declared by the passed in type.  This will construct the channel
          * on first access.  This allows loose and deferred binding between accessorys
          *
          * @tparam ChannelDecl - @ref appbase::channel_decl
          * @return reference to the channel described by the declaration
          */
         template<typename ChannelDecl>
         auto get_channel() -> std::enable_if_t<is_channel_decl<ChannelDecl>::value, typename ChannelDecl::channel_type&>
         {
            using channel_type = typename ChannelDecl::channel_type;
            auto key = std::type_index(typeid(ChannelDecl));
            auto itr = channels.find(key);
            if(itr != channels.end()) {
               return *channel_type::get_channel(itr->second);
            } else {
               channels.emplace(std::make_pair(key, channel_type::make_unique(io_serv)));
               return  *channel_type::get_channel(channels.at(key));
            }
         }

         boost::asio::io_service& get_io_service() { return *io_serv; }
      protected:
         template<typename Impl>
         friend class accessory;

         bool initialize_impl(int argc, char** argv, vector<abstract_accessory*> autostart_accessorys);

         /** these notifications get called from the accessory when their state changes so that
          * the application can call shutdown in the reverse order.
          */
         ///@{
         void accessory_initialized(abstract_accessory& plug){ initialized_accessorys.push_back(&plug); }
         void accessory_started(abstract_accessory& plug){ running_accessorys.push_back(&plug); }
         ///@}

      private:
         application(); ///< private because application is a singleton that should be accessed via instance()
         map<string, std::unique_ptr<abstract_accessory>> accessorys; ///< all registered accessorys
         vector<abstract_accessory*>                  initialized_accessorys; ///< stored in the order they were started running
         vector<abstract_accessory*>                  running_accessorys; ///< stored in the order they were started running

         map<std::type_index, erased_method_ptr>   methods;
         map<std::type_index, erased_channel_ptr>  channels;

         std::shared_ptr<boost::asio::io_service>  io_serv;

         void set_program_options();
         void write_default_config(const bfs::path& cfg_file);
         void print_default_config(std::ostream& os);
         std::unique_ptr<class application_impl> my;

   };

   application& app();


   template<typename Impl>
   class accessory : public abstract_accessory {
      public:
         accessory():_name(boost::core::demangle(typeid(Impl).name())){}
         virtual ~accessory(){}

         virtual state get_state()const override         { return _state; }
         virtual const std::string& name()const override { return _name; }

         virtual void register_dependencies() {
            static_cast<Impl*>(this)->accessory_requires([&](auto& plug){});
         }

         virtual void initialize(const variables_map& options) override {
            if(_state == registered) {
               _state = initialized;
               static_cast<Impl*>(this)->accessory_requires([&](auto& plug){ plug.initialize(options); });
               static_cast<Impl*>(this)->accessory_initialize(options);
               //ilog( "initializing accessory ${name}", ("name",name()) );
               app().accessory_initialized(*this);
            }
            assert(_state == initialized); /// if initial state was not registered, final state cannot be initiaized
         }

         virtual void startup() override {
            if(_state == initialized) {
               _state = started;
               static_cast<Impl*>(this)->accessory_requires([&](auto& plug){ plug.startup(); });
               static_cast<Impl*>(this)->accessory_startup();
               app().accessory_started(*this);
            }
            assert(_state == started); // if initial state was not initialized, final state cannot be started
         }

         virtual void shutdown() override {
            if(_state == started) {
               _state = stopped;
               //ilog( "shutting down accessory ${name}", ("name",name()) );
               static_cast<Impl*>(this)->accessory_shutdown();
            }
         }

      protected:
         accessory(const string& name) : _name(name){}

      private:
         state _state = abstract_accessory::registered;
         std::string _name;
   };
}
