file(GLOB HEADERS "include/datxos/bp2p_accessory/*.hpp")
add_library( bp2p_accessory
             bp2p_accessory.cpp
             ${HEADERS} )

target_link_libraries( bp2p_accessory core_accessory datxos_chain appbase )
target_include_directories( bp2p_accessory PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/include" )
