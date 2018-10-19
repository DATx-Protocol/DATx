file(GLOB HEADERS "include/datxos/p2p_net_api_accessory/*.hpp")
add_library( p2p_net_api_accessory
             p2p_net_api_accessory.cpp
             ${HEADERS} )

target_link_libraries( p2p_net_api_accessory p2p_net_accessory http_accessory appbase )
target_include_directories( p2p_net_api_accessory PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/include" )
