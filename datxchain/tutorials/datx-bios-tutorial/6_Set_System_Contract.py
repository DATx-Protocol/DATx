#!/usr/bin/env python3


#######################################################################
#     
# Scrip Created by DATX Team
# For DATX bios node set system contract.
# 
# Be sure that you have private key and public key of datxos,fix them in 0_bios_config.conf
# Be sure that you are root user if you are bios node.
#
# 1. If:
#     
#    ./6_Set_System_Contract.py
#    
#    At first,The datxos.publickey and datxos.publickey should be repalced by real key. 
#    This will set system contract.
#    
# 
# 2. If:
#
#    ./6_Set_System_Contract.py -c your.configfile
#    
#     will read informations in your configfile.
#     
#
#######################################################################

import util
import argparse



def SetSystemContract():
    util.retry('cldatx ' + 'set contract datxos ' + contracts_dir + 'DatxSystem/ -x 3500')
    util.sleep(1)
    util.run('cldatx ' + 'push action datxos setpriv' + util.jsonArg(['datxos.msig', 1]) + '-p datxos@active')

if __name__=="__main__":
    
    parser = argparse.ArgumentParser()
    parser.add_argument('-c','--config', type=str, default='0_bios_config.conf')
    args = parser.parse_args()

    config=util.confArg(args.config)
    contracts_dir=config.get('PATH','contracts-dir')
    SetSystemContract()
