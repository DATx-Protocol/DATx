#!/usr/bin/env python3


#######################################################################
#     
# Scrip Created by DATX Team
# For DATX bios node register some new producers.
# 
# Be sure that you have private key and public key ,fix them in 0_producer_config.conf
# 
#
# 1. To do:
#     
#    ./9_Reg_Producers.py
#     
#    This will create register some new producers.
#    
#    
#
#######################################################################

import util
import argparse
import json

def RegProducers():
    
    util.retry('cldatx system regproducer ' + producer_name + ' ' + public_key + ' '+ producer_url + ' ' + verifier_url) 


if __name__=="__main__":
    
    parser = argparse.ArgumentParser()
    parser.add_argument('-c','--config', type=str, default='0_producer_config.conf')
    args = parser.parse_args()

    config=util.confArg(args.config)

    verifier_url=config.get('IP','verifier-url')
    producer_url=config.get('IP','producer-url')

    producer_name=config.get('NAME','producer-name')
    public_key=config.get('KEY','public-key')


    RegProducers()