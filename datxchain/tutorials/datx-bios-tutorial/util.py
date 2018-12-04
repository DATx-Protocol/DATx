#!/usr/bin/env python3
import argparse
import json
import os
import random
import re
import subprocess
import sys
import time
import configparser

# args = None
# logFile = None
log_path='./output.log'

def confArg(config_file):
    cf = configparser.ConfigParser()
    cf.read(config_file)
    return cf

def conf_logpath(config_file):
    cf = configparser.ConfigParser()
    cf.read(config_file)
    log_path=cf.get('PATH','log-path')
    return log_path



def jsonArg(a):
    return " '" + json.dumps(a) + "' "

def run(args):
    print(args)
    with open(log_path,'a') as logFile:
        logFile.write(args + '\n')   
    if subprocess.call(args, shell=True):
        print('exiting because of error')
        sys.exit(1)

def retry(args):
    while True:
        print(args)
        with open(log_path,'a') as logFile:
            logFile.write(args + '\n') 
        if subprocess.call(args, shell=True):
            print('*** Retry')
        else:
            break

def background(args):
    print(args)
    with open(log_path,'a') as logFile:
        logFile.write(args + '\n') 
    return subprocess.Popen(args, shell=True)

def getOutput(args):
    print(args)
    with open(log_path,'a') as logFile:
        logFile.write(args + '\n') 
    proc = subprocess.Popen(args, shell=True, stdout=subprocess.PIPE)
    return proc.communicate()[0].decode('utf-8')

def getJsonOutput(args):
    print(args)
    with open(log_path,'a') as logFile:
        logFile.write(args + '\n') 
    proc = subprocess.Popen(args, shell=True, stdout=subprocess.PIPE)
    return json.loads(proc.communicate()[0])

def sleep(t):
    print('sleep', t, '...')
    time.sleep(t)
    print('resume')


