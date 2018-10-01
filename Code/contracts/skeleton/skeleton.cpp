#include <DatxioLib/DatxioLib.hpp>

using namespace datxio;

class hello : public datxio::contract {
  public:
      using contract::contract;

      /// @abi action 
      void hi( account_name user ) {
         print( "Hello, ", name{user} );
      }
};

DATXIO_ABI( hello, (hi) )
