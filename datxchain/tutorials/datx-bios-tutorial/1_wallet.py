#!/usr/bin/env python3


#######################################################################
#     
# Scrip Created by DATX Team
# For DATX wallet import 
# 
# Be sure that you have a private key for datxos account, fix it in 0_bios_config.conf
# Be sure that you are root user if you are datxos.
#
# 1. if :
#
#    ./1_wallet.py 
#     
#    will kill another noddatx and kdatxd firstly,and then it will create a default wallet in ~/datxos-wallet, and import your private key
#    your wallet password will be saved in ~/datxos-wallet/password.ini.
# 
#
# 2. if :
#
#    ./1_wallet.py -c your.configfile

#    will read informations in your configfile.
#    will just create a default wallet in ~/datxos-wallet, and import your private key.
#    Similarly,your wallet password will be saved in ~/datxos-wallet/password.ini
#
#######################################################################


import util
import argparse
import json



def StartWallet():
    util.run('rm -rf ~/datxos-wallet' )
    util.run('mkdir -p ~/datxos-wallet')
    util.sleep(3)
    util.run('cldatx ' + 'wallet create --file ~/datxos-wallet/password.ini' )

def ImportKeys(private_key):
    
    util.run('cldatx ' + 'wallet import --private-key ' + private_key)

    keys = {}
    for a in accounts:
        key = a['pvt']
        if not key in keys:
            if len(keys) >= max_user_keys:
                break
            keys[key] = True
            util.run('cldatx wallet import --private-key ' + key)

def KillAll():
    util.run('killall kdatxd noddatx || true')
    util.sleep(1.5)




if __name__=="__main__":
    
    parser = argparse.ArgumentParser()
    parser.add_argument('-c','--config', type=str, default='0_bios_config.conf')
    args = parser.parse_args()

    config=util.confArg(args.config)
    private_key=config.get('KEY','private-key')
    log_path=config.get('PATH','log-path')
    user_limit=config.getint('PARAMETER','user-limit')
    max_user_keys=config.getint('PARAMETER','max-user-keys')

    with open(log_path,'a') as logFile:
        logFile.write('\n\n' + '*' * 80 + '\n\n\n')
    

    with open('accounts.json') as f:
        a = json.load(f)
        if user_limit:
            del a['users'][user_limit:]#10

        accounts = a['users']

    
    KillAll()
    StartWallet()
    ImportKeys(private_key)





