file(GLOB HEADERS "include/datxio/core_plugin/*.hpp")
add_library( core_plugin
             core_plugin.cpp
             ${HEADERS} )

target_link_libraries( core_plugin datxio_chain appbase )
target_include_directories( core_plugin PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/include" "${CMAKE_CURRENT_SOURCE_DIR}/../chain_interface/include" "${CMAKE_CURRENT_SOURCE_DIR}/../../libraries/appbase/include")
