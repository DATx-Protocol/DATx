#pragma once

#include <DatxioLib/DatxioLib.hpp>

namespace datxio {

   class sudo : public contract {
      public:
         sudo( account_name self ):contract(self){}

         void exec();

   };

} /// namespace datxio
