file(GLOB HEADERS "include/datxos/faucet_testp2p_net_accessory/*.hpp")
add_library( faucet_testp2p_net_accessory
             faucet_testp2p_net_accessory.cpp
             ${HEADERS} )

target_link_libraries( faucet_testp2p_net_accessory appbase fc http_accessory core_accessory )
target_include_directories( faucet_testp2p_net_accessory PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/include" )
