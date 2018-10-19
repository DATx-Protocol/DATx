#include <DatxosLib/DatxosLib.hpp>
#include <DatxosLib/print.hpp>
using namespace datxos;

class payloadless : public datxos::contract {
  public:
      using contract::contract;

      void doit() {
         print( "Im a payloadless action" );
      }
};

DATXOS_ABI( payloadless, (doit) )
