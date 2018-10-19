file(GLOB HEADERS "include/datxos/core_accessory/*.hpp")
add_library( core_accessory
             core_accessory.cpp
             ${HEADERS} )

target_link_libraries( core_accessory datxos_chain appbase )
target_include_directories( core_accessory PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/include" "${CMAKE_CURRENT_SOURCE_DIR}/../chain_interface/include" "${CMAKE_CURRENT_SOURCE_DIR}/../../libraries/appbase/include")
