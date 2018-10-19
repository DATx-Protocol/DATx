file(GLOB HEADERS "include/datxos/history_accessory/*.hpp")
add_library( history_accessory
             history_accessory.cpp
             ${HEADERS} )

target_link_libraries( history_accessory core_accessory datxos_chain appbase )
target_include_directories( history_accessory PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/include" )
