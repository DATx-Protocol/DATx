file(GLOB HEADERS "include/datxio/faucet_testp2p_net_plugin/*.hpp")
add_library( faucet_testp2p_net_plugin
             faucet_testp2p_net_plugin.cpp
             ${HEADERS} )

target_link_libraries( faucet_testp2p_net_plugin appbase fc http_plugin core_plugin )
target_include_directories( faucet_testp2p_net_plugin PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/include" )
