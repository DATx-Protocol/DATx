#!/usr/bin/env python3


#######################################################################
#     
# Scrip Created by DATX Team
# For DATX bios node install contracts.
# 
# Be sure that you have private key and public key of datxos,fix them in 0_bios_config.conf
# Be sure that you are root user if you are bios node.
#
# 1. If:
#     
#    ./4_Install_Contracts.py 
#    
#    At first,The datxos.publickey and datxos.publickey should be repalced by real key. 
#    This will create install contracts.
#    
# 
# 2. If:
#
#    ./4_Install_Contracts.py -c your.configfile
#    
#     will read informations in your configfile.
#     This will create install contracts.
#
#######################################################################
import util
import argparse


def InstallContracts():
    util.run('cldatx ' + 'set contract datxos.token ' + contracts_dir + 'DatxToken/')
    util.run('cldatx ' + 'set contract datxos.msig ' + contracts_dir + 'DatxMsig/')
    util.run('cldatx ' + 'set contract datxos.charg ' + contracts_dir +'DatxRecharge/')

    util.run('cldatx ' + 'set contract datxos.dtoke ' + contracts_dir + 'DatxDToken/')

    util.run('cldatx ' + 'set contract datxos.extra ' + contracts_dir + 'DatxExtract/')

    util.run('cldatx ' + 'set account permission datxos.dtoke active \'{"threshold": 1,"keys": [{"key": "DATX8Znrtgwt8TfpmbVpTKvA2oB8Nqey625CLN8bCN3TEbgx86Dsvr","weight": 1}],"accounts": [{"permission":{"actor":"datxos.charg","permission":"datxos.code"},"weight":1}]}\' owner -p datxos.dtoke')
    
    util.run('cldatx ' + 'set account permission datxos.dtoke active \'{"threshold": 1,"keys": [{"key": "DATX8Znrtgwt8TfpmbVpTKvA2oB8Nqey625CLN8bCN3TEbgx86Dsvr","weight": 1}],"accounts": [{"permission":{"actor":"datxos.extra","permission":"datxos.code"},"weight":1}]}\' owner -p datxos.dtoke')


if __name__=="__main__":
    
    parser = argparse.ArgumentParser()
    parser.add_argument('-c','--config', type=str, default='0_bios_config.conf')
    args = parser.parse_args()

    config=util.confArg(args.config)
    contracts_dir=config.get('PATH','contracts-dir')
    public_key=config.get('KEY','public-key')
    InstallContracts()
