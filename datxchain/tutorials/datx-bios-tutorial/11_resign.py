#!/usr/bin/env python3


#######################################################################
#     
# Scrip Created by DATX Team
# For DATX bios node resign auth to producers.
# 
# Be sure that you have private key and public key of datxos,fix them in 0_bios_config.conf
# Be sure that you are root user if you are bios node.
#
# 1. If:
#     
#    ./11_resign.py
#    
#    will resign auth of datxos and another system accounts to producers.
#     
#
#######################################################################

import util
import json


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



def updateAuth(account, permission, parent, controller):
    util.run('cldatx ' + 'push action datxos updateauth' + util.jsonArg({
        'account': account,
        'permission': permission,
        'parent': parent,
        'auth': {
            'threshold': 1, 'keys': [], 'waits': [],
            'accounts': [{
                'weight': 1,
                'permission': {'actor': controller, 'permission': 'active'}
            }]
        }
    }) + '-p ' + account + '@' + permission)

def updateAuth_spe(account, permission, parent, controller):
    util.run('cldatx ' + 'push action datxos updateauth' + util.jsonArg({
        'account': account,
        'permission': permission,
        'parent': parent,
        'auth': {
            'threshold': 1, 'keys': [], 'waits': [],
            'accounts': [{
                'weight': 1,
                'permission': {'actor': controller, 'permission': 'active'}
            },{
               "weight":1,
               "permission":{"actor":"datxos.charg","permission":"datxos.code"}
            },{
               "weight":1,
               "permission":{"actor":"datxos.extra","permission":"datxos.code"}
            }]
        }
    }) + '-p ' + account + '@' + permission)

def updateAuth_extra(account, permission, parent, controller):
    util.run('cldatx ' + 'push action datxos updateauth' + util.jsonArg({
        'account': account,
        'permission': permission,
        'parent': parent,
        'auth': {
            'threshold': 1, 'keys': [], 'waits': [],
            'accounts': [{
                'weight': 1,
                'permission': {'actor': controller, 'permission': 'active'}
            }, {
                "weight": 1,
                "permission": {"actor": "datxos.extra", "permission": "datxos.code"}
            }]
        }
    }) + '-p ' + account + '@' + permission)

def resign(account, controller):
    list_con=['datxos.dbtc','datxos.deos','datxos.deth']
    if account in list_con:
        updateAuth_spe(account, 'owner', '', controller)
        updateAuth_spe(account, 'active', 'owner', controller)
    elif account == "datxos.extra":
        updateAuth_extra(account, 'owner', '', controller)
        updateAuth_extra(account, 'active', 'owner', controller)  
    else:
        updateAuth(account, 'owner', '', controller)
        updateAuth(account, 'active', 'owner', controller)

    util.sleep(1)
    util.run('cldatx ' + 'get account ' + account)


if __name__=="__main__":
    
    resign('datxos', 'datxos.prods')

    for a in systemAccounts:
        resign(a, 'datxos')