add_subdirectory(bp2p_accessory)
add_subdirectory(p2p_net_accessory)
add_subdirectory(p2p_net_api_accessory)
add_subdirectory(http_accessory)
add_subdirectory(http_client_accessory)
add_subdirectory(core_accessory)
add_subdirectory(core_api_accessory)
add_subdirectory(producer_accessory)
add_subdirectory(producer_api_accessory)
add_subdirectory(history_accessory)
add_subdirectory(history_api_accessory)

add_subdirectory(wallet_accessory)
add_subdirectory(wallet_api_accessory)
#add_subdirectory(txn_test_gen_accessory)
add_subdirectory(db_size_api_accessory)
#add_subdirectory(faucet_testp2p_net_accessory)
add_subdirectory(mongo_db_accessory)
add_subdirectory(login_accessory)

# Forward variables to top level so packaging picks them up
set(CPACK_DEBIAN_PACKAGE_DEPENDS ${CPACK_DEBIAN_PACKAGE_DEPENDS} PARENT_SCOPE)
