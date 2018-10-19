#! /bin/bash

if [ -d "/usr/local/include/datxos" ]; then
   printf "\n\tOld datxos install needs to be removed.\n\n"
   printf "\tDo you wish to remove this install? (requires sudo)\n"
   select yn in "Yes" "No"; do
      case $yn in
         [Yy]* )
            if [ "$(id -u)" -ne 0 ]; then
               printf "\n\tThis requires sudo, please run ./scripts/clean_old_install.sh with sudo\n\n"
               exit -1
            fi
            pushd /usr/local &> /dev/null

            pushd include &> /dev/null
            rm -rf appbase chainbase datxos datxos.system DatxosLib fc libc++ musl &> /dev/null
            popd &> /dev/null

            pushd bin &> /dev/null
            rm cldatx datxos-abigen datxos-applesdemo datxos-launcher datxos-s2wasm datxos-wast2wasm datxoscpp kdatxd noddatx &> /dev/null
            popd &> /dev/null

            libraries=(libdatxos_testing
            libdatxos_chain
            libfc
            libbinaryen
            libWAST
            libWASM
            libRuntime
            libPlatform
            libIR
            libLogging
            libsoftfloat
            libchainbase
            libappbase
            libbuiltins)
            pushd lib &> /dev/null
            for lib in ${libraries[@]}; do
               rm ${lib}.a ${lib}.dylib ${lib}.so &> /dev/null
            done
            popd &> /dev/null

            pushd etc &> /dev/null
            rm datxos &> /dev/null
            popd &> /dev/null

            pushd share &> /dev/null
            rm datxos &> /dev/null
            popd &> /dev/null

            pushd usr/share &> /dev/null
            rm datxos &> /dev/null
            popd &> /dev/null

            pushd var/lib &> /dev/null
            rm datxos &> /dev/null
            popd &> /dev/null

            pushd var/log &> /dev/null
            rm datxos &> /dev/null
            popd &> /dev/null

            popd &> /dev/null
            break;;
         [Nn]* )
            printf "\tAborting uninstall\n\n"
            exit -1;;
      esac
   done
fi
