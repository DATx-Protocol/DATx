/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */

#include <DatxioLib/DatxioLib.hpp>

namespace asserter {
   struct assertdef {
      int8_t      condition;
      std::string message;

      DATXLIB_SERIALIZE( assertdef, (condition)(message) )
   };
}
