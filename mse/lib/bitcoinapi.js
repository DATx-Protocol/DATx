
const bitcoin = require('bitcoinjs-lib');
const http2 = require('./http');
const urlParse = require('url').parse;
const pushtx = require('blockchain.info/pushtx').usingNetwork(3);

const redis = require('./redis.js');

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
    return {};
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
    return {};
  }
}

function genP2PKHAddr(publicKey, options) {
  try{
    publicKey = Buffer.from(publicKey, 'hex');
    IsTestnet = options && options.IsTestnet;
    network = IsTestnet ? bitcoin.networtexks.testnet : bitcoin.networks.bitcoin;
    address = bitcoin.payments.p2pkh({pubkey: publicKey, network: network}).address;
    return address;
  }catch(e){
    return '';
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
    return {};
  }
}

function getTrxDetail(trxid,options){
  IsTestnet = options && options.IsTestnet;
  network = IsTestnet ? bitcoin.networks.testnet : bitcoin.networks.bitcoin;
  var url = (IsTestnet ? 'https://testnet.blockchain.info/rawtx/' :
                         'https://blockchain.info/rawtx/') + trxid;
  var netOptions = urlParse(url);

  let data = '';
  return new Promise(function(resolve, reject) {
    let req = http2.get(netOptions, function(res) {
      res.setEncoding('utf8');
      res.on('data', function(chunk) {
        data += chunk;
      });

      res.on('end', function() {
        resolve({result: true, data: data});
      });
    });

    req.on('error', (e) => {
      resolve({result: false, errmsg: e.message});
    });
    setTimeout(_=>{
      resolve({result: false, errmsg: 'Promise time out'});
     }, 5000);
    req.end();
  });
}

function getUTXOS(from, options) {
  IsTestnet = options && options.IsTestnet;
  network = IsTestnet ? bitcoin.networks.testnet : bitcoin.networks.bitcoin;
  var url = (IsTestnet ? 'https://testnet.blockchain.info/unspent?active=' :
                         'https://blockchain.info/unspent?active=') +
      from;
  var netOptions = urlParse(url);

  let data = '';
  return new Promise(function(resolve, reject) {
    let req = http2.get(netOptions, function(res) {
      res.setEncoding('utf8');
      res.on('data', function(chunk) {
        data += chunk;
      });

      res.on('end', function() {
        resolve({result: true, data: data});
      });
    });

    req.on('error', (e) => {
      resolve({result: false, errmsg: e.message});
    });
    setTimeout(_=>{
      resolve({result: false, errmsg: 'Promise time out'});
     }, 5000);
    req.end();
  });
}

async function buildTrx(from, to, value, fee, options,memo) {
  IsTestnet = options && options.IsTestnet;
  network = IsTestnet ? bitcoin.networks.testnet : bitcoin.networks.bitcoin;
  chunk = await getUTXOS(from, options);

  utxos = JSON.parse(chunk.data);
  utxosUsed = new Array();
  sum = 0;
  for (let unspent of utxos.unspent_outputs) {
    if (await redis.sismemberAsync("spentUTXOS",unspent.tx_hash_big_endian + unspent.tx_output_n))
      continue;

    utxosUsed.push(unspent);
    redis.client.sadd("spentUTXOS",unspent.tx_hash_big_endian + unspent.tx_output_n);
    sum += unspent.value;
    if (sum >= value) break;
  }

  const txb = new bitcoin.TransactionBuilder(network);
  for (let utxo of utxosUsed) {
    txb.addInput(utxo.tx_hash_big_endian, utxo.tx_output_n);
  }

  txb.addOutput(to, value - fee);
  txb.addOutput(from, sum - value);

  data = Buffer.from(memo, 'utf8');
  embed = bitcoin.payments.embed({ data : [data],network : network });
  txb.addOutput(embed.output, 0);
  return txb.buildIncomplete();

}

async function signTrx(trx, wif, script, options) {
  insDecode = decodeInput(trx);
  IsTestnet = options && options.IsTestnet;
  network = IsTestnet ? bitcoin.networks.testnet : bitcoin.networks.bitcoin;
  keyPair = bitcoin.ECPair.fromWIF(wif, network);
  scriptbuffer = Buffer.from(script, 'hex');
  inputCount = trx.ins.length;
  txb = bitcoin.TransactionBuilder.fromTransaction(trx, network);
  for (var i = 0; i < inputCount; i++) {
    try{
      txb.sign(i, keyPair, scriptbuffer);
      if (! await redis.sismemberAsync("spentUTXOS",insDecode[i].txid + insDecode[i].n.toString()))
      redis.client.sadd("spentUTXOS",insDecode[i].txid + insDecode[i].n.toString());
    }
    catch(e){
      // do nothing
    }
  }

  return txb.buildIncomplete();
}

async function broadcastTrx(trx, options) {
  IsTestnet = options && options.IsTestnet;
  network = IsTestnet ? bitcoin.networks.testnet : bitcoin.networks.bitcoin;
  txb = bitcoin.TransactionBuilder.fromTransaction(trx, network);
  trx = txb.build();
  result = pushtx.pushtx(trx.toHex());
  return trx.getId();
}

function getTrxFromHex(trxHex) {
  trxbuffer = Buffer.from(trxHex, 'hex');
  trx = bitcoin.Transaction.fromBuffer(trxbuffer);
  return trx;
}

function getTrxBuilder(trx,options){
  IsTestnet = options && options.IsTestnet;
  network = IsTestnet ? bitcoin.networks.testnet : bitcoin.networks.bitcoin;
  txb = bitcoin.TransactionBuilder.fromTransaction(trx, network);
  return txb;
}

function decodeInput(tx){
  var result = [];
  tx.ins.forEach(function(input, n){
      var vin = {
          txid: input.hash.reverse().toString('hex'),
          n : input.index,
          script: bitcoin.script.toASM(input.script),
          sequence: input.sequence,
      };
      input.hash.reverse();
      result.push(vin);
  })
  return result
}

function decodeOutput(tx, options){
  IsTestnet = options && options.IsTestnet;
  network = IsTestnet ? bitcoin.networks.testnet : bitcoin.networks.bitcoin;

  var format = function(out, n, network){
      var vout = {
          satoshi: out.value,
          value: (1e-8 * out.value).toFixed(8),
          n: n,
          scriptPubKey: {
              asm: bitcoin.script.toASM(out.script),
              hex: out.script.toString('hex'),
              type: bitcoin.classify.output(out.script),
              addresses: [],
          },
      };
      switch(vout.scriptPubKey.type){
      case 'pubkeyhash':
      case 'scripthash':
          vout.scriptPubKey.addresses.push(bitcoin.address.fromOutputScript(out.script, network));
          break;
      case 'nulldata':
          vout.scriptPubKey.addresses.push(bitcoin.script.decompile(out.script).slice(1).toString('utf8'));
          break;
      }
      return vout
  }

  var result = [];
  tx.outs.forEach(function(out, n){
      result.push(format(out, n, network));
  })
  return result
}

function decodeMemo(script){
  result = bitcoin.script.decompile(Buffer.from(script,'hex')).slice(1).toString('utf8');
  return result;
}

module.exports = {
  genKeyPairs,
  getKeysFromWIF,
  genP2PKHAddr,
  genMulSigAddr,
  buildTrx,
  signTrx,
  broadcastTrx,
  getTrxFromHex,
  getTrxBuilder,
  decodeInput,
  decodeOutput,
  getTrxDetail,
  decodeMemo
}
