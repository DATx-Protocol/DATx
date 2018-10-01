#list( APPEND CMAKE_MODULE_PATH "${CMAKE_SOURCE_DIR}/libraries/fc/GitVersionGen" )
#include( GetGitRevisionDescription )
#get_git_head_revision(GIT_REFSPEC DATX_GIT_REVISION_SHA)
#get_git_unix_timestamp(DATX_GIT_REVISION_UNIX_TIMESTAMP)
#git_describe(DATX_GIT_REVISION_DESCRIPTION --tags)
if(NOT DATX_GIT_REVISION_DESCRIPTION)
    set(DATX_GIT_REVISION_DESCRIPTION "unknown")
endif(NOT DATX_GIT_REVISION_DESCRIPTION)

file(GLOB HEADERS "include/datxio/utilities/*.hpp")

set(sources
   key_conversion.cpp
   string_escape.cpp
   tempdir.cpp
   words.cpp
   ${HEADERS})

#configure_file("${CMAKE_CURRENT_SOURCE_DIR}/git_revision.cpp.in" "${CMAKE_CURRENT_BINARY_DIR}/git_revision.cpp" @ONLY)
#list(APPEND sources "${CMAKE_CURRENT_BINARY_DIR}/git_revision.cpp")

add_library( datx_utilities
             ${sources}
             ${HEADERS} )
target_link_libraries( datx_utilities fc WAST WASM )
target_include_directories( datx_utilities
                            PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/include"
                                   "${CMAKE_CURRENT_SOURCE_DIR}/../wasm-jit/Include" )
if (USE_PCH)
  set_target_properties(datx_utilities PROPERTIES COTIRE_ADD_UNITY_BUILD FALSE)
  cotire(datx_utilities)
endif(USE_PCH)
