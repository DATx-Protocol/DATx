ALL


# Note:please make sure that you have installed the following dependence:
# node (v8.12.0 suggested)
# npm (we recommend that you install npm with a Node version manager(nvm) to prevent permission errors )
# go version : 1.10 or more
# python version: 2.7 
===========================================================

# setup system environment variable: $GOPATH
# into project path and compile all software

sudo ./datx_all.sh build

# install all software

sudo ./datx_all.sh install

# setup noddatx config.ini with producer-name and start noddatx,'-p' is not allowed 
# Note:  you must start noddatx first to generate the default config file before start other software.

## config.ini can be found at the following locations:
#       Mac OS: ~/Library/Application Support/datxos/noddatx/config/config.ini
#       Linux: ~/.local/share/datxos/noddatx/config/config.ini

# we do not support custom config.ini file.
# Note : Setup your producer-name in the config.ini
# Note : Setup your wallet-namer and wallet-password  in the ~/datxos-wallet/wallet_password.ini
noddatx --accessory  datxos::core_api_accessory datxos::history_accessory datxos::history_api_accessory --replay-blockchain --verbose-http-errors

# start all other software
sudo ./datx_all.sh start

# stop all other software
sudo ./datx_all.sh stop

