file(GLOB HEADERS "include/datxio/bp2p_plugin/*.hpp")
add_library( bp2p_plugin
             bp2p_plugin.cpp
             ${HEADERS} )

target_link_libraries( bp2p_plugin core_plugin datxio_chain appbase )
target_include_directories( bp2p_plugin PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/include" )
