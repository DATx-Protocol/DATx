file(GLOB HEADERS "include/datxos/wallet_api_accessory/*.hpp")
add_library( wallet_api_accessory
             wallet_api_accessory.cpp
             ${HEADERS} )

target_link_libraries( wallet_api_accessory wallet_accessory http_accessory appbase )
target_include_directories( wallet_api_accessory PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/include" )
