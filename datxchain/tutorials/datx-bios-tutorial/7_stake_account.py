#!/usr/bin/env python3


#######################################################################
#     
# Scrip Created by DATX Team
# For DATX bios node create stake account.
# 
# Be sure that you have private key and public key of datxos,fix them in 0_bios_config.conf
# Be sure that you are root user if you are bios node.
#
# 1. If:
#     
#    ./7_stake_account.py
#    
#    At first,The datxos.publickey and datxos.publickey should be repalced by real key. 
#    This will create stake account.
#    
# 
# 2. If:
#
#    ./7_stake_account.py -c your.configfile
#    
#     will read informations in your configfile.
#     
#
#######################################################################

import util
import argparse
import json


def intToCurrency(i):
    return '%d.%04d %s' % (i // 10000, i % 10000, symbol)


def createStakedAccounts(begin, end):
    #(b,e) =(0,len(accounts))=(0,7),(0~3)users,(4~7)producers
    for i in range(begin, end):
        a = accounts[i]

        stakeNet = 500000000000  #50000000DATX
        stakeCpu = 500000000000  #50000000DATX
        stakeRam = 500000000000  #50000000DATX
        small_stake=5000000 #500DATX

        util.retry('cldatx ' + 'system newaccount --transfer datxos %s %s --stake-net "%s" --stake-cpu "%s" --buy-ram "%s"   ' % 
            (a['name'], a['pub'], intToCurrency(stakeNet), intToCurrency(stakeCpu), intToCurrency(stakeRam)))
        util.sleep(1)
    
    for i in range(user_limit):
        a = accounts[i]
        util.retry('cldatx ' + 'push action datxos.token transfer \'["datxos",%s,"%s","vote"]\' -p datxos'% (a['name'],intToCurrency(5000000000000)))
        util.sleep(1)
        util.retry('cldatx ' + 'system delegatebw datxos %s "%s" "%s" --transfer'%(a['name'],intToCurrency(300000000000),intToCurrency(350000000000)))
        util.sleep(1)

    util.run('cldatx '+ 'push action datxos.dtoke transfer \'{"from":"datxos.dbtc","to":"alice","quantity":"300.0000 DBTC","memo":"test"}\' -p datxos.dbtc')
    util.run('cldatx '+ 'push action datxos.dtoke transfer \'{"from":"datxos.deth","to":"alice","quantity":"300.0000 DETH","memo":"test"}\' -p datxos.deth')
    util.run('cldatx '+ 'push action datxos.dtoke transfer \'{"from":"datxos.deos","to":"alice","quantity":"300.0000 DEOS","memo":"test"}\' -p datxos.deos')


if __name__=="__main__":
    
    parser = argparse.ArgumentParser()
    parser.add_argument('-c','--config', type=str, default='0_bios_config.conf')
    args = parser.parse_args()

    config=util.confArg(args.config)
    symbol=config.get('PARAMETER','symbol')
    contracts_dir=config.get('PATH','contracts-dir')
    user_limit=config.getint('PARAMETER','user-limit')
    producer_limit=config.getint('PARAMETER','producer-limit')


    with open('accounts.json') as f:
        a = json.load(f)
        if user_limit:
            del a['users'][user_limit:]#10
        if producer_limit:
            del a['producers'][producer_limit:]#3
        firstProducer = len(a['users'])
        numProducers = len(a['producers'])
        accounts = a['users'] + a['producers']

    createStakedAccounts(0, len(accounts))