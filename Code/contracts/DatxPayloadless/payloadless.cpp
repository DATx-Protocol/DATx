#include <DatxioLib/DatxioLib.hpp>
#include <DatxioLib/print.hpp>
using namespace datxio;

class payloadless : public datxio::contract {
  public:
      using contract::contract;

      void doit() {
         print( "Im a payloadless action" );
      }
};

DATXIO_ABI( payloadless, (doit) )
