#!/usr/bin/env python3


#######################################################################
#     
# Scrip Created by DATX Team
# For producer node join the DATX chain 
# 
# Be sure that you have private key and public key, you can find your role in accounts.json to test. 
# Fixed your private key and public key into the 0_producer_config.conf.
#
# 1. create your wallet 
#      
#    cldatx wallet create
#    
# 2. import your private-key 
#   
#    cldatx wallet import --private-key  (private key)
#
# 3. To do 
# 
#    sudo ./8_Join_DAIX.py
#    
#    This will start producer node, and the data will be saved at ./node/yourname/ .
#    
# 
#
#
#######################################################################


import util
import argparse
import json
import os


args = None

def startNode(account):
    dir = nodes_dir  + account['name'] + '/'
    util.run('rm -rf ' + dir)
    util.run('mkdir -p ' + dir)
    ##################################    other p2p-peer-address     ###############################################
    otherOpts=''
    otherOpts += '    --p2p-peer-address ' +  p2p_peer_address

    ##################################    other p2p-peer-address     ###############################################

    otherOpts += (
        '    --accessory datxos::history_accessory'
        '    --accessory datxos::history_api_accessory'
    )
    cmd = (
        'noddatx ' +
        '    --max-irreversible-block-age -1'
        '    --genesis-json ' + os.path.abspath(genesis) +
        '    --blocks-dir ' + os.path.abspath(dir) + '/blocks'
        '    --data-dir ' + os.path.abspath(dir) +
        '    --chain-state-db-size-mb 1024'+
        '    --verbose-http-errors'+
        '    --sync-fetch-span 100'+
        '    --max-transaction-time 1000'+
        '    --contracts-console'+
        '    --filter-on \"*\"'+
        '    --max-clients ' + str(maxClients) +
        '    --p2p-max-nodes-per-host ' + str(maxClients) +
        '    --enable-stale-production'
        '    --producer-name ' + account['name'] +
        '    --private-key \'["' + account['pub'] + '","' + account['pvt'] + '"]\''
        '    --accessory datxos::http_accessory'
        '    --accessory datxos::core_api_accessory'
        '    --accessory datxos::producer_accessory' +
        '    --accessory datxos::core_accessory'
        '    --accessory datxos::p2p_net_accessory' +
        '    --accessory datxos::p2p_net_api_accessory' +
        otherOpts)
    with open(dir + 'stderr', mode='w') as f:
        f.write(cmd + '\n\n')
    util.background(cmd + '    2>>' + dir + 'stderr')


def stepStartBoot():
    startNode({'name': producer_name, 'pvt': private_key, 'pub': public_key})
    util.sleep(1.5)

def stepLog():
    util.run('tail -n 20 ' + nodes_dir  + producer_name + '/'+ 'stderr')

def KillNoddatx():
    util.run('killall noddatx || true')
    util.sleep(1.5)

if __name__=="__main__":
    
    parser = argparse.ArgumentParser()
    parser.add_argument('-c','--config', type=str, default='0_producer_config.conf')
    args = parser.parse_args()

    config=util.confArg(args.config)

    producer_name=config.get('NAME','producer-name')
    private_key=config.get('KEY','private-key')
    public_key=config.get('KEY','public-key')

    genesis= config.get('PATH','genesis') 
    nodes_dir= config.get('PATH','nodes-dir')
    
    p2p_peer_address= config.get('IP','p2p-peer-address')
     
    maxClients= config.getint('PARAMETER','maxClients')

    KillNoddatx()
    stepStartBoot()
    stepLog()