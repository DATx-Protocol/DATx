/**
 *  @file
 *  @copyright defined in datx/LICENSE.txt
 */

#include <datxio/utilities/tempdir.hpp>

#include <cstdlib>

namespace datxio { namespace utilities {

fc::path temp_directory_path()
{
   const char* datx_tempdir = getenv("DATX_TEMPDIR");
   if( datx_tempdir != nullptr )
      return fc::path( datx_tempdir );
   return fc::temp_directory_path() / "datx-tmp";
}

} } // datxio::utilities
