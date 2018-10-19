file(GLOB HEADERS "include/datxos/http_client_accessory/*.hpp")
add_library( http_client_accessory
             http_client_accessory.cpp
             ${HEADERS} )

target_link_libraries( http_client_accessory appbase datxos_chain fc )
target_include_directories( http_client_accessory PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/include" )
