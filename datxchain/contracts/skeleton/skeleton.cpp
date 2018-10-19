#include <DatxosLib/DatxosLib.hpp>

using namespace datxos;

class hello : public datxos::contract {
  public:
      using contract::contract;

      /// @abi action 
      void hi( account_name user ) {
         print( "Hello, ", name{user} );
      }
};

DATXOS_ABI( hello, (hi) )
