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

#删除掉原有默认钱包，创建新默认钱包
def startWallet():
    #run('rm -rf ' + os.path.abspath(args.wallet_dir))
    run('rm -rf ~/datxio-wallet' )
    #run('mkdir -p ' + os.path.abspath(args.wallet_dir))
    run('mkdir -p ~/datxio-wallet')
    print(args.kdatxd + ' --unlock-timeout %d --http-server-address http://%s:8000 ' % (unlockTimeout, args.http_server))
    background(args.kdatxd + ' --unlock-timeout %d --http-server-address http://%s:8000 --wallet-dir ~/datxio-wallet/' %( unlockTimeout, args.http_server))
    
    sleep(3)
    run(args.cldatx + 'wallet create --to-console -f ~/datxio-wallet/password.txt' )

#导入datxio的私钥，导入本地的私钥
def importKeys():
    run(args.cldatx + 'wallet import --private-key ' + accounts['pvt'])

#cldatx --wallet-url  http://127.0.0.1:8899

def startNode(nodeIndex, account):
    dir = args.nodes_dir  + account['name'] + '/'
    run('rm -rf ' + dir)
    run('mkdir -p ' + dir)
    ##################################    other p2p-peer-address     ###############################################

    otherOpts = '    --p2p-peer-address 172.31.3.30:' + str(9002)
    otherOpts = otherOpts + '    --p2p-peer-address 172.31.3.40:' + str(9002)
    otherOpts = otherOpts + '    --p2p-peer-address 172.31.3.41:' + str(9002)
    otherOpts = otherOpts + '    --p2p-peer-address 172.31.3.42:' + str(9002)

    ##################################    other p2p-peer-address     ###############################################

    if not nodeIndex: otherOpts += (
        '    --plugin datxio::history_plugin'
        '    --plugin datxio::history_api_plugin'
    )
    cmd = (
        args.noddatx +
        '    --max-irreversible-block-age -1'
        '    --contracts-console'
        '    --genesis-json ' + os.path.abspath(args.genesis) +
        '    --blocks-dir ' + os.path.abspath(dir) + '/blocks'
        '    --config-dir ' + os.path.abspath(dir) +
        '    --data-dir ' + os.path.abspath(dir) +
        '    --chain-state-db-size-mb 1024'
        ###########################   http-server-address  bnet-endpoint  p2p-listen-endpoint    #########################################

        '    --http-server-address %s:' %(args.http_server) + str(8000) +
        '    --bnet-endpoint %s:'%(args.http_server) + str(9001) +
        '    --p2p-listen-endpoint %s:'%(args.http_server) + str(9002) +

        ##########################    http-server-address  bnet-endpoint  p2p-listen-endpoint    #########################################
        '    --filter-on \"*\"'+
        '    --max-clients ' + str(maxClients) +
        '    --p2p-max-nodes-per-host ' + str(maxClients) +
        '    --enable-stale-production'
        '    --producer-name ' + account['name'] +
        '    --private-key \'["' + account['pub'] + '","' + account['pvt'] + '"]\''
        '    --plugin datxio::http_plugin'
        '    --plugin datxio::core_api_plugin'
        '    --plugin datxio::producer_plugin' +
        '    --plugin datxio::core_plugin'
        '    --plugin datxio::p2p_net_plugin' +
        '    --plugin datxio::p2p_net_api_plugin' +
        otherOpts)
    with open(dir + 'stderr', mode='w') as f:
        f.write(cmd + '\n\n')
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
    run('tail -n 60 ' + args.nodes_dir  + accounts['name'] + '/'+ 'stderr')


parser = argparse.ArgumentParser()

commands = [
    ('k', 'kill',           stepKillAll,                True,    "Kill all noddatx and kdatxd processes"),
    ('w', 'wallet',         stepStartWallet,            True,    "Start kdatxd, create wallet, fill with keys"),
    ('P', 'start-prod',     stepStartProducer,          True,    "Start producer"),
    ('l', 'log',            stepLog,                    True,    "Show tail of node's log"),
]



parser.add_argument('--public-key', metavar='', help="datxIO Public Key", default='DATX8Znrtgwt8TfpmbVpTKvA2oB8Nqey625CLN8bCN3TEbgx86Dsvr', dest="public_key")
parser.add_argument('--private-Key', metavar='', help="datxIO Private Key", default='5K463ynhZoCDDa4RDcr63cUwWLTnKqmdcoTKTHBjqoKfv4u5V7p', dest="private_key")
parser.add_argument('--cldatx', metavar='', help="Cldatx command", default='../../build/programs/cldatx/cldatx --wallet-url http://127.0.0.1:8899 ')
parser.add_argument('--noddatx', metavar='', help="Path to noddatx binary", default='../../build/programs/noddatx/noddatx')
parser.add_argument('--kdatxd', metavar='', help="Path to kdatxd binary", default='../../build/programs/kdatxd/kdatxd')
parser.add_argument('--contracts-dir', metavar='', help="Path to contracts directory", default='../../build/contracts/')
parser.add_argument('--user-limit', metavar='', help="Max number of users. (0 = no limit)", type=int, default=3)
parser.add_argument('--producer-limit', metavar='', help="Maximum number of producers. (0 = no limit)", type=int, default=3)
parser.add_argument('--max-user-keys', metavar='', help="Maximum user keys to import into wallet", type=int, default=10)
parser.add_argument('--nodes-dir', metavar='', help="Path to nodes directory", default='./nodes/')
parser.add_argument('--genesis', metavar='', help="Path to genesis.json", default="./genesis.json")
parser.add_argument('--wallet-dir', metavar='', help="Path to wallet directory", default='./wallet/')
parser.add_argument('--log-path', metavar='', help="Path to log file", default='./output.log')
parser.add_argument('--symbol', metavar='', help="The datxio.system symbol", default='SYS')
parser.add_argument('--producer-sync-delay', metavar='', help="Time (s) to sleep to allow producers to sync", type=int, default=100)
parser.add_argument('-a', '--all', action='store_true', help="Do everything marked with (*)")
parser.add_argument('--http-server',default='172.31.3.39', metavar='', help='HTTP address for cldatx')
parser.add_argument('--http-port',type=int,default=8000, metavar='', help='HTTP port for cldatx')
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

args.cldatx += '--url http://%s:%d ' %(args.http_server, args.http_port)

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
