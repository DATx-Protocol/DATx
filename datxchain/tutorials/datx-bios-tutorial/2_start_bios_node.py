#!/usr/bin/env python3


#######################################################################
#     
# Scrip Created by DATX Team
# For DATX bios node start 
# 
# Be sure that you have private key and public key of datxos,fix them in 0_bios_config.conf
# Be sure that you are root user if you are bios node.
#
# 1. If:
#     
#    ./2_start_bios_node.py 
#    
#    At first,The datxos.publickey and datxos.publickey should be repalced by real key. 
#    This will start datxos node,and the data will be saved at ./node/datxos/ .
#    
# 
# 2. If:
#
#    ./2_start_bios_node.py -c your.configfile
#    
#     will read informations in your configfile.
#     This will start datxos node,and the data will be saved at ./node/datxos/
#
#######################################################################


import util
import argparse
import json
import os


args = None

def startNode(account):
    dir = nodes_dir + account['name'] + '/'
    util.run('rm -rf ' + dir)
    util.run('mkdir -p ' + dir)
    otherOpts =''
    otherOpts += '    --accessory datxos::history_accessory'
    otherOpts += '    --accessory datxos::history_api_accessory' 
    cmd = (
        'noddatx ' +
        '    --max-irreversible-block-age -1'
        '    --genesis-json ' + os.path.abspath(genesis) +
        '    --blocks-dir ' + os.path.abspath(dir) + '/blocks'+
        '    --data-dir ' + os.path.abspath(dir) +
        '    --chain-state-db-size-mb 1024'+
        '    --verbose-http-errors'+
        '    --contracts-console'+
        '    --max-transaction-time 1000'
        '    --max-clients ' + str(maxClients) +
        '    --p2p-max-nodes-per-host ' + str(maxClients) +
        '    --enable-stale-production '
        '    --filter-on \"*\"'+
        '    --producer-name ' + account['name'] +
        '    --private-key \'["' + account['pub'] + '","' + account['pvt'] + '"]\''
        '    --accessory datxos::http_accessory'
        '    --accessory datxos::core_api_accessory'
        '    --accessory datxos::producer_accessory' +
        '    --accessory datxos::core_accessory'+
        '    --accessory datxos::p2p_net_accessory' +
        '    --accessory datxos::p2p_net_api_accessory' +
        otherOpts)
    with open(dir + 'stderr', mode='w') as f:
        f.write(cmd + '\n\n')
    util.background(cmd + '    2>>' + dir + 'stderr')


def stepStartBoot():
    startNode({'name': 'datxos', 'pvt': private_key, 'pub': public_key})
    util.sleep(1.5)

def stepLog():
    util.run('tail -n 20 ' + nodes_dir + 'datxos' + '/'+ 'stderr')

def KillNoddatx():
    util.run('killall noddatx || true')
    util.sleep(1.5)

if __name__=="__main__":
    
    parser = argparse.ArgumentParser()
    parser.add_argument('-c','--config', type=str, default='0_bios_config.conf')
    args = parser.parse_args()

    config=util.confArg(args.config)
    private_key=config.get('KEY','private-key')
    public_key=config.get('KEY','public-key')

    genesis= config.get('PATH','genesis') 
    nodes_dir= config.get('PATH','nodes-dir')

    maxClients= config.getint('PARAMETER','maxClients')
    
    KillNoddatx()
    stepStartBoot()
    stepLog()