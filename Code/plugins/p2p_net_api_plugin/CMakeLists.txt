file(GLOB HEADERS "include/datxio/p2p_net_api_plugin/*.hpp")
add_library( p2p_net_api_plugin
             p2p_net_api_plugin.cpp
             ${HEADERS} )

target_link_libraries( p2p_net_api_plugin p2p_net_plugin http_plugin appbase )
target_include_directories( p2p_net_api_plugin PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/include" )
