file( GLOB HEADERS "include/datxos/history_api_accessory/*.hpp" )
add_library( history_api_accessory
             history_api_accessory.cpp
             ${HEADERS} )

target_link_libraries( history_api_accessory history_accessory core_accessory http_accessory appbase )
target_include_directories( history_api_accessory PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/include" )
