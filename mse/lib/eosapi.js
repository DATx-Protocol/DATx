const Eos = require('eosjs');

const INI = require("../lib/ini-file-loader");
const se = INI.getConfigFile();

const BigNumber = require("bignumber.js")

const chainId = se["eos-chainid"];
const httpEndpoint = se["eos-endpoint"];
const eos = Eos({httpEndpoint, chainId, keyProvider: () => userProvidedKey,expireInSeconds: 600});

let userProvidedKey = null;

async function propose(proposer, MultiSigAccount, to, value, memo, auths, ProvidedKey) {
  userProvidedKey = ProvidedKey;
  transfer = await eos.transfer(MultiSigAccount, to, value, memo, {broadcast: false, sign: false});
  transfer.transaction.transaction.max_net_usage_words = 0;

  msig = await eos.contract('eosio.msig');
  randomName = String(Math.round(Math.random() * 1e12)).replace(/[0,6-9]/g, '');
  result = await msig.propose(proposer, randomName, auths, transfer.transaction.transaction);

  return randomName;
}

async function confirm(propser, proposeName, actor, ProvidedKey) {
  userProvidedKey = ProvidedKey;
  msig = await eos.contract('eosio.msig');
  confirm = await msig.approve(
      propser, proposeName, {actor: actor, permission: 'active'},
      {authorization: actor + '@active'});
  return confirm;
}

async function exec(proposer, proposeName, executer, ProvidedKey) {
  userProvidedKey = ProvidedKey;
  msig = await eos.contract('eosio.msig');
  result = await msig.exec(proposer, proposeName, executer, {authorization: executer + '@active'});
  return result;
}

async function getProposeAction(proposer,proposeName){
  try{
    encodedName = new BigNumber(Eos.modules.format.encodeName(proposeName, false));
    result = await eos.getTableRows({
        code: 'eosio.msig',
        json: true,
        limit: 1,
        lower_bound: encodedName.toString(),
        scope: proposer,
        table: 'proposal'
      });
    trxData = result.rows[0].packed_transaction;
    trx = eos.fc.fromBuffer('transaction', trxData);
    if(trx.actions.length != 1 || trx.actions[0].account != 'eosio.token' || trx.actions[0].name != 'transfer') throw new Error('invalid trx');
    actData = trx.actions[0].data;
    token = await eos.contract('eosio.token');
    act = token.fc.fromBuffer('transfer',actData);
    return {result:true,data:JSON.stringify(act)};    
  }
  catch(e){
    return {result:false,data:e.toString()};
  }
}

function EccSIgn(data,privateKey){
  return Eos.modules.ecc.sign(data,privateKey);
}

function EccVerify(sign,data,publicKey){
  return Eos.modules.ecc.verify(sign,data,publicKey);
}


module.exports = {
  propose,
  confirm,
  exec,
  getProposeAction,
  EccSIgn,
  EccVerify
}
