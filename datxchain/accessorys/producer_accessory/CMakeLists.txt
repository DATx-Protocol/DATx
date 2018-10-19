file(GLOB HEADERS "include/datxos/producer_accessory/*.hpp")

add_library( producer_accessory
             producer_accessory.cpp
             ${HEADERS}
           )

target_link_libraries( producer_accessory core_accessory http_client_accessory appbase datxos_chain datx_utilities )
target_include_directories( producer_accessory
                            PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/include" "${CMAKE_CURRENT_SOURCE_DIR}/../chain_interface/include" )
