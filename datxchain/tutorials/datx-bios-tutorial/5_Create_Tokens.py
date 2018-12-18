#!/usr/bin/env python3


#######################################################################
#     
# Scrip Created by DATX Team
# For DATX bios node create tokens.
# 
# Be sure that you have private key and public key of datxos,fix them in 0_bios_config.conf
# Be sure that you are root user if you are bios node.
#
# 1. If:
#     
#    ./5_Create_Tokens.py 
#    
#    At first,The datxos.publickey and datxos.publickey should be repalced by real key. 
#    This will create tokens.
#    
# 
# 2. If:
#
#    ./5_Create_Tokens.py -c your.configfile
#    
#     will read informations in your configfile.
#     
#
#######################################################################

import util
import argparse

def intToCurrency(i):
    return '%d.%04d %s' % (i // 10000, i % 10000, symbol)

def CreateTokens():
    #datx
    util.run('cldatx ' + 'push action datxos.token create \'["datxos", "10000000000.0000 %s"]\' -p datxos.token' % (symbol)+' -x 3500')
    #totalAllocation = allocateFunds(0, len(accounts))   
    
    util.run('cldatx ' + 'push action datxos.token issue \'["datxos", "%s", "memo"]\' -p datxos' % intToCurrency(totalAllocation))

    #btc
    util.run('cldatx ' + 'push action datxos.dtoke create \'["datxos.dbtc", "21000000.0000 DBTC",0,0,0]\' -p datxos.dtoke' )
    util.run('cldatx ' + 'push action datxos.dtoke issue \'["datxos.dbtc", "21000000.0000 DBTC", "memo"]\' -p datxos.dbtc')
    
    #util.run('cldatx '+ 'push action datxos.dtoke transfer \'{"from":"datxos.dtoke","to":"alice","quantity":"100.0000 DBTC","memo":"test"}\' -p datxos.dtoke')
    
    
    #eth
    util.run('cldatx ' + 'push action datxos.dtoke create \'["datxos.deth", "102000000.0000 DETH",0,0,0]\' -p datxos.dtoke' )
    util.run('cldatx ' + 'push action datxos.dtoke issue \'["datxos.deth", "102000000.0000 DETH", "memo"]\' -p datxos.deth')
    #eos
    util.run('cldatx ' + 'push action datxos.dtoke create \'["datxos.deos", "1000000000.0000 DEOS",0,0,0]\' -p datxos.dtoke' )
    util.run('cldatx ' + 'push action datxos.dtoke issue \'["datxos.deos", "1000000000.0000 DEOS", "memo"]\' -p datxos.deos')
    
    util.run('cldatx ' + 'set account permission datxos.dbtc active \'{"threshold": 1,"keys": [{"key": "DATX8Znrtgwt8TfpmbVpTKvA2oB8Nqey625CLN8bCN3TEbgx86Dsvr","weight": 1}],"accounts": [{"permission":{"actor":"datxos.charg","permission":"datxos.code"},"weight":1},{"permission":{"actor":"datxos.extra","permission":"datxos.code"},"weight":1}]}\' owner -p datxos.dbtc')
    util.run('cldatx ' + 'set account permission datxos.deth active \'{"threshold": 1,"keys": [{"key": "DATX8Znrtgwt8TfpmbVpTKvA2oB8Nqey625CLN8bCN3TEbgx86Dsvr","weight": 1}],"accounts": [{"permission":{"actor":"datxos.charg","permission":"datxos.code"},"weight":1},{"permission":{"actor":"datxos.extra","permission":"datxos.code"},"weight":1}]}\' owner -p datxos.deth')
    util.run('cldatx ' + 'set account permission datxos.deos active \'{"threshold": 1,"keys": [{"key": "DATX8Znrtgwt8TfpmbVpTKvA2oB8Nqey625CLN8bCN3TEbgx86Dsvr","weight": 1}],"accounts": [{"permission":{"actor":"datxos.charg","permission":"datxos.code"},"weight":1},{"permission":{"actor":"datxos.extra","permission":"datxos.code"},"weight":1}]}\' owner -p datxos.deos')
    
    util.sleep(1)



if __name__=="__main__":
    
    parser = argparse.ArgumentParser()
    parser.add_argument('-c','--config', type=str, default='0_bios_config.conf')
    args = parser.parse_args()

    config=util.confArg(args.config)
    totalAllocation=76000000000000
    symbol=config.get('PARAMETER','symbol')
    CreateTokens()
