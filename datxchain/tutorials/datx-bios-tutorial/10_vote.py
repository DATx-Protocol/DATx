#!/usr/bin/env python3


#######################################################################
#     
# Scrip Created by DATX Team
# For DATX vote producers.
# 
# Be sure that you have private key and public key of datxos,fix them in 0_bios_config.conf
# Be sure that you are root user if you are bios node.
#
# 1. If:
#     
#    ./10_vote.py
#    
#    At first,The datxos.publickey and datxos.publickey should be repalced by real key. 
#   
#    
# 
# 2. If:
#
#    ./10_vote.py -c your.configfile
#    
#     will read informations in your configfile.
#     
#
#######################################################################

import util
import argparse
import json
import random

def vote(b, e):
    for i in range(b, e):
        voter = accounts[i]['name']
        prods = random.sample(range(firstProducer, firstProducer + numProducers), num_producers_vote-1)
        prods = ' '.join(map(lambda x: accounts[x]['name'], prods))
        util.retry('cldatx ' + 'system voteproducer prods ' + voter + ' ' + prods)


if __name__=="__main__":
    
    parser = argparse.ArgumentParser()
    parser.add_argument('-c','--config', type=str, default='0_bios_config.conf')
    args = parser.parse_args()

    config=util.confArg(args.config)
    symbol=config.get('PARAMETER','symbol')
    contracts_dir=config.get('PATH','contracts-dir')
    user_limit=config.getint('PARAMETER','user-limit')
    num_voters=config.getint('PARAMETER','num-voters')
    producer_limit=config.getint('PARAMETER','producer-limit')
    num_producers_vote=config.getint('PARAMETER','num-producers-vote')


    with open('accounts.json') as f:
        a = json.load(f)
        if user_limit:
            del a['users'][user_limit:]#10
        if producer_limit:
            del a['producers'][producer_limit:]#3
        firstProducer = len(a['users'])
        numProducers = len(a['producers'])
        accounts = a['users'] + a['producers']

    vote(0, 0 + num_voters)