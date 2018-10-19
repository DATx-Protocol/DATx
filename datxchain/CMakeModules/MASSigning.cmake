macro(mas_sign target)


if(APPLE AND MAS_CERT_FINGERPRINT AND MAS_BASE_APPID AND MAS_PROVISIONING_PROFILE AND MAS_KEYCHAIN_GROUP)

  add_custom_command(TARGET ${target} POST_BUILD
                     COMMAND ${CMAKE_SOURCE_DIR}/tools/mas_sign.sh ${MAS_CERT_FINGERPRINT} ${MAS_BASE_APPID}$<TARGET_FILE_NAME:${target}> ${MAS_PROVISIONING_PROFILE} ${MAS_KEYCHAIN_GROUP} $<TARGET_FILE_NAME:${target}>
                     WORKING_DIRECTORY ${CMAKE_CURRENT_BINARY_DIR}
                     VERBATIM
                     )

endif()

endmacro(mas_sign)
