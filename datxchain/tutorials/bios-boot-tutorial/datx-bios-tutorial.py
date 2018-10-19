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
wait_other_node_time=10
push_action_times=20

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
    run('rm -rf ~/datxos-wallet' )
    #run('mkdir -p ' + os.path.abspath(args.wallet_dir))
    run('mkdir -p ~/datxos-wallet')
    #print(args.kdatxd + ' --unlock-timeout %d --http-server-address http://192.168.152.135:8888 --wallet-dir %s' % (unlockTimeout, os.path.abspath(args.wallet_dir)))
    background(args.kdatxd + ' --unlock-timeout %d --http-server-address http://%s:8888 --wallet-dir ~/datxos-wallet/' % (unlockTimeout,args.http_server))
    
    sleep(3)
    run(args.cldatx + 'wallet create --file ~/datxos-wallet/password.txt' ) #--to-console

#导入datxos的私钥，导入本地的私钥
def importKeys():
    run(args.cldatx + 'wallet import --private-key ' + args.private_key)
    keys = {}
    for a in accounts:
        key = a['pvt']
        if not key in keys:
            if len(keys) >= args.max_user_keys:
                break
            keys[key] = True
            run(args.cldatx + 'wallet import --private-key ' + key)
    for i in range(firstProducer, firstProducer + numProducers):
        a = accounts[i]
        key = a['pvt']
        if not key in keys:
            keys[key] = True
            run(args.cldatx + 'wallet import --private-key ' + key)

#cldatx --wallet-url  http://127.0.0.1:8899

def startNode(nodeIndex, account):
    dir = args.nodes_dir + ('%02d-' % nodeIndex) + account['name'] + '/'
    run('rm -rf ' + dir)
    run('mkdir -p ' + dir)
    ##################################     p2p-peer-address     ###############################################
    otherOpts =''
    #otherOpts = otherOpts +'    --p2p-peer-address 172.31.3.39:' + str(9002)
    #otherOpts = otherOpts + '    --p2p-peer-address 172.31.3.170:' + str(9002)
    #otherOpts = otherOpts + '    --p2p-peer-address 172.31.7.68:' + str(9002)
    #otherOpts = otherOpts + '    --p2p-peer-address 192.168.152.139:' + str(9002)

    ##################################     p2p-peer-address     ###############################################

    #if not nodeIndex: otherOpts += (
    #    '    --accessory datxos::history_accessory'
    #    '    --accessory datxos::history_api_accessory'
    #)
    otherOpts += '    --accessory datxos::history_accessory'
    otherOpts += '    --accessory datxos::history_api_accessory' 
        
    
    cmd = (
        args.noddatx +
        '    --max-irreversible-block-age -1'
        '    --contracts-console'
        '    --genesis-json ' + os.path.abspath(args.genesis) +
        '    --blocks-dir ' + os.path.abspath(dir) + '/blocks'+
        #'    --config-dir ' + os.path.abspath(dir) +
        '    --data-dir ' + os.path.abspath(dir) +
        '    --chain-state-db-size-mb 1024'+
        '    --verbose-http-errors'+
        ###########################   http-server-address  bnet-endpoint  p2p-listen-endpoint    #########################################

        '    --http-server-address %s:' %(args.http_server) + str(8888) +
        '    --bnet-endpoint %s:' %(args.http_server)+ str(9001) +
        '    --p2p-listen-endpoint %s:'%(args.http_server) + str(9002) +

        ##########################    http-server-address  bnet-endpoint  p2p-listen-endpoint    #########################################

        '    --max-clients ' + str(maxClients) +
        '    --p2p-max-nodes-per-host ' + str(maxClients) +
        '    --enable-stale-production'
        '    --filter-on \"*\"'+
        '    --producer-name ' + account['name'] +
        '    --private-key \'["' + account['pub'] + '","' + account['pvt'] + '"]\''
        '    --accessory datxos::http_accessory'
        '    --accessory datxos::core_api_accessory'
        '    --accessory datxos::producer_accessory' +
        '    --accessory datxos::core_accessory'+
        '    --accessory datxos::core_api_accessory'
        '    --accessory datxos::p2p_net_accessory' +
        '    --accessory datxos::p2p_net_api_accessory' +
        otherOpts)
    with open(dir + 'stderr', mode='w') as f:
        f.write(cmd + '\n\n')
    background(cmd + '    2>>' + dir + 'stderr')
    
# def startProducer(nodeIndex, account):
#     startNode(nodeIndex, account)

def createStakedAccounts(b, e):
#(b,e) =(0,len(accounts))=(0,7),(0~3)users,(4~6)producers
    for i in range(b, e):
        a = accounts[i]

        
        stakeNet = 3000000000000  #3,0000,0000 
        stakeCpu = 3000000000000  #3,0000,0000 
        stakeRam = 2000000000000  

        retry(args.cldatx + 'system newaccount --transfer datxos %s %s --stake-net "%s" --stake-cpu "%s" --buy-ram "%s"   ' % 
            (a['name'], a['pub'], intToCurrency(stakeNet), intToCurrency(stakeCpu), intToCurrency(stakeRam)))
        retry(args.cldatx + 'transfer datxos %s "%s" test' % (a['name'],intToCurrency(2000000000000)))
        sleep(1)
    retry(args.cldatx + 'transfer datxos dotc "%s" test' % (intToCurrency(6000000000000)))
    #>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
    # ramFunds = round(args.ram_funds * 100000000)
    # configuredMinStake = round(args.min_stake * 10000)
    # maxUnstaked = round(args.max_unstaked * 10000)
    # for i in range(b, e):
    #     a = accounts[i]
    #     funds = a['funds']
    #     print('#' * 80)
    #     print('# %d/%d %s %s' % (i, e, a['name'], intToCurrency(funds)))
    #     print('#' * 80)
    #     if funds < ramFunds:
    #         print('skipping %s: not enough funds to cover ram' % a['name'])
    #         continue
    #     minStake = min(funds - ramFunds, configuredMinStake)
    #     unstaked = min(funds - ramFunds - minStake, maxUnstaked)
    #     stake = funds - ramFunds - unstaked
    #     stakeNet = round(stake / 2)
    #     stakeCpu = stake - stakeNet
    #     print('%s: total funds=%s, ram=%s, net=%s, cpu=%s, unstaked=%s' % (a['name'], intToCurrency(a['funds']), intToCurrency(ramFunds), intToCurrency(stakeNet), intToCurrency(stakeCpu), intToCurrency(unstaked)))
    #     assert(funds == ramFunds + stakeNet + stakeCpu + unstaked)
    #     retry(args.cldatx + 'system newaccount --transfer datxos %s %s --stake-net "%s" --stake-cpu "%s" --buy-ram "%s"   ' % 
    #         (a['name'], a['pub'], intToCurrency(stakeNet), intToCurrency(stakeCpu), intToCurrency(ramFunds)))
    #     if unstaked:
    #         retry(args.cldatx + 'transfer datxos %s "%s"' % (a['name'], intToCurrency(unstaked)))
    #>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>    
    
    run(args.cldatx+ 'push action datxos.dtoke transfer \'{"from":"datxos.dbtc","to":"alice","quantity":"300.0000 DBTC","memo":"test"}\' -p datxos.dbtc')
    run(args.cldatx+ 'push action datxos.dtoke transfer \'{"from":"datxos.dbtc","to":"bob","quantity":"300.0000 DBTC","memo":"test"}\' -p datxos.dbtc')
    run(args.cldatx+ 'push action datxos.dtoke transfer \'{"from":"datxos.dbtc","to":"charlie","quantity":"300.0000 DBTC","memo":"test"}\' -p datxos.dbtc')
    for i in range(push_action_times):
        run(args.cldatx+ 'push action datxos.dtoke transfer \'{"from":"alice","to":"datxos.dbtc","quantity":"6.0000 DBTC","memo":"test%d"}\'  -p alice' %(i) )   

#创建系统账户
def createSystemAccounts():
    print(">>>>>>>>>>>>>       wait other nodes ,about %d secs."%(wait_other_node_time))
    print(">>>>>>>>>>>>>       please add other nodes into the net.")
    sleep(wait_other_node_time)
   
    for a in systemAccounts:
        run(args.cldatx + 'create account datxos ' + a + ' ' + args.public_key)

def intToCurrency(i):
    return '%d.%04d %s' % (i // 10000, i % 10000, args.symbol)

def allocateFunds(b, e):
    dist = numpy.random.pareto(1.161, e - b).tolist() # 1.161 = 80/20 rule
    dist.sort()
    dist.reverse()
    factor = 7_600_000_000 / sum(dist)
    total = 0
    for i in range(b, e):
        funds = round(factor * dist[i - b] * 10000)
        if i >= firstProducer and i < firstProducer + numProducers:
            funds = max(funds, round(args.min_producer_funds * 10000))
        total += funds
        accounts[i]['funds'] = funds
    print(total)
    return total



def listProducers():
    run(args.cldatx + 'system listproducers')

def regProducers(b, e):
    for i in range(b, e):
        a = accounts[i]
        retry(args.cldatx + 'system regproducer ' + a['name'] + ' ' + a['pub'] + ' https://' + a['name'] + '.com' + '/' + a['pub'])


def vote(b, e):
    for i in range(b, e):
        voter = accounts[i]['name']
        prods = random.sample(range(firstProducer, firstProducer + numProducers), args.num_producers_vote)
        prods = ' '.join(map(lambda x: accounts[x]['name'], prods))
        retry(args.cldatx + 'system voteproducer prods ' + voter + ' ' + prods)

def stepKillAll():
    run('killall kdatxd noddatx || true')
    sleep(1.5)
def stepStartWallet():
    startWallet()
    importKeys()
def stepStartBoot():
    startNode(0, {'name': 'datxos', 'pvt': args.private_key, 'pub': args.public_key})
    sleep(1.5)
def stepInstallSystemContracts():

    run(args.cldatx + 'set contract datxos.token ' + args.contracts_dir + 'DatxToken/')
    run(args.cldatx + 'set contract datxos.msig ' + args.contracts_dir + 'DatxMsig/')
    run(args.cldatx + 'set contract datxos.charg ' + args.contracts_dir +'DatxRecharge/')

    #run(args.cldatx + 'set contract datxos.dtoke ' + args.contracts_dir + 'DatxToken/')
    #run(args.cldatx + 'set contract datxos.deth ' + args.contracts_dir + 'DatxToken/')
    #run(args.cldatx + 'set contract datxos.deos ' + args.contracts_dir + 'DatxToken/')

    run(args.cldatx + 'set contract datxos.dtoke ' + args.contracts_dir + 'DatxDToken/')
    run(args.cldatx + 'set contract datxos.extra ' + args.contracts_dir + 'DatxExtract/')

    run(args.cldatx + 'set account permission datxos.dtoke active \'{"threshold": 1,"keys": [{"key": "DATX8Znrtgwt8TfpmbVpTKvA2oB8Nqey625CLN8bCN3TEbgx86Dsvr","weight": 1}],"accounts": [{"permission":{"actor":"datxos.charg","permission":"datxos.code"},"weight":1}]}\' owner -p datxos.dtoke')
    
    run(args.cldatx + 'set account permission datxos.dtoke active \'{"threshold": 1,"keys": [{"key": "DATX8Znrtgwt8TfpmbVpTKvA2oB8Nqey625CLN8bCN3TEbgx86Dsvr","weight": 1}],"accounts": [{"permission":{"actor":"datxos.extra","permission":"datxos.code"},"weight":1}]}\' owner -p datxos.dtoke')
    
   

def stepCreateTokens():
    #datx
    run(args.cldatx + 'push action datxos.token create \'["datxos", "10000000000.0000 %s"]\' -p datxos.token' % (args.symbol)+' -x 3500')
    #totalAllocation = allocateFunds(0, len(accounts))   
    totalAllocation=76000000000000
    run(args.cldatx + 'push action datxos.token issue \'["datxos", "%s", "memo"]\' -p datxos' % intToCurrency(totalAllocation))

   


    #btc
    run(args.cldatx + 'push action datxos.dtoke create \'["datxos.dbtc", "21000000.0000 DBTC",0,0,0]\' -p datxos.dtoke' )
    run(args.cldatx + 'push action datxos.dtoke issue \'["datxos.dbtc", "21000000.0000 DBTC", "memo"]\' -p datxos.dbtc')
    
    #run(args.cldatx+ 'push action datxos.dtoke transfer \'{"from":"datxos.dtoke","to":"alice","quantity":"100.0000 DBTC","memo":"test"}\' -p datxos.dtoke')
    
    
    #eth
    run(args.cldatx + 'push action datxos.dtoke create \'["datxos.deth", "102000000.0000 DETH",0,0,0]\' -p datxos.dtoke' )
    run(args.cldatx + 'push action datxos.dtoke issue \'["datxos.deth", "102000000.0000 DETH", "memo"]\' -p datxos.deth')
    #eos
    run(args.cldatx + 'push action datxos.dtoke create \'["datxos.deos", "1000000000.0000 DEOS",0,0,0]\' -p datxos.dtoke' )
    run(args.cldatx + 'push action datxos.dtoke issue \'["datxos.deos", "1000000000.0000 DEOS", "memo"]\' -p datxos.deos')
    
    run(args.cldatx + 'set account permission datxos.dbtc active \'{"threshold": 1,"keys": [{"key": "DATX8Znrtgwt8TfpmbVpTKvA2oB8Nqey625CLN8bCN3TEbgx86Dsvr","weight": 1}],"accounts": [{"permission":{"actor":"datxos.charg","permission":"datxos.code"},"weight":1}]}\' owner -p datxos.dbtc')
    run(args.cldatx + 'set account permission datxos.deth active \'{"threshold": 1,"keys": [{"key": "DATX8Znrtgwt8TfpmbVpTKvA2oB8Nqey625CLN8bCN3TEbgx86Dsvr","weight": 1}],"accounts": [{"permission":{"actor":"datxos.charg","permission":"datxos.code"},"weight":1}]}\' owner -p datxos.deth')
    run(args.cldatx + 'set account permission datxos.deos active \'{"threshold": 1,"keys": [{"key": "DATX8Znrtgwt8TfpmbVpTKvA2oB8Nqey625CLN8bCN3TEbgx86Dsvr","weight": 1}],"accounts": [{"permission":{"actor":"datxos.charg","permission":"datxos.code"},"weight":1}]}\' owner -p datxos.deos')
    

    sleep(1)
def stepSetSystemContract():
    retry(args.cldatx + 'set contract datxos ' + args.contracts_dir + 'DatxSystem/ -x 3500')
    sleep(1)
    run(args.cldatx + 'push action datxos setpriv' + jsonArg(['datxos.msig', 1]) + '-p datxos@active')
def stepCreateStakedAccounts():
    createStakedAccounts(0, len(accounts))

def stepRegProducers():
    regProducers(firstProducer, firstProducer + numProducers)
    sleep(1)
    #listProducers()
def stepVote():
    vote(0, 0 + args.num_voters)
    sleep(1)
    #listProducers()
    sleep(5)

def stepLog():
    run('tail -n 60 ' + args.nodes_dir + '00-datxos/stderr')


parser = argparse.ArgumentParser()

commands = [
    ('k', 'kill',           stepKillAll,                True,    "Kill all noddatx and kdatxd processes"),
    ('w', 'wallet',         stepStartWallet,            True,    "Start kdatxd, create wallet, fill with keys"),
    ('b', 'boot',           stepStartBoot,              True,    "Start boot node"),
    ('s', 'sys',            createSystemAccounts,       True,    "Create system accounts (datxos.*)"),
    ('c', 'contracts',      stepInstallSystemContracts, True,    "Install system contracts (token, msig)"),
    ('t', 'tokens',         stepCreateTokens,           True,    "Create tokens"),
    ('S', 'sys-contract',   stepSetSystemContract,      True,    "Set system contract"),
    ('T', 'stake',          stepCreateStakedAccounts,   True,    "Create staked accounts"),
    ('p', 'reg-prod',       stepRegProducers,           True,    "Register producers"),
    ('v', 'vote',           stepVote,                   True,    "Vote for producers"),
    ('l', 'log',            stepLog,                    True,    "Show tail of node's log"),
]

parser.add_argument('--public-key', metavar='', help="datxOS Public Key", default='DATX8Znrtgwt8TfpmbVpTKvA2oB8Nqey625CLN8bCN3TEbgx86Dsvr', dest="public_key")
parser.add_argument('--private-Key', metavar='', help="datxOS Private Key", default='5K463ynhZoCDDa4RDcr63cUwWLTnKqmdcoTKTHBjqoKfv4u5V7p', dest="private_key")
parser.add_argument('--cldatx', metavar='', help="Cldatx command", default='../../build/programs/cldatx/cldatx --wallet-url http://127.0.0.1:8899 ')
parser.add_argument('--noddatx', metavar='', help="Path to noddatx binary", default='../../build/programs/noddatx/noddatx')
parser.add_argument('--kdatxd', metavar='', help="Path to kdatxd binary", default='../../build/programs/kdatxd/kdatxd')
parser.add_argument('--contracts-dir', metavar='', help="Path to contracts directory", default='../../build/contracts/')
parser.add_argument('--user-limit', metavar='', help="Max number of users. (0 = no limit)", type=int, default=4)
parser.add_argument('--producer-limit', metavar='', help="Maximum number of producers. (0 = no limit)", type=int, default=3)
parser.add_argument('--max-user-keys', metavar='', help="Maximum user keys to import into wallet", type=int, default=10)
parser.add_argument('--num-producers-vote', metavar='', help="Number of producers for which each user votes", type=int, default=3)
parser.add_argument('--ram-funds', metavar='', help="How much funds for each user to spend on ram", type=float, default=0.1)
parser.add_argument('--min-stake', metavar='', help="Minimum stake before allocating unstaked funds", type=float, default=0.9)
parser.add_argument('--max-unstaked', metavar='', help="Maximum unstaked funds", type=float, default=10)
parser.add_argument('--min-producer-funds', metavar='', help="Minimum producer funds", type=float, default=1000.0000)
parser.add_argument('--num-voters', metavar='', help="Number of voters", type=int, default=3)
parser.add_argument('--nodes-dir', metavar='', help="Path to nodes directory", default='./nodes/')
parser.add_argument('--genesis', metavar='', help="Path to genesis.json", default="./genesis.json")
parser.add_argument('--wallet-dir', metavar='', help="Path to wallet directory", default='./wallet/')
parser.add_argument('--log-path', metavar='', help="Path to log file", default='./output.log')
parser.add_argument('--symbol', metavar='', help="The datxos.system symbol", default='DATX')
parser.add_argument('--producer-sync-delay', metavar='', help="Time (s) to sleep to allow producers to sync", type=int, default=100)
parser.add_argument('-a', '--all', action='store_true', help="Do everything marked with (*)")
parser.add_argument('--http-server',default='127.0.0.1', metavar='', help='HTTP server for cldatx')
parser.add_argument('-H', '--http-port',type=int,default=8888, metavar='', help='HTTP port for cldatx')

for (flag, command, function, inAll, help) in commands:
    prefix = ''
    if inAll: prefix += '*'
    if prefix: help = '(' + prefix + ') ' + help
    if flag:
        parser.add_argument('-' + flag, '--' + command, action='store_true', help=help, dest=command)
    else:
        parser.add_argument('--' + command, action='store_true', help=help, dest=command)


args = parser.parse_args()

args.cldatx += '--url http://%s:%d ' % (args.http_server,args.http_port)

logFile = open(args.log_path, 'a')

logFile.write('\n\n' + '*' * 80 + '\n\n\n')

maxClients =  50

# with open('accounts.json') as f:
#     a = json.load(f)
#     if args.user_limit:
#         del a['users'][args.user_limit:]#4
#     if args.producer_limit:
#         del a['producers'][args.producer_limit:]#3
#     # firstProducer = len(a['users'])
#     # numProducers = len(a['producers'])
#     firstvoter=len(a['producers'])
#     numVoters = len(a['users'])
#     firstProducer = 0
#     numProducers = len(a['producers'])
#     accounts =a['producers']+  a['users'] 

with open('accounts.json') as f:
    a = json.load(f)
    if args.user_limit:
        del a['users'][args.user_limit:]#4
    if args.producer_limit:
        del a['producers'][args.producer_limit:]#3
    firstProducer = len(a['users'])
    numProducers = len(a['producers'])
    accounts = a['users'] + a['producers']

haveCommand = False
for (flag, command, function, inAll, help) in commands:
    if getattr(args, command) or inAll and args.all:
        if function:
            haveCommand = True
            function()
if not haveCommand:
    print('bios-boot-tutorial.py: Tell me what to do. -a does almost everything. -h shows options.')


