/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */

#include <datxos/utilities/tempdir.hpp>

#include <cstdlib>

namespace datxos { namespace utilities {

fc::path temp_directory_path()
{
   const char* datx_tempdir = getenv("DATX_TEMPDIR");
   if( datx_tempdir != nullptr )
      return fc::path( datx_tempdir );
   return fc::temp_directory_path() / "datx-tmp";
}

} } // datxos::utilities
