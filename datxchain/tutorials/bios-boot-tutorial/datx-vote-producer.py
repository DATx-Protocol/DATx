#!/usr/bin/env python3

import argparse
import json
import os
import random
import re
import subprocess
import sys
import time

args = None
logFile = None

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



# def vote(b, e):
#     for i in range(b, e):
#         voter = accounts[i]['name']
#         prods = random.sample(range(firstProducer, firstProducer + numProducers), args.num_producers_vote)
#         prods = ' '.join(map(lambda x: accounts[x]['name'], prods))
#         retry(args.cldatx + 'system voteproducer prods ' + voter + ' ' + prods)

def vote(name,b,e):
    for i in range(b,e):
        voter=accounts[i]['name']
        retry(args.cldatx + 'system voteproducer prods ' + voter + ' ' + name)

def updateAuth(account, permission, parent, controller):
    run(args.cldatx + 'push action datxos updateauth' + jsonArg({
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

def resign(account, controller):
    updateAuth(account, 'owner', '', controller)
    updateAuth(account, 'active', 'owner', controller)
    sleep(1)
    run(args.cldatx + 'get account ' + account)


def setproducer(b,e):
    list_account=[]
    for i in range(b, e):
        a = accounts[i]
        list_account.append(a)
        #retry(args.cldatx + 'system regproducer ' + a['name'] + ' ' + a['pub'] + ' https://' + a['name'] + '.com' + '/' + a['pub'])  
    retry(args.cldatx + 'push action datxos setprods' + jsonArg({
            "schedule":[
                {
                "producer_name":list_account[0]['name'],
                "block_signing_key": list_account[0]['pub']
            },{
                "producer_name":list_account[1]['name'],
                "block_signing_key": list_account[1]['pub']
            },{
                "producer_name":list_account[2]['name'],
                "block_signing_key": list_account[2]['pub']
            }]
            }) + '-p datxos@active')

def stepVote():
    vote(args.pro_name,0, 0 + args.num_voters)
    sleep(1)

def stepSetProducers():
    
    setproducer(firstProducer, firstProducer + numProducers)

def stepResign():
    resign('datxos', 'datxos.prods')
    for a in systemAccounts:
        resign(a, 'datxos')



parser = argparse.ArgumentParser()

commands = [
    #('e', 'set-prod',       stepSetProducers,           True,    "set producers"),
    ('v', 'vote',           stepVote,                   True,    "Vote for producers"),
    ('q', 'resign',         stepResign,                 True,    "Resign datxos"),
]

parser.add_argument('--public-key', metavar='', help="datxOS Public Key", default='DATX8Znrtgwt8TfpmbVpTKvA2oB8Nqey625CLN8bCN3TEbgx86Dsvr', dest="public_key")
parser.add_argument('--private-Key', metavar='', help="datxOS Private Key", default='5K463ynhZoCDDa4RDcr63cUwWLTnKqmdcoTKTHBjqoKfv4u5V7p', dest="private_key")
parser.add_argument('--cldatx', metavar='', help="Cldatx command", default='cldatx ')
parser.add_argument('--pro-name',default='dotc', metavar='', help='producer name')
parser.add_argument('--user-limit', metavar='', help="Max number of users. (0 = no limit)", type=int, default=4)
parser.add_argument('--producer-limit', metavar='', help="Maximum number of producers. (0 = no limit)", type=int, default=3)
parser.add_argument('--max-user-keys', metavar='', help="Maximum user keys to import into wallet", type=int, default=10)
parser.add_argument('--num-producers-vote', metavar='', help="Number of producers for which each user votes", type=int, default=3)
parser.add_argument('--ram-funds', metavar='', help="How much funds for each user to spend on ram", type=float, default=0.1)
parser.add_argument('--min-stake', metavar='', help="Minimum stake before allocating unstaked funds", type=float, default=0.9)
parser.add_argument('--max-unstaked', metavar='', help="Maximum unstaked funds", type=float, default=10)
parser.add_argument('--min-producer-funds', metavar='', help="Minimum producer funds", type=float, default=1000.0000)
parser.add_argument('--num-voters', metavar='', help="Number of voters", type=int, default=3)
parser.add_argument('--log-path', metavar='', help="Path to log file", default='./output.log')
parser.add_argument('--symbol', metavar='', help="The datxos.system symbol", default='DATX')
parser.add_argument('--producer-sync-delay', metavar='', help="Time (s) to sleep to allow producers to sync", type=int, default=100)
parser.add_argument('-a', '--all', action='store_true', help="Do everything marked with (*)")
# parser.add_argument('--http-server',default='0.0.0.0', metavar='', help='HTTP server for cldatx')

# parser.add_argument('-H', '--http-port',type=int,default=8888, metavar='', help='HTTP port for cldatx')

for (flag, command, function, inAll, help) in commands:
    prefix = ''
    if inAll: prefix += '*'
    if prefix: help = '(' + prefix + ') ' + help
    if flag:
        parser.add_argument('-' + flag, '--' + command, action='store_true', help=help, dest=command)
    else:
        parser.add_argument('--' + command, action='store_true', help=help, dest=command)


args = parser.parse_args()

#args.cldatx += '--url http://%s:%d ' % (args.http_server,args.http_port)

logFile = open(args.log_path, 'a')

logFile.write('\n\n' + '*' * 80 + '\n\n\n')

maxClients =  50

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


