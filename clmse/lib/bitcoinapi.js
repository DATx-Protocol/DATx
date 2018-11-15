const bitcoin = require('bitcoinjs-lib');

function genKeyPairs(options) {
    try{
      IsTestnet = options && options.IsTestnet;
      network = IsTestnet ? bitcoin.networks.testnet : bitcoin.networks.bitcoin;
      keyPair = bitcoin.ECPair.makeRandom({network: network});
      return {
        wif: keyPair.toWIF(),
        pubkey: keyPair.publicKey.toString('hex'),
        prikey: keyPair.privateKey.toString('hex')
      };
    }catch(e){
      return {Error:e.toString()};
    }
  }
  
  function getKeysFromWIF(wif, options) {
    try{
      IsTestnet = options && options.IsTestnet;
      network = IsTestnet ? bitcoin.networks.testnet : bitcoin.networks.bitcoin;
      keyPair = bitcoin.ECPair.fromWIF(wif, network);
      return {
        pubkey: keyPair.publicKey.toString('hex'),
        prikey: keyPair.privateKey.toString('hex')
      };
    }catch(e){
      return {Error:e.toString()};
    }
  }
  
  function genP2PKHAddr(publicKey, options) {
    try{
      publicKey = Buffer.from(publicKey, 'hex');
      IsTestnet = options && options.IsTestnet;
      network = IsTestnet ? bitcoin.networks.testnet : bitcoin.networks.bitcoin;
      address = bitcoin.payments.p2pkh({pubkey: publicKey, network: network}).address;
      return address;
    }catch(e){
      return e.toString();
    }
  }
  
  function genMulSigAddr(pubkeys, num, options) {
    try{
      pubkeys = pubkeys.map(pubkey => Buffer.from(pubkey, 'hex'));
      IsTestnet = options && options.IsTestnet;
      network = IsTestnet ? bitcoin.networks.testnet : bitcoin.networks.bitcoin;
      p2ms = bitcoin.payments.p2ms({m: num, pubkeys: pubkeys, network: network});
      p2sh = bitcoin.payments.p2sh({redeem: p2ms, network: network});
      return {address: p2sh.address, script: p2sh.redeem.output.toString('hex')};
    }catch(e){
      return {Error:e.toString()};
    }
  }


  module.exports = {
    genKeyPairs,
    getKeysFromWIF,
    genP2PKHAddr,
    genMulSigAddr,
  }