file(GLOB HEADERS "include/datxos/login_accessory/*.hpp")
add_library( login_accessory
             login_accessory.cpp
             ${HEADERS} )

target_link_libraries( login_accessory core_accessory http_accessory appbase )
target_include_directories( login_accessory PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/include" )
