
var Web3 = require('web3');
var Tx = require('ethereumjs-tx');
var rf = require('fs');
var path = require('path');

var INI = require("../lib/ini-file-loader");
var confPath = path.resolve(__dirname, '../config/config.ini');
var ini___ = INI.loadFileSync(confPath);
var se = ini___.getOrCreateSection("node config");

var eth_endpoint = se["eth-endpoint"];
var eth_ws_provider = se["eth-endpoint-ws"];

var web3 = new Web3(new Web3.providers.HttpProvider(eth_endpoint));
var web3ws = new Web3(new Web3.providers.WebsocketProvider(eth_ws_provider))

async function withdraw(myAddr, contractAddr, to, value, privateKey, data) {
  var nonce = await web3.eth.getTransactionCount(myAddr);
  var privateKey = Buffer.from(privateKey, 'hex');

  var abi = JSON.parse(rf.readFileSync(path.resolve(__dirname, '../config/Wallet.json'),'utf-8')).abi;

  //var contract = new web3.eth.Contract(abi,contractAddr);
  var contractInstance = new web3.eth.Contract(abi,contractAddr);

  var rawTx = {
    nonce: web3.utils.toHex(nonce),
    to: contractAddr,
    gasPrice: web3.utils.toHex(se["eth-gasprice"]),
    gasLimit: web3.utils.toHex(se["eth-gaslimit"]),
    value: '',
    data: contractInstance.methods.execute(to, value, data).encodeABI(),
    chainId: 3
  };
  var tx = new Tx(rawTx);
  tx.sign(privateKey);
  var serializedTx = tx.serialize();
  var rawparam = serializedTx.toString('hex');

  let obj = await web3.eth.sendSignedTransaction('0x' + rawparam);
  return obj;
}

function getContractInstanceWithSocket(){

  var abi = JSON.parse(rf.readFileSync(path.resolve(__dirname, '../config/Wallet.json'),'utf-8')).abi;
  //var contract = web3.eth.contract(abi);
  var contractInstance = new web3ws.eth.Contract(abi,se["eth-muladdress"]);//contract.at(se["eth-muladdress"]);
  return contractInstance;
}

function fromAscii(str) {
  return web3.utils.fromAscii(str);
}

function toAscii(hex){
  return web3.utils.toAscii(hex);
}

module.exports = {
  withdraw,
  fromAscii,
  toAscii,
  getContractInstanceWithSocket
}