// bitcoinapi = require('../lib/bitcoinapi');

// (async function(){
//     script = '522103a4ac53ded034de0ce8e1a5aa8cae967a7c33f8ef807ee31d0a972fbcd912c8cb21038e6c355aa3a7b0a3338215e1fb952c1c255eab07012c800a151f8fd7bb9feac92102c8a936b526d91e6047569ec8fd53779a2368a150d63cea655fc9c7ba66d2199e53ae'
//     //buildTrx
//     trx = await bitcoinapi.buildTrx('2MsimupueVskjJMy79kKGP5uzfCWfZuK8TD','n4fc3bKTVrVRveBrHZf5Zv4wGGBTf5sdHg',2e5,1e5,{IsTestnet:1})

//     //sign locally
//     trx = bitcoinapi.signTrx(trx,'cV6dJVkPZfmmDg5XHNpExd8fgMTaoYgv7NusH26ZEFj2vcA9xvhz',script,{IsTestnet:1})

//     //sign remotly
//     trxSerialize = trx.toHex()
//     trx = bitcoinapi.getTrxFromHex(trxSerialize)
//     trx = bitcoinapi.signTrx(trx,'cQX7uYfhPJJGPdoqgHQXLC4WQWcDXGH7DgLs4m89xWykYkg1V4iA',script,{IsTestnet:1})
//     trxSerialize = trx.toHex()
//     result = bitcoinapi.decodeOutput(trx,{IsTestnet:1})
//     //broadcast trx
//     trx = bitcoinapi.getTrxFromHex(trxSerialize)
//     //bitcoinapi.broadcastTrx(trx) 
// })()


const BigNumber = require("bignumber.js")
Eos = require('eosjs');
var path = require('path');

var INI = require("../lib/ini-file-loader");
var confPath = path.resolve(__dirname, '../config/config.ini');
var ini___ = INI.loadFileSync(confPath);
var se = ini___.getOrCreateSection("node config");

let userProvidedKey = '5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3';

chainid = se["eos-chainid"];
httpEndpoint = se["eos-endpoint"];
eos = Eos({httpEndpoint, chainid, keyProvider: () => userProvidedKey});
(async function(){
    // msig = await eos.contract('eosio.msig');
    // result = await msig.review('alice','22154');


// const accountName = '22154';
// const encodedName = new BigNumber(Eos.modules.format.encodeName(accountName, false));
// result = await eos.getTableRows({
//     code: 'eosio.msig',
//     json: true,
//     limit: 1,
//     lower_bound: encodedName.toString(),
//     scope: 'alice',
//     table: 'proposal'
//   });
//   console.log(result);
  //d43d865b4217a1efc88b000000000100a6823403ea3055000000572d3ccdcd0100000000000f917900000000a8ed32322100000000000f91790000000000000e3d20a107000000000004454f53000000000000
  

    // transfer = await eos.transfer('jacky', 'bob', '50.0000 EOS', '', {broadcast: false, sign: false});
    // console.log(transfer);
    //data = 'd43d865b4217a1efc88b000000000100a6823403ea3055000000572d3ccdcd0100000000000f917900000000a8ed32322100000000000f91790000000000000e3d20a107000000000004454f53000000000000';
    //result = eos.fc.fromBuffer('transaction', data);
    // type = {type: 1, data: '00ff'};
    // buffer = eos.fc.toBuffer('extensions_type', type);
    // result = eos.fc.fromBuffer('extensions_type', buffer);
    //console.log(result);
    // token = await eos.contract('eosio.token');
    // result = token.fc.fromBuffer('transfer','00000000000f91790000000000000e3d20a107000000000004454f530000000000');
    // //result = eos.fc.fromBuffer('transfer','00000000000f91790000000000000e3d20a107000000000004454f530000000000');
    // console.log(result);

})();


var Web3 = require('web3');
var Tx = require('ethereumjs-tx');
var rf = require('fs');
var path = require('path');

var INI = require("../lib/ini-file-loader");
var confPath = path.resolve(__dirname, '../config/config.ini');
var ini___ = INI.loadFileSync(confPath);
var se = ini___.getOrCreateSection("node config");

var eth_endpoint = se["eth-endpoint"];

var web3 = new Web3(new Web3.providers.HttpProvider(eth_endpoint));
result = web3.utils.toAscii('0x6173646164616461');
console.log(result);

