/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */

#include <DatxosLib/DatxosLib.hpp>

namespace asserter {
   struct assertdef {
      int8_t      condition;
      std::string message;

      DATXLIB_SERIALIZE( assertdef, (condition)(message) )
   };
}
