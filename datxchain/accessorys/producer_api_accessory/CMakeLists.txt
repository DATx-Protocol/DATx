file(GLOB HEADERS "include/datxos/producer_api_accessory/*.hpp")
add_library( producer_api_accessory
             producer_api_accessory.cpp
             ${HEADERS} )

target_link_libraries( producer_api_accessory producer_accessory http_accessory appbase )
target_include_directories( producer_api_accessory PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/include" )
