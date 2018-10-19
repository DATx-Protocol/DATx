file(GLOB HEADERS "include/datxos/core_api_accessory/*.hpp")
add_library( core_api_accessory
             core_api_accessory.cpp
             ${HEADERS} )

target_link_libraries( core_api_accessory core_accessory http_accessory appbase )
target_include_directories( core_api_accessory PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/include" )
