file(GLOB HEADERS "include/datxos/p2p_net_accessory/*.hpp" )
add_library( p2p_net_accessory
             p2p_net_accessory.cpp
             ${HEADERS} )

target_link_libraries( p2p_net_accessory core_accessory producer_accessory appbase fc )
target_include_directories( p2p_net_accessory PUBLIC ${CMAKE_CURRENT_SOURCE_DIR}/include ${CMAKE_CURRENT_SOURCE_DIR}/../chain_interface/include  "${CMAKE_CURRENT_SOURCE_DIR}/../../libraries/appbase/include")
