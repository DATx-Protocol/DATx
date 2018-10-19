macro(datxos_additional_accessory)
    set(ADDITIONAL_accessoryS_TARGET "${ADDITIONAL_accessoryS_TARGET};${ARGN}" PARENT_SCOPE)
endmacro()

foreach(ADDITIONAL_accessory_SOURCE_DIR ${DATXOS_ADDITIONAL_accessoryS})
    string(UUID ADDITIONAL_accessory_SOURCE_DIR_MD5 NAMESPACE "00000000-0000-0000-0000-000000000000" NAME ${ADDITIONAL_accessory_SOURCE_DIR} TYPE MD5)
    message(STATUS "[Additional Accessory] ${ADDITIONAL_accessory_SOURCE_DIR} => ${CMAKE_BINARY_DIR}/additional_accessorys/${ADDITIONAL_accessory_SOURCE_DIR_MD5}")
    add_subdirectory(${ADDITIONAL_accessory_SOURCE_DIR} ${CMAKE_BINARY_DIR}/additional_accessorys/${ADDITIONAL_accessory_SOURCE_DIR_MD5})
endforeach()

foreach(ADDITIONAL_accessory_TARGET ${ADDITIONAL_accessoryS_TARGET})
    target_link_libraries( noddatx PRIVATE -Wl,${whole_archive_flag} ${ADDITIONAL_accessory_TARGET} -Wl,${no_whole_archive_flag} )
endforeach()
