file(GLOB HEADERS "include/datxos/http_accessory/*.hpp")
add_library( http_accessory
             http_accessory.cpp
             ${HEADERS} )

target_link_libraries( http_accessory datxos_chain appbase fc )
target_include_directories( http_accessory PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/include" )
