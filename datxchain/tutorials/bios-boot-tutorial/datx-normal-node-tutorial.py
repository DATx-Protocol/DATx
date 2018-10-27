#!/usr/bin/env python3

import argparse
import json
import numpy
import os
import random
import re
import subprocess
import sys
import time

args = None
logFile = None

unlockTimeout = 999999999

http_server=""



def jsonArg(a):
    return " '" + json.dumps(a) + "' "

def run(args):
    print('datx-bios-tutorial.py:', args)
    logFile.write(args + '\n')
    if subprocess.call(args, shell=True):
        print('datx-bios-tutorial.py: exiting because of error')
        sys.exit(1)

def retry(args):
    while True:
        print('datx-bios-tutorial.py:', args)
        logFile.write(args + '\n')
        if subprocess.call(args, shell=True):
            print('*** Retry')
        else:
            break

def background(args):
    print('datx-bios-tutorial.py:', args)
    logFile.write(args + '\n')
    return subprocess.Popen(args, shell=True)

def getOutput(args):
    print('datx-bios-tutorial.py:', args)
    logFile.write(args + '\n')
    proc = subprocess.Popen(args, shell=True, stdout=subprocess.PIPE)
    return proc.communicate()[0].decode('utf-8')

def getJsonOutput(args):
    print('datx-bios-tutorial.py:', args)
    logFile.write(args + '\n')
    proc = subprocess.Popen(args, shell=True, stdout=subprocess.PIPE)
    return json.loads(proc.communicate()[0])

def sleep(t):
    print('sleep', t, '...')
    time.sleep(t)
    print('resume')

#删除掉上次创建钱包，创建新钱包
def startWallet():
    #run('rm -rf ' + os.path.abspath(args.wallet_dir))
    #run('mkdir -p ' + os.path.abspath(args.wallet_dir))
    #run('rm -rf ~/datxos-wallet' )
    #run('mkdir -p ~/datxos-wallet')

    home_dir=os.environ['HOME']
    if os.path.exists(home_dir+'/datxos-wallet'):
        if os.path.exists(home_dir+'/datxos-wallet/my_wallet.wallet'):
            run('rm  ~/datxos-wallet/my_wallet.wallet' )
            
        if os.path.exists(home_dir+'/datxos-wallet/my_wallet_password.txt'):    
            run('rm  ~/datxos-wallet/my_wallet_password.txt' )
    else:
         
        run('mkdir -p ~/datxos-wallet')


    
    #background(args.kdatxd + ' --unlock-timeout %d --http-server-address http://%s:8888 --wallet-dir ~/datxos-wallet/' %( unlockTimeout, http_server))
    #background(args.kdatxd + ' --unlock-timeout %d --wallet-dir ~/datxos-wallet/' %( unlockTimeout, http_server))


    sleep(3)
    run(args.cldatx + 'wallet create -n my_wallet --file ~/datxos-wallet/my_wallet_password.txt' )

#导入datxos的私钥，导入本地的私钥
def importKeys():
    run(args.cldatx + 'wallet import -n my_wallet --private-key ' + accounts['pvt'])



def startNode(nodeIndex, account):
    dir = args.nodes_dir  + account['name'] + '/'
    run('rm -rf ' + dir)
    run('mkdir -p ' + dir)
    ##################################    other p2p-peer-address     ###############################################
    otherOpts=''
    otherOpts = '    --p2p-peer-address 172.31.3.5:' + str(9876)
    #otherOpts = otherOpts + '    --p2p-peer-address 172.31.3.39:' + str(9002)

    ##################################    other p2p-peer-address     ###############################################

    if not nodeIndex: otherOpts += (
        '    --accessory datxos::history_accessory'
        '    --accessory datxos::history_api_accessory'
    )
    cmd = (
        args.noddatx +
        '    --max-irreversible-block-age -1'
        '    --contracts-console'
        '    --genesis-json ' + os.path.abspath(args.genesis) +
        '    --blocks-dir ' + os.path.abspath(dir) + '/blocks'
        #'    --config-dir ' + os.path.abspath(dir) +
        '    --data-dir ' + os.path.abspath(dir) +
        '    --chain-state-db-size-mb 1024'+
        '    --verbose-http-errors'+
        '    --sync-fetch-span 100'+
        '    --max-transaction-time 1000'+
        ###########################   http-server-address  bnet-endpoint  p2p-listen-endpoint    #########################################

        #'    --http-server-address %s:' %(http_server) + str(8888) +
        #'    --bnet-endpoint %s:'%(http_server) + str(9001) +
        #'    --p2p-listen-endpoint %s:'%(http_server) + str(9002) +

        # '    --http-server-address '  +
        # '    --bnet-endpoint ' +
        # '    --p2p-listen-endpoint ' +

        ##########################    http-server-address  bnet-endpoint  p2p-listen-endpoint    #########################################
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
    #background(cmd + '    2>>' + dir + 'stderr')
    background(cmd + '    2>>' + dir + 'stderr')
    
def listProducers():
    run(args.cldatx + 'system listproducers')


def stepKillAll():
    run('killall kdatxd noddatx || true')
    sleep(1.5)
def stepStartWallet():
    startWallet()
    importKeys()
def stepStartProducer():
    startNode(0, {'name': accounts["name"], 'pvt': accounts['pvt'], 'pub': accounts['pub']})
    sleep(1.5)
def stepLog():
    run('tail -n 20 ' + args.nodes_dir  + accounts['name'] + '/'+ 'stderr')


parser = argparse.ArgumentParser()

commands = [
    ('k', 'kill',           stepKillAll,                True,    "Kill all noddatx and kdatxd processes"),
    ('w', 'wallet',         stepStartWallet,            True,    "Start kdatxd, create wallet, fill with keys"),
    ('P', 'start-prod',     stepStartProducer,          True,    "Start producer"),
    ('l', 'log',            stepLog,                    True,    "Show tail of node's log"),
]



parser.add_argument('--public-key', metavar='', help="datxOS Public Key", default='DATX8Znrtgwt8TfpmbVpTKvA2oB8Nqey625CLN8bCN3TEbgx86Dsvr', dest="public_key")
parser.add_argument('--private-Key', metavar='', help="datxOS Private Key", default='5K463ynhZoCDDa4RDcr63cUwWLTnKqmdcoTKTHBjqoKfv4u5V7p', dest="private_key")
parser.add_argument('--cldatx', metavar='', help="Cldatx command", default='cldatx  ')
parser.add_argument('--noddatx', metavar='', help="Path to noddatx binary", default='noddatx')
parser.add_argument('--kdatxd', metavar='', help="Path to kdatxd binary", default='kdatxd')
parser.add_argument('--contracts-dir', metavar='', help="Path to contracts directory", default='../../build/contracts/')
parser.add_argument('--user-limit', metavar='', help="Max number of users. (0 = no limit)", type=int, default=3)
parser.add_argument('--producer-limit', metavar='', help="Maximum number of producers. (0 = no limit)", type=int, default=3)
parser.add_argument('--max-user-keys', metavar='', help="Maximum user keys to import into wallet", type=int, default=10)
parser.add_argument('--nodes-dir', metavar='', help="Path to nodes directory", default='./nodes/')
parser.add_argument('--genesis', metavar='', help="Path to genesis.json", default="./genesis.json")
parser.add_argument('--wallet-dir', metavar='', help="Path to wallet directory", default='./wallet/')
parser.add_argument('--log-path', metavar='', help="Path to log file", default='./output.log')
parser.add_argument('--symbol', metavar='', help="The datxos.system symbol", default='SYS')
parser.add_argument('--producer-sync-delay', metavar='', help="Time (s) to sleep to allow producers to sync", type=int, default=100)
parser.add_argument('-a', '--all', action='store_true', help="Do everything marked with (*)")
#parser.add_argument('--http-server',default='127.0.0.1', metavar='', help='HTTP address for cldatx')
parser.add_argument('--http-port',type=int,default=8888, metavar='', help='HTTP port for cldatx')
#parser.add_argument('--producer-name',default='producer111a', metavar='', help='default producer name.')


for (flag, command, function, inAll, help) in commands:
    prefix = ''
    if inAll: prefix += '*'
    if prefix: help = '(' + prefix + ') ' + help
    if flag:
        parser.add_argument('-' + flag, '--' + command, action='store_true', help=help, dest=command)
    else:
        parser.add_argument('--' + command, action='store_true', help=help, dest=command)


args = parser.parse_args()

#args.cldatx += '--url http://%s:%d ' %(args.http_server, args.http_port)

logFile = open(args.log_path, 'a')

logFile.write('\n\n' + '*' * 80 + '\n\n\n')

maxClients =  50

with open('my_account.json') as f:
    accounts = json.load(f)
    #accounts = a['producers']


haveCommand = False
for (flag, command, function, inAll, help) in commands:
    if getattr(args, command) or inAll and args.all:
        if function:
            haveCommand = True
            function()
if not haveCommand:
    print('bios-boot-tutorial.py: Tell me what to do. -a does almost everything. -h shows options.')
