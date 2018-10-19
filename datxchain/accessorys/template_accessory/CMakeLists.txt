file(GLOB HEADERS "include/datxos/template_accessory/*.hpp")
add_library( template_accessory
             template_accessory.cpp
             ${HEADERS} )

target_link_libraries( template_accessory appbase fc )
target_include_directories( template_accessory PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/include" )
