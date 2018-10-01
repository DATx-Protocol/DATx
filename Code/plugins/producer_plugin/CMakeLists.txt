file(GLOB HEADERS "include/datxio/producer_plugin/*.hpp")

add_library( producer_plugin
             producer_plugin.cpp
             ${HEADERS}
           )

target_link_libraries( producer_plugin chain_plugin http_client_plugin appbase datxio_chain datx_utilities )
target_include_directories( producer_plugin
                            PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/include" "${CMAKE_CURRENT_SOURCE_DIR}/../chain_interface/include" )
