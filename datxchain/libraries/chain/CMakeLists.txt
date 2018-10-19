configure_file(${CMAKE_CURRENT_SOURCE_DIR}/include/datxos/chain/core_symbol.hpp.in ${CMAKE_CURRENT_BINARY_DIR}/include/datxos/chain/core_symbol.hpp)
configure_file(${CMAKE_CURRENT_SOURCE_DIR}/genesis_state_root_key.cpp.in ${CMAKE_CURRENT_BINARY_DIR}/genesis_state_root_key.cpp)

file(GLOB HEADERS "include/datxos/chain/*.hpp"
                  "include/datxos/chain/webassembly/*.hpp"
                  "${CMAKE_CURRENT_BINARY_DIR}/include/datxos/chain/core_symbol.hpp" )


add_library( datxos_chain
             merkle.cpp
             name.cpp
             transaction.cpp
             block_header.cpp
             block_header_state.cpp
             block_state.cpp
             fork_database.cpp
             controller.cpp
             authorization_manager.cpp
             resource_limits.cpp
             block_log.cpp
             transaction_context.cpp
             datxos_contract.cpp
             datxos_contract_abi.cpp
             chain_config.cpp
             chain_id_type.cpp
             genesis_state.cpp
             ${CMAKE_CURRENT_BINARY_DIR}/genesis_state_root_key.cpp
              wast_to_wasm.cpp
              wasm_interface.cpp
              wasm_datxos_validation.cpp
              wasm_datxos_injection.cpp
              apply_context.cpp
              abi_serializer.cpp
              asset.cpp

             webassembly/wavm.cpp
             webassembly/binaryen.cpp

             ${HEADERS}
             )

target_link_libraries( datxos_chain datx_utilities fc chainbase Logging IR WAST WASM Runtime
      wasm asmjs passes cfg ast emscripten-optimizer support softfloat builtins
                     )
target_include_directories( datxos_chain
                            PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/include" "${CMAKE_CURRENT_BINARY_DIR}/include"
                                   "${CMAKE_CURRENT_SOURCE_DIR}/../wasm-jit/Include"
                                   "${CMAKE_CURRENT_SOURCE_DIR}/../../externals/binaryen/src"
                            )

install( TARGETS datxos_chain
   RUNTIME DESTINATION ${CMAKE_INSTALL_FULL_BINDIR}
   LIBRARY DESTINATION ${CMAKE_INSTALL_FULL_LIBDIR}
   ARCHIVE DESTINATION ${CMAKE_INSTALL_FULL_LIBDIR}
)
install( DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/include/datxos/chain/
      DESTINATION ${CMAKE_INSTALL_FULL_INCLUDEDIR}/datxos/chain
      FILES_MATCHING PATTERN "*.hpp" PATTERN "*.h" PATTERN "webassembly" EXCLUDE
)
install(FILES ${CMAKE_CURRENT_BINARY_DIR}/include/datxos/chain/core_symbol.hpp DESTINATION ${CMAKE_INSTALL_FULL_INCLUDEDIR}/datxos/chain)
#if(MSVC)
#  set_source_files_properties( db_init.cpp db_block.cpp database.cpp block_log.cpp PROPERTIES COMPILE_FLAGS "/bigobj" )
#endif(MSVC)
