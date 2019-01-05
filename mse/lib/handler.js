const url = require('url');
const querystring = require('querystring');

const bitcoinapi = require('./bitcoinapi.js');
const ethapi = require('./ethapi.js');
const eosapi = require('./eosapi.js');

const INI = require("./ini-file-loader.js");
const se = INI.getConfigFile();
const httpClient = require('./client.js');
const redis = require('./redis.js');
const util = require("./util")

// The callback to handle requests
function handler(request, response) {
  try {
    var pathName = url.parse(request.url).pathname;
    if (pathName == '/btc/withdraw'){
      util.log('get btc withdraw,URl: ' + request.url);
      var IsTestnet = 0;
      var trxid = '';
      var fee = 0;
      var nodeName = '';
      var sign = '';
      var query = url.parse(request.url).query;
      if(query){
        var params = querystring.parse(query)
        if(params.isTestnet == 1){
          IsTestnet = 1;
        }
        trxid = params.trxid;
        fee = params.fee;
        nodeName = params.nodeName;
        sign = params.sign;
        if(!trxid){
          throw new Error('invalid params,url:' + request.url);
        }
        else{
          (async function(){
          try{ 
            let {to,value,symbol} = await checkDATXWithdraw(trxid);
            var checkData = trxid;
            var checkResult = await checkRequestSIgn(nodeName,sign,checkData);
            if(!checkResult){
              throw new Error('invalid sign ,trxid:' + trxid);
            }
            checkResult = symbol == 'DBTC';
            if(!checkResult){
              throw new Error('transcation rejected,trxid:' + trxid);
            }
            //check extract right
            checkResult = await checkExtractRight(nodeName,trxid);
            if(!checkResult){
              throw new Error('Extract rejected,trxid:' + trxid);
            }

            //buildTrx
            var trx = await bitcoinapi.buildTrx(se["btc-muladdress"],to,parseInt(value * 1e8),fee,{IsTestnet:IsTestnet},trxid); 
            util.log('build trx for ' + trxid + '********:' + trx.toHex());
            //call other nodes
            util.runDatxCmd("cldatx push action datxos.extra setdoing '[\"" + trxid + "\",\"" + nodeName + "\",\"" + se["producer-name"] + "\"]' -p " + se["producer-name"]);
            var trxSerialize = await btcGatherSign(trx.toHex(),trxid,nodeName,sign,IsTestnet);

            //broadcast
            trx = bitcoinapi.getTrxFromHex(trxSerialize);
            result = await bitcoinapi.broadcastTrx(trx); 
            redis.client.set(trxid,trxSerialize);
            util.log('broadcast success,datx trxid:' + trxid);

            rep(response,200,result);
          }catch(e){
            util.log('Error:' + e.toString());

            rep(response,500,'Error:' + e.toString());
          }
          })()
        }
      }
    }

    else if(pathName == '/btc/signTrx'){
      util.log('get btc signTrx,URl: ' + request.url);
      var IsTestnet = 0;
      var trxSerialize = '';
      var trxid = '';
      var nodeName = '';
      var sign = '';
      var script = se["btc-mulscript"];
      var query = url.parse(request.url).query;
      if(query){
        var params = querystring.parse(query)
        if(params.isTestnet == 1){
          IsTestnet = 1;
        }
        trxSerialize = params.trxSerialize;
        trxid = params.trxid;
        nodeName = params.nodeName;
        sign = params.sign;
        if(!trxSerialize || !trxid){
          throw new Error('invalid params,url:' + request.url);
        }
        else{
          (async function(){
          try{
            var checkData = trxid; 
            var checkResult = await checkRequestSIgn(nodeName,sign,checkData);
            if(!checkResult){
              throw new Error('invalid sign,trxid:' + trxid);
            }

            //check extract right
            checkResult = await checkExtractRight(nodeName,trxid);
            if(!checkResult){
              throw new Error('Extract rejected,trxid:' + trxid);
            }

            var recordTrxSerialize = await redis.getAsync(trxid);
            var trx = bitcoinapi.getTrxFromHex(trxSerialize);
            var ins = bitcoinapi.decodeInput(trx);
            var outs = bitcoinapi.decodeOutput(trx,{IsTestnet:IsTestnet});
            var isFirst = true;
            //do some check
            if(recordTrxSerialize == null || recordTrxSerialize == undefined || recordTrxSerialize == ""){
              var sum = 0;
              var input;
              var refTrx;
              for (var i =0; i < ins.length; i++){
                input = ins[i];
                refTrx = await bitcoinapi.getTrxDetail(input.txid,{IsTestnet:IsTestnet});
                if(refTrx.result == false) throw new Error('transcation rejected,can not get trx detail');
                refTrx = JSON.parse(refTrx.data);
                if(refTrx.out[input.n].addr != se["btc-muladdress"]) throw new Error('transcation rejected,invalid inputs multi sig address');
                sum += Number(refTrx.out[input.n].value);
              }

              if(outs.length != 3) throw new Error('transcation rejected invalid outs length');
              if(outs[1].scriptPubKey.addresses[0] != se["btc-muladdress"])  throw new Error('transcation rejected invalid outs multi sig address');
              if(bitcoinapi.decodeMemo(outs[2].scriptPubKey.hex) != trxid) throw new Error('transcation rejected,invalid memo');
              
              var to1 = outs[0].scriptPubKey.addresses[0];
              var fee = sum - outs[0].satoshi - outs[1].satoshi;
              var value1 = outs[0].satoshi;
              if(fee > Number(se["btc-maxfee"])) throw new Error('transcation rejected,invalid fee');

              let {to,value,symbol} = await checkDATXWithdraw(trxid);
              checkResult = (to1 == to && value == parseFloat((value1+fee)/1e8).toFixed(4) && symbol == 'DBTC')
              if(!checkResult)  throw new Error('transcation rejected,check trx detail fail');
            }
            else{
              isFirst = false;
              var recordTrx = bitcoinapi.getTrxFromHex(recordTrxSerialize);
              var recordIns = bitcoinapi.decodeInput(recordTrx);
              var recordOuts = bitcoinapi.decodeOutput(recordTrx,{IsTestnet:IsTestnet});
              if(recordIns.length != ins.length || recordOuts.length != outs.length) throw new Error('transcation rejected');

              for(var i =0; i < ins.length; i++){
                if(ins[i].txid != recordIns[i].txid ||
                   ins[i].n != recordIns[i].n ||
                   ins[i].sequence != recordIns[i].sequence)
                   throw new Error('transcation rejected,ins not the same');
              }

              for(var i =0; i < outs.length; i++){
                if(outs[i].scriptPubKey.hex != recordOuts[i].scriptPubKey.hex ||
                   outs[i].scriptPubKey.satoshi != recordOuts[i].scriptPubKey.satoshi)
                   throw new Error('transcation rejected,outs not the same');
              }
            }
            trx = await bitcoinapi.signTrx(trx,se["btc-wif"],script,{IsTestnet:IsTestnet});
            trxSerialize = trx.toHex();
            redis.client.set(trxid,trxSerialize);

            rep(response,200,trxSerialize);

            //call other nodes
            if(isFirst){
              util.runDatxCmd("cldatx push action datxos.extra setdoing '[\"" + trxid + "\",\"" + nodeName + "\",\"" + se["producer-name"] + "\"]' -p " + se["producer-name"]);

              trxSerialize = await btcGatherSign(trxSerialize,trxid,nodeName,sign,IsTestnet);
              trx = bitcoinapi.getTrxFromHex(trxSerialize);
              result = await bitcoinapi.broadcastTrx(trx); 
            }
          }catch(e){
            util.log(e.toString());
            rep(response,500,'Error:' + e.toString());
          }
          })();
        }
      }
    }

    else if(pathName == '/btc/decodeMemo'){
      var query = url.parse(request.url).query;
      if(query){
        var params = querystring.parse(query)
        var script = params.script;
        var result = bitcoinapi.decodeMemo(script);

        rep(response,200,result);
      }
    }

    //eth withdraw
    else if(pathName == '/eth/withdraw'){
      util.log('get eth withdraw,URl: ' + request.url);
      var myAddr = se["eth-myaddress"];
      var contractAddr = se["eth-muladdress"];
      var myPrivatekey = se["eth-privatekey"];
      var trxid = '';
      var sign = '';
      var nodeName  = '';
      var isInform = false;
      var query = url.parse(request.url).query;
      if(query){
        var params = querystring.parse(query)
        trxid = params.trxid;
        sign = params.sign;
        nodeName = params.nodeName;
        isInform = params.isInform;
        if(!trxid){
          throw new Error('invalid params');
        }
        else{
          (async function(){
          try{
            //do some check
            let {to,value,symbol} = await checkDATXWithdraw(trxid);
            var checkData = trxid;
            var checkResult = await checkRequestSIgn(nodeName,sign,checkData);
            if(!checkResult){
              throw new Error('invalid sign,trxid:' + trxid);
            }
            if(await redis.getAsync(trxid) == "done"){
              throw new Error('handled trx id,trxid:' + trxid);
            }

            checkResult = symbol == 'DETH';
            if(!checkResult)  throw new Error('transcation rejected,trxid:' + trxid);

            //check extract right
            checkResult = await checkExtractRight(nodeName,trxid);
            if(!checkResult){
              throw new Error('Extract rejected,trxid:' + trxid);
            }

            //inform others to configm
            if(!isInform){
              var verifyNodes = await getVerifiers();
              for(var i = 0;i < verifyNodes.length;i++){
                if(verifyNodes[i].account == se["producer-name"]){
                  util.runDatxCmd("cldatx push action datxos.extra setdoing '[\"" + trxid + "\",\"" + nodeName + "\",\"" + se["producer-name"] + "\"]' -p " + se["producer-name"]);
                  ethapi.withdraw(myAddr,contractAddr,to,parseInt(value*1e18),myPrivatekey,ethapi.fromAscii(trxid));
                  continue;
                }

                var URL = verifyNodes[i].url + '/eth/withdraw?trxid=' + trxid + '&to=' + to + '&value=' + value + '&isInform=true&sign=' + sign + '&nodeName=' + nodeName;
                httpClient.requestAsync(URL,'GET',null,function (error, response, body) {
                  if(error) util.log(error.toString());
                });
              }
            }
            else{
              util.runDatxCmd("cldatx push action datxos.extra setdoing '[\"" + trxid + "\",\"" + nodeName + "\",\"" + se["producer-name"] + "\"]' -p " + se["producer-name"]);
              ethapi.withdraw(myAddr,contractAddr,to,parseInt(value*1e18),myPrivatekey,ethapi.fromAscii(trxid));
            }

            redis.client.set(trxid,"done");
            rep(response,200,'success');
          }catch(e){
            rep(response,500,'Error:' + e.toString());
          }
          })();
        }
      }
    }

    else if(pathName == '/eos/getAccount'){
      rep(response,200,se["eos-account"]);
    }

    //propose
    else if(pathName == '/eos/withdraw'){
      util.log('get eos withdraw,URl: ' + request.url);
      var proposer = se["eos-account"]; 
      var ProvidedKey = se["eos-privateKey"]
      var auths = []; 
      var MultiSigAccount = se["eos-mulAccount"]; 
      var trxid = '';
      var sign = '';
      var nodeName = '';
      var query = url.parse(request.url).query;
      if(query){
        var params = querystring.parse(query);
        trxid = params.trxid;
        sign = params.sign;
        nodeName = params.nodeName;
        if(!trxid){
          throw new Error('invalid params');
        }
        else{
          (async function(){
          try{
            //check 
            let {to,value,symbol} = await checkDATXWithdraw(trxid);
            var checkData = trxid;
            var checkResult = await checkRequestSIgn(nodeName,sign,checkData);
            if(!checkResult){
              throw new Error('invalid sign,trxid:' + trxid);
            }
            checkResult = symbol == 'DEOS';
            if(!checkResult)  throw new Error('transcation rejected,trxid:' + trxid);

            //check extract right
            checkResult = await checkExtractRight(nodeName,trxid);
            if(!checkResult){
              throw new Error('Extract rejected,trxid:' + trxid);
            }

            if(await redis.getAsync(trxid) == "done"){
              throw new Error('handled trx id,trxid:' + trxid);
            }

            var verifyNodes = await getVerifiers();
            for (var i = 0; i < verifyNodes.length; i++){
              if(verifyNodes[i].account == se["producer-name"]){
                auths.push({actor: se["eos-account"], permission: 'active'});
                continue;
              }
              URL = verifyNodes[i].url + '/eos/getAccount';
              var result = await httpClient.requestWithOvertime(URL,3000,'GET');
              
              if(result.result){
                auths.push({actor: result.data, permission: 'active'});
              }
            }
            util.log("proposer auths count is" + auths.length);
            util.runDatxCmd("cldatx push action datxos.extra setdoing '[\"" + trxid + "\",\"" + nodeName + "\",\"" + se["producer-name"] + "\"]' -p " + se["producer-name"]);
            var proposeName = await eosapi.propose(proposer,MultiSigAccount,to,value + ' EOS',trxid,auths,ProvidedKey);
            util.log('proposeName for ' + trxid + ' is ' + proposeName);

            //just for now,should use a delay queue
            util.log("wait 10 sec,make sure this propose broadcast well")
            await sleep(10000);
            //broadcast to other nodes
            for (var i = 0; i < verifyNodes.length; i++){
              if(verifyNodes[i].account == se["producer-name"]){
                eosapi.confirm(proposer,proposeName,proposer,ProvidedKey);
                continue;
              }

              var URL = verifyNodes[i].url + '/eos/confirm?proposer=' + proposer + '&proposeName=' + proposeName + '&trxid=' + trxid + '&nodeName=' + nodeName + '&sign=' + sign;;
              util.log('request eos confirm,url:' + URL);
              httpClient.requestAsync(URL,'GET',null,function (error, response, body) {
                if (error) util.log(error.toString());
              });
            }

            //wait 15 second then exec
            setTimeout(function(proposer,proposeName,proposer,ProvidedKey,response){
              eosapi.exec(proposer,proposeName,proposer,ProvidedKey).then(result =>{ 
                if(result && result.transaction_id){
                  redis.client.set(trxid,"done");
                  rep(response,200,result.transaction_id);
                }
                else{
                  util.log("eos withdraw err happened:" + result.toString());
                  rep(response,200,result);
                }
              })       
            },15000,proposer,proposeName,proposer,ProvidedKey,response);
          }catch(e){
            util.log('Error:' + e.toString());
            rep(response,500,'Error:' + e.toString());
          }
          })();
        }
      }
    }

    //confirm
    else if(pathName == '/eos/confirm'){
      util.log('get eos confirm,URL:' + request.url);
      var proposer = '';
      var proposeName = '';
      var trxid = '';
      var sign = '';
      var nodeName = '';
      var query = url.parse(request.url).query;
      if(query){
        var params = querystring.parse(query);
        proposer = params.proposer;
        proposeName = params.proposeName;
        trxid = params.trxid;
        nodeName = params.nodeName;
        sign = params.sign;
        if(!proposer || !proposeName || !trxid){
          throw new Error('invalid params,request url:'+ request.url);
        }
        else{
          (async function(){
          try{
            //do some check
            checkData = trxid; 
            checkResult = await checkRequestSIgn(nodeName,sign,checkData);
            if(!checkResult){
              throw new Error('invalid sign,trxid:' + trxid);
            }
            if(await redis.getAsync(trxid) == "done"){
              throw new Error('handled trx id,trxid:' + trxid);
            }

            //check extract right
            checkResult = await checkExtractRight(nodeName,trxid);
            if(!checkResult){
              throw new Error('Extract rejected,trxid:' + trxid);
            }

            var check = await eosapi.getProposeAction(proposer,proposeName);
            if(!check.result) throw new Error(check.data);
            var action = JSON.parse(check.data);
            var to1 = action.to;
            var from = action.from;
            var quantity = action.quantity;

            if(from != se["eos-mulAccount"])  throw new Error('transcation rejected,trxid:' + trxid + ",from is " + from);
            let{to,value,symbol} = await checkDATXWithdraw(trxid);
            var checkResult = (to == to1 && quantity == value.toString() + ' ' + symbol.replace('D',''));
            if(!checkResult) throw new Error('transcation rejected,trxid:' + trxid + ",to is " + to + ",value is " + value);

            var result = await eosapi.confirm(proposer,proposeName,se["eos-account"],se["eos-privateKey"]);
            util.runDatxCmd("cldatx push action datxos.extra setdoing '[\"" + trxid + "\",\"" + nodeName + "\",\"" + se["producer-name"] + "\"]' -p " + se["producer-name"]);
            redis.client.set(trxid,"done");
            util.log('eos confirmed,datx trxid:' + trxid);
            response.writeHead(200,{'Content-Type': 'application/json'});
            response.write(JSON.stringify(result));
            response.end();
          }catch(e){
            util.log('Error:' + e.toString());
            rep(response,500,'Error:' + e.toString());
          }
          })();
        } 
      }
    }
  }
  catch (e) {
    util.log('Error:' + e.toString());
    rep(response,500,'Error:' + e.toString());
  }
}

function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
};

async function getVerifiers(){
  var verifyNodes = await httpClient.requestWithOvertime('http://' + se["http-server-address"] + '/v1/chain/get_table_rows',3000,'POST',
    '{"scope":"datxos.extra","code":"datxos.extra","table":"verifier","json":"true","limit":30}');
  verifyNodes = JSON.parse(verifyNodes.data).rows;
  return verifyNodes;
}

async function getProducers(){
  var produceNodes = await httpClient.requestWithOvertime('http://' + se["http-server-address"] + '/v1/chain/get_table_rows',3000,'POST',
    '{"scope":"datxos","code":"datxos","table":"producer","json":"true","limit":30}');
  produceNodes = JSON.parse(produceNodes.data).rows;
  return produceNodes;
}

async function btcGatherSign(trxSerialize,trxid,nodeName,sign,IsTestnet){
  returnVal = trxSerialize;
  var trx = bitcoinapi.getTrxFromHex(returnVal);
  var verifyNodes = await getVerifiers();
  for(var i = 0;i < verifyNodes.length;i++){
    if(verifyNodes[i].account == se["producer-name"]){
      trx = await bitcoinapi.signTrx(trx,se["btc-wif"],se["btc-mulscript"],{IsTestnet:IsTestnet})
      returnVal = trx.toHex();
      continue;
    };
    var URL = verifyNodes[i].url + '/btc/signTrx?trxid=' + trxid + '&trxSerialize=' + returnVal + '&isTestnet=' + IsTestnet +  '&nodeName=' + nodeName + '&sign=' + sign;
    util.log('btc request others to sign,url:' + URL);
    var result = await httpClient.requestWithOvertime(URL,10000,'GET');
    if(!result.result)util.log(result.errmsg);
    if(result.result){               
      //check return trx
      var isSame = false;

      var returnTrx = bitcoinapi.getTrxFromHex(result.data);
      var originalTrx = bitcoinapi.getTrxFromHex(returnVal);

      var returnDecodeInput = bitcoinapi.decodeInput(returnTrx);
      var originalDecodeInput = bitcoinapi.decodeInput(originalTrx);

      var returnDecodeOutput = bitcoinapi.decodeOutput(returnTrx,{IsTestnet:IsTestnet});
      var originalDecodeOutput = bitcoinapi.decodeOutput(originalTrx,{IsTestnet:IsTestnet});
      if(returnDecodeInput.length == originalDecodeInput.length 
        && returnDecodeOutput.length == originalDecodeOutput.length){
        isSame = true;
        for(var j = 0;j < returnDecodeInput.length;j ++){
          if(returnDecodeInput[j].txid != originalDecodeInput[j].txid
            || returnDecodeInput[j].n != originalDecodeInput[j].n){
              isSame = false;
            }
        }

        for(var k = 0;k < returnDecodeOutput.length;k++){
          if(returnDecodeOutput[k].scriptPubKey.hex != returnDecodeOutput[k].scriptPubKey.hex
            || returnDecodeOutput[k].satoshi != returnDecodeOutput[k].satoshi){
              isSame = false;
            }
        }
        if(isSame){
          returnVal = result.data;
        }
      } 
    }           
  }
  return returnVal;
}

async function checkDATXWithdraw(trxid){
  try
  {
    let to = "";
    let value = 0;
    let symbol = "";
    var URL ='http://' + se["http-server-address"] + '/v1/history/get_transaction';
    var result = await httpClient.requestWithOvertime(URL,5000,'POST','{"id":"' + trxid +'"}');

    if(result.result){
      var trx = JSON.parse(result.data);
      let quantity = trx.trx.trx.actions[0].data.quantity;
      to = trx.trx.trx.actions[0].data.memo;
      value = parseFloat(quantity.split(' ')[0]).toFixed(4);
      symbol = quantity.split(' ')[1];
    }
    return {to:to,value:value,symbol:symbol};
  }
  catch(e)
  {
    return {to:"",value:"",symbol:""};
  }
}

async function checkRequestSIgn(nodeName,sign,data){
  return true;
  if(!nodeName) return false;

  var verifyNodes = await getVerifiers();
  for(var i = 0;i < verifyNodes.length; i++){
    if(nodeName == verifyNodes[i].account){
      var pubkey = verifyNodes[i].verifier_key;
      pubkey = 'EOS' + pubkey.substring(4);
      return eosapi.EccVerify(sign,data,pubkey);
    }
  }

  var produceNodes = await getProducers();
  for(var i = 0;i < produceNodes.length; i++){
    if(nodeName == produceNodes[i].account){
      var pubkey = produceNodes[i].producer_key;
      return eosapi.EccVerify(sign,data,pubkey);
    }
  }
  return false;
}

async function checkExtractRight(nodeName,trxid){
  try{
    //return true;
    var records = await httpClient.requestWithOvertime('http://' + se["http-server-address"] + '/v1/chain/get_table_rows',3000,'POST',
    '{"scope":"datxos.extra","code":"datxos.extra","table":"record","json":"true","limit":3000}');
    records = JSON.parse(records.data).rows;

    for(let i = 0;i<records.length;i++){
      if(records[i].trxid == trxid){
        return nodeName == records[i].producer;
      }
    }
    return false;
  }
  catch(e){
    return false;
  }
}

function rep(response,head,msg){
  response.writeHead(head);
  response.write(msg);
  response.end();
}

module.exports = {
  handler
}




