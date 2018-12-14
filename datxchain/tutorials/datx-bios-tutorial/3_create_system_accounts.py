#!/usr/bin/env python3


#######################################################################
#     
# Scrip Created by DATX Team
# For DATX bios node create system accounts.
# 
# Be sure that you have private key and public key of datxos,fix them in 0_bios_config.conf
# Be sure that you are root user if you are bios node.
#
# 1. If:
#     
#    ./3_create_system_accounts.py 
#    
#    At first,The datxos.publickey and datxos.publickey should be repalced by real key. 
#    This will create accounts as follows:
#    'datxos.bpay',
#    'datxos.msig',
#    'datxos.names',
#    'datxos.ram',
#    'datxos.save',
#    'datxos.stake',
#    'datxos.token',
#    'datxos.vpay',
#    'datxos.veri',
#    'datxos.charg',
#    'datxos.deth',
#    'datxos.deos',
#    'datxos.dbtc',
#    'datxos.extra',
#    'datxos.dtoke'
#    
# 
# 2. If:
#
#    ./3_create_system_accounts.py -c your.configfile
#    
#     will read informations in your configfile.
#     
#
#######################################################################
import util
import argparse

systemAccounts = [
    'datxos.bpay',
    'datxos.msig',
    'datxos.names',
    'datxos.ram',
    'datxos.save',
    'datxos.stake',
    'datxos.token',
    'datxos.vpay',
    'datxos.veri',
    'datxos.charg',
    'datxos.deth',
    'datxos.deos',
    'datxos.dbtc',
    'datxos.extra',
    'datxos.dtoke'
]

def createSystemAccounts():
    for a in systemAccounts:
        util.run('cldatx ' + 'create account datxos ' + a + ' ' + public_key)


if __name__=="__main__":
    
    parser = argparse.ArgumentParser()
    parser.add_argument('-c','--config', type=str, default='0_bios_config.conf')
    args = parser.parse_args()

    config=util.confArg(args.config)
    public_key=config.get('KEY','public-key')

    createSystemAccounts()
