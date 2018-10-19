AppBase
--------------

The AppBase library provides a basic framework for building applications from
a set of accessorys. AppBase manages the accessory life-cycle and ensures that all
accessorys are configured, initialized, started, and shutdown in the proper order.

## Key Features

- Dynamically Specify accessorys to Load
- Automatically Load Dependent accessorys in Order
- accessorys can specify commandline arguments and configuration file options
- Program gracefully exits from SIGINT and SIGTERM
- Minimal Dependencies (Boost 1.60, c++14)

## Defining a Accessory

A simple example of a 2-accessory application can be found in the /examples directory. Each accessory has
a simple life cycle:

1. Initialize - parse configuration file options
2. Startup - start executing, using configuration file options
3. Shutdown - stop everything and free all resources

All accessorys complete the Initialize step before any accessory enters the Startup step. Any dependent accessory specified
by `APPBASE_accessory_REQUIRES` will be Initialized or Started prior to the accessory being Initialized or Started. 

Shutdown is called in the reverse order of Startup. 

```
class p2p_net_accessory : public appbase::accessory<p2p_net_accessory>
{
   public:
     p2p_net_accessory(){};
     ~p2p_net_accessory(){};

     APPBASE_accessory_REQUIRES( (core_accessory) );

     virtual void set_program_options( options_description& cli, options_description& cfg ) override
     {
        cfg.add_options()
              ("listen-endpoint", bpo::value<string>()->default_value( "127.0.0.1:9876" ), "The local IP address and port to listen for incoming connections.")
              ("remote-endpoint", bpo::value< vector<string> >()->composing(), "The IP address and port of a remote peer to sync with.")
              ("public-endpoint", bpo::value<string>()->default_value( "0.0.0.0:9876" ), "The public IP address and port that should be advertized to peers.")
              ;
     }

     void accessory_initialize( const variables_map& options ) { std::cout << "initialize net accessory\n"; }
     void accessory_startup()  { std::cout << "starting net accessory \n"; }
     void accessory_shutdown() { std::cout << "shutdown net accessory \n"; }

};

int main( int argc, char** argv ) {
   try {
      appbase::app().register_accessory<p2p_net_accessory>(); // implict registration of core_accessory dependency
      if( !appbase::app().initialize( argc, argv ) )
         return -1;
      appbase::app().startup();
      appbase::app().exec();
   } catch ( const boost::exception& e ) {
      std::cerr << boost::diagnostic_information(e) << "\n";
   } catch ( const std::exception& e ) {
      std::cerr << e.what() << "\n";
   } catch ( ... ) {
      std::cerr << "unknown exception\n";
   }
   std::cout << "exited cleanly\n";
   return 0;
}
```

This example can be used like follows:

```
./examples/appbase_example --accessory p2p_net_accessory
initialize chain accessory
initialize net accessory
starting chain accessory
starting net accessory
^C
shutdown net accessory
shutdown chain accessory
exited cleanly
```

### Boost ASIO 

AppBase maintains a singleton `application` instance which can be accessed via `appbase::app()`.  This 
application owns a `boost::asio::io_service` which starts running when `appbase::exec()` is called. If 
a accessory needs to perform IO or other asynchronous operations then it should dispatch it via 
`app().get_io_service().post( lambda )`.  

Because the app calls `io_service::run()` from within `application::exec()` all asynchronous operations
posted to the io_service should be run in the same thread.  

## Graceful Exit 

To trigger a graceful exit call `appbase::app().quit()` or send SIGTERM or SIGINT to the process.

## Dependencies 

1. c++14 or newer  (clang or g++)
2. Boost 1.60 or newer compiled with C++14 support

To compile boost with c++14 use:

```
./b2 ...  cxxflags="-std=c++0x -stdlib=libc++" linkflags="-stdlib=libc++" ...
```

