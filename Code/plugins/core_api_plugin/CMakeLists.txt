file(GLOB HEADERS "include/datxio/core_api_plugin/*.hpp")
add_library( core_api_plugin
             core_api_plugin.cpp
             ${HEADERS} )

target_link_libraries( core_api_plugin core_plugin http_plugin appbase )
target_include_directories( core_api_plugin PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/include" )
