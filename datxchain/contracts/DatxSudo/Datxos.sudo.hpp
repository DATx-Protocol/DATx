#pragma once

#include <DatxosLib/DatxosLib.hpp>

namespace datxos {

   class sudo : public contract {
      public:
         sudo( account_name self ):contract(self){}

         void exec();

   };

} /// namespace datxos
