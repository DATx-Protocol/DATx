file(GLOB HEADERS "include/datxos/db_size_api_accessory/*.hpp")
add_library( db_size_api_accessory
             db_size_api_accessory.cpp
             ${HEADERS} )

target_link_libraries( db_size_api_accessory http_accessory core_accessory )
target_include_directories( db_size_api_accessory PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/include" )
