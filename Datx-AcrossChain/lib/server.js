var fs = require('fs');
var path = require('path');
var http2 = require('..');
var url = require('url');
var querystring = require('querystring')

var bitcoinapi = require('./bitcoinapi.js')
var ethapi = require('./ethapi.js')
var eosapi = require('./eosapi.js')

var path = require('path');
var INI = require("../lib/ini-file-loader");
var confPath = path.resolve(__dirname, '../config/config.ini');
var ini___ = INI.loadFileSync(confPath);
var se = ini___.getOrCreateSection("node config");

var httpClient = require('./client.js')

// The callback to handle requests
function onRequest(request, response) {
  try {
    var pathName = url.parse(request.url).pathname;

    //genKeyPairs
    //params isTestnet = 1
    if (pathName == '/btc/genKeyPairs'){
      var IsTestnet = 0;
      var query = url.parse(request.url).query;
      if (query) {
        var params = querystring.parse(query)
        if(params.isTestnet == 1){
          IsTestnet = 1;
        }
      }
      var result = bitcoinapi.genKeyPairs({IsTestnet:IsTestnet})
      response.writeHead(200, {'Content-Type': 'application/json'});
      response.write(JSON.stringify(result));
      response.end();
    }

    //getKeysFromWIF
    //params isTestnet=1&wif=asdasdad
    else if (pathName == '/btc/getKeysFromWIF'){
      var IsTestnet = 0;
      var wif = '';
      var query = url.parse(request.url).query;
      if (query) {
        var params = querystring.parse(query)
        if(params.isTestnet == 1){
          IsTestnet = 1;
        }
        wif = params.wif;
        if(!wif){
          throw new Error('invalid params');
        }
        else{
          var result = bitcoinapi.getKeysFromWIF(wif,{IsTestnet:IsTestnet})
          response.writeHead(200, {'Content-Type': 'application/json'});
          response.write(JSON.stringify(result));
        }
        response.end();
      }
    }

    //genP2PKHAddr
    //params isTestnet=1&pubkey=adadada
    else if (pathName == '/btc/genP2PKHAddr') {
      var IsTestnet = 0;
      var pubkey = '';
      var query = url.parse(request.url).query;
      if (query) {
        var params = querystring.parse(query)
        if(params.isTestnet == 1){
          IsTestnet = 1;
        }
        pubkey = params.pubkey;
        if(!pubkey){
          throw new Error('invalid params');
        }
        else{
          var result = bitcoinapi.genP2PKHAddr(pubkey,{IsTestnet:IsTestnet})
          response.writeHead(200);
          response.write(result);
        }
      }
      response.end();
    }
    //genMulSigAddr
    //params isTestnet=1&pubkeys=asdasd,asdadgg&num=2
    else if(pathName == '/btc/genMulSigAddr'){
      var IsTestnet = 0;
      var pubkeys = [];
      var num = 0;
      var query = url.parse(request.url).query;
      if (query) {
        var params = querystring.parse(query)
        if(params.isTestnet == 1){
          IsTestnet = 1;
        }
        pubkeys = params.pubkeys.split(',');
        num = Number(params.num);
      }
      if(pubkeys.length <=0 || !num){
        throw new Error('invalid params');
      }
      else{
        var result = bitcoinapi.genMulSigAddr(pubkeys,num,{IsTestnet:IsTestnet})
        response.writeHead(200, {'Content-Type': 'application/json'});
        response.write(JSON.stringify(result));
      }
      response.end();
    }

    //withdraw
    //isTestnet=1&to=ffffff&value=10000000&fee=100000
    else if (pathName == '/btc/withdraw'){
      var IsTestnet = 0;
      var from = se["btc-muladdress"];
      var trxid = '';
      var to = '';
      var value = 0;
      var fee = 0;
      var nodeName = '';
      var sign = '';
      var script = se["btc-mulscript"];
      var query = url.parse(request.url).query;
      if(query){
        var params = querystring.parse(query)
        if(params.isTestnet == 1){
          IsTestnet = 1;
        }
        trxid = params.trxid;
        to = params.to;
        value = params.value;
        fee = params.fee;
        nodeName = params.nodeName;
        sgin = params.sign;
        if(!trxid|| !to || !value){
          throw new Error('invalid params');
        }
        else{
          (async function(){
            var checkData = to + value + fee + trxid + isTestnet ;
            var checkResult = await checkRequestSIgn(nodeName,sgin,checkData);
            if(!checkResult){
              throw new Error('invalid sign');
            }
            checkResult = await checkDATXWithdraw(trxid,to,value)
            if(!checkResult){
              throw new Error('invalid trx id');
            }
            //check 提币权


            //buildTrx
            var trx = await bitcoinapi.buildTrx(from,to,value,fee,{IsTestnet:IsTestnet}); 

            //call other nodes
            var verifyNodes = await httpClient.requestWithOvertime(se["datx-endpoint"] + '/v1/chain/get_table_rows',3000,'POST',
            '{"scope":"eosio","code":"eosio","table":"verifiers","json":"true","limit":30}');
            verifyNodes = JSON.parse(verifyNodes).rows;
            var trxSerialize = trx.toHex();
            for(var i = 0;i < verifyNodes.length;i++){
              if(verifyNodes[i].owner == se["node-name"]){
                trx = await bitcoinapi.signTrx(trx,se["btc-wif"],script,{IsTestnet:IsTestnet})
                trxSerialize = trx.toHex();
                continue;
              };
              checkData = trxSerialize + trxid + isTestnet; 
              sgin = eosapi.EccSIgn(checkData,se[datx-privateKey]);
              var URL = verifyNodes[i].url + '/btc/signTrx?trxid=' + trxid + '&trx=' + trxSerialize + '&isTestnet=' + IsTestnet +  '&nodeName=' + se[node-name] + '&sign=' + sign;
              var result = await httpClient.requestWithOvertime(URL,3000,'GET');
              if(result.result){               
                //check return trx
                var isSame = false;

                var returnTrx = bitcoinapi.getTrxFromHex(result.data);
                var originalTrx = bitcoinapi.getTrxFromHex(trxSerialize);

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
                    trxSerialize = result.data;
                  }
                } 
              }           
            }

            //broadcast
            trx = bitcoinapi.getTrxFromHex(trxSerialize);
            result = await bitcoinapi.broadcastTrx(trx); 
            redis.client.set(trxid,trxSerialize);

            response.writeHead(200);
            response.write(result);
            response.end();
          })()
        }
      }
    }

    //signTrx
    //isTestnet=1&trx=asdafafff
    //need other params to judge if should sign
    else if(pathName == '/btc/signTrx'){
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
        trxSerialize = params.trx;
        trxid = params.trxid;
        nodeName = params.nodeName;
        sign = params.sign;
        if(!trxSerialize || !trxid){
          throw new Error('invalid params');
        }
        else{
          (async function(){
            var checkData = trxSerialize + trxid + isTestnet; 
            var checkResult = await checkRequestSIgn(nodeName,sign,checkData);
            if(!checkResult){
              throw new Error('invalid sign');
            }

            var recordTrxSerialize = redis.client.get(trxid);
            var trx = bitcoinapi.getTrxFromHex(trxSerialize);
            var ins = bitcoinapi.decodeInput(trx);
            var outs = bitcoinapi.decodeOutput(trx,{IsTestnet:IsTestnet});
            //do some check
            if(recordTrxSerialize == null || recordTrxSerialize == undefined || recordTrxSerialize == ""){
              var sum = 0;
              var input;
              var refTrx;
              for (var i =0; i < ins.length; i++){
                input = ins[i];
                refTrx = await bitcoinapi.getTrxDetail(input.txid,{IsTestnet:IsTestnet});
                if(refTrx.result == false) throw new Error('transcation rejected');
                refTrx = JSON.parse(refTrx.data);
                if(refTrx.outs[input.n].addr != se["btc-muladdress"]) throw new Error('transcation rejected');
                sum += Number(refTrx.outs[input.n].value);
              }

              if(outs.length != 2) throw new Error('transcation rejected');
              if(outs[1].addresses[0] != se["btc-muladdress"])  throw new Error('transcation rejected');
              if(outs[2].addresses[0] != trxid) throw new Error('transcation rejected');

              var to = outs[0].addresses[0];
              var fee = sum - outs[0].satoshi - outs[1].satoshi;
              var value = outs[0].satoshi;
              if(fee > Number(se["btc-maxfee"])) throw new Error('transcation rejected');
              checkResult = await checkDATXWithdraw(trxid,to,value+fee);
              if(!checkResult)  throw new Error('transcation rejected');
            }
            else{
              var recordTrx = bitcoinapi.getTrxFromHex(recordTrxSerialize);
              var recordIns = bitcoinapi.decodeInput(recordTrx);
              var recordOuts = bitcoinapi.decodeOutput(recordTrx,{IsTestnet:IsTestnet});
              if(recordIns.length != ins.length || recordOuts.length != outs.length) throw new Error('transcation rejected');

              for(var i =0; i < ins.length; i++){
                if(ins[i].txid != recordIns[i].txid ||
                   ins[i].n != recordIns[i].n ||
                   ins[i].sequence != recordIns[i].sequence)
                   throw new Error('transcation rejected');
              }

              for(var i =0; i < outs.length; i++){
                if(outs[i].scriptPubKey.hex != recordOuts[i].scriptPubKey.hex ||
                   outs[i].scriptPubKey.satoshi != recordOuts[i].scriptPubKey.satoshi)
                   throw new Error('transcation rejected');
              }
            }
            trx = await bitcoinapi.signTrx(trx,se["btc-wif"],script,{IsTestnet:IsTestnet});
            trxSerialize = trx.toHex();
            redis.client.set(trxid,trxSerialize);

            response.writeHead(200);
            response.write(trxSerialize);
            response.end();

            //call other nodes
            var verifyNodes = await httpClient.requestWithOvertime(se["datx-endpoint"] + '/v1/chain/get_table_rows',3000,'POST',
            '{"scope":"eosio","code":"eosio","table":"verifiers","json":"true","limit":30}');
            verifyNodes = JSON.parse(verifyNodes).rows;
            for(var i = 0;i < verifyNodes.length;i++){
              if(verifyNodes[i].owner == se["node-name"]){
                continue;
              };
              checkData = trxSerialize + trxid + isTestnet; 
              sgin = eosapi.EccSIgn(checkData,se[datx-privateKey]);
              var URL = verifyNodes[i].url + '/btc/signTrx?trxid=' + trxid + '&trx=' + trxSerialize + '&isTestnet=' + IsTestnet +  '&nodeName=' + se[node-name] + '&sign=' + sign;
              var result = await httpClient.requestWithOvertime(URL,3000,'GET');
              if(result.result){               
                //check return trx
                var isSame = false;

                var returnTrx = bitcoinapi.getTrxFromHex(result.data);
                var originalTrx = bitcoinapi.getTrxFromHex(trxSerialize);

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
                    trxSerialize = result.data;
                  }
                } 
              }           
            }
            //broadcast
            trx = bitcoinapi.getTrxFromHex(trxSerialize);
            result = await bitcoinapi.broadcastTrx(trx); 
            redis.client.set(trxid,trxSerialize);
          })();
        }
      }
    }

    //eth withdraw
    else if(pathName == '/eth/withdraw'){
      var myAddr = se["eth-myaddress"];
      var contractAddr = se["eth-muladdress"];
      var myPrivatekey = se["eth-privateKey"];
      var to = '';
      var value = '';
      var trxid = '';
      var sign = '';
      var nodeName  = '';
      var isInform = false;
      var query = url.parse(request.url).query;
      if(query){
        var params = querystring.parse(query)
        to = params.to;
        value = Number(params.value);
        trxid = params.trxid;
        sign = params.sign;
        nodeName = params.nodeName;
        isInform = params.isInform;
        if(!to || !value){
          throw new Error('invalid params');
        }
        else{
          (async function(){
            //do some check
            var checkData = to + value + trxid;
            var checkResult = await checkRequestSIgn(nodeName,sgin,checkData);
            if(!checkResult){
              throw new Error('invalid sign');
            }
            if(redis.client.get(trxid) == "done"){
              throw new Error('handled trx id');
            }

            checkResult = await checkDATXWithdraw(trxid,to,value)
            if(!checkResult)  throw new Error('transcation rejected');

            //通知其他节点configm
            if(!isInform){
              var verifyNodes = await httpClient.requestWithOvertime(se["datx-endpoint"] + '/v1/chain/get_verifier_schedule',3000,'POST','');
              verifyNodes = JSON.parse(verifyNodes).rows;
              for(var i = 0;i < verifyNodes.length;i++){
                if(verifyNodes[i].owner == se["node-name"]){
                  ethapi.withdraw(myAddr,contractAddr,to,value,myPrivatekey,ethapi.fromAscii(trxid));
                  continue;
                }

                var URL = verifyNodes[i].url + '/eth/withdraw?trxid=' + trxid + '&to=' + to + '&value=' + value + '&isInform=true&sign=' + sign + '&nodeName=' + nodeName;
                httpClient.requestAsync(URL);
              }
            }
            else{
              ethapi.withdraw(myAddr,contractAddr,to,value,myPrivatekey,ethapi.fromAscii(trxid));
            }

            redis.client.set(trxid,"done");
            response.writeHead(200);
            response.write('success');
            response.end();
          })();
        }
      }
    }

    //createAccount
    else if(pathName == '/eos/createAccount'){
      var creator = '';
      var creatorKey = '';
      var accountName = '';
      var auth = '';
      var query = url.parse(request.url).query;
      if(query){
        var params = querystring.parse(query)
        creator = params.creator;
        creatorKey = params.creatorKey;
        accountName = params.accountName;
        auth = JSON.parse(auth);
        if(!creator || !creatorKey || !accountName || !auth){
          throw new Error('invalid params');
        }
        else{
          (async function(){
            var result = await eosapi.createAccount(creator,creatorKey,accountName,auth);

            response.writeHead(200,{'Content-Type': 'application/json'});
            response.write(JSON.stringify(result));
            response.end();
          })();
        }
      }
    }

    //issue
    else if(pathName == '/eos/issue'){
      var to = '';
      var value = '';
      var memo = '';
      query = url.parse(request.url).query;
      if(query){
        var params = querystring.parse(query);
        to = params.to;
        value = params.value;
        memo = params.memo;
        if(!to || !value){
          throw new Error('invalid params');
        }
        else{
          (async function(){
            var result = await eosapi.issue(to,value,memo);

            response.writeHead(200,{'Content-Type': 'application/json'});
            response.write(JSON.stringify(result));
            response.end();
          })();
          
        }
      }
    }

    else if(pathName == '/eos/getAccount'){
      response.writeHead(200);
      response.write(se["eos-account"]);
      response.end();
    }

    //propose
    else if(pathName == '/eos/withdraw'){
      var proposer = se["eos-account"]; 
      var ProvidedKey = se["eos-privateKey"]
      var auths = []; //根据见证节点组合成的权限要求
      var MultiSigAccount = se["eos-mulAccount"]; 
      var to = '';
      var value = '';
      var trxid = '';
      var sign = '';
      var nodeName = '';
      var query = url.parse(request.url).query;
      if(query){
        var params = querystring.parse(query);
        to = params.to;
        value = params.value;
        trxid = params.trxid;
        sign = params.sign;
        nodeName = params.nodeName;
        if(!to || !value || !trxid){
          throw new Error('invalid params');
        }
        else{
          (async function(){
            //chenck 
            var checkData = to + value + trxid;
            var checkResult = await checkRequestSIgn(nodeName,sgin,checkData);
            if(!checkResult){
              throw new Error('invalid sign');
            }
            checkResult = await checkDATXWithdraw(trxid,to,value)
            if(!checkResult)  throw new Error('transcation rejected');

            var verifyNodes = await httpClient.requestWithOvertime(se["datx-endpoint"] + '/v1/chain/get_verifier_schedule',3000,'POST','');
            verifyNodes = JSON.parse(verifyNodes).rows;
            for (var i = 0; i < verifyNodes.length; i++){
              if(verifyNodes[i].owner == se["node-name"]){
                auths.push({actor: se["eos-account"], permission: 'active'});
                isVerifier = true;
                continue;
              }
              URL = verifyNodes[i].url + '/eos/getAccount';
              var result = await httpClient.requestWithOvertime(URL,3000,'GET');
              
              if(result.result){
                auths.push({actor: result.data, permission: 'active'});
              }
            }
      
            var proposeName = await eosapi.propose(proposer,MultiSigAccount,to,value,trxid,auths,ProvidedKey);
            //broadcast to other nodes
            for (var i = 0; i < verifyNodes.length; i++){
              if(verifyNodes[i].owner == se["node-name"]){
                eosapi.confirm(proposer,proposeName,proposer,ProvidedKey);
                continue;
              }

              checkData = proposer + proposeName + trxid; 
              sgin = eosapi.EccSIgn(checkData,se[datx-privateKey]);
              var URL = verifyNodes[i].url + '/eos/confirm?proposer=' + proposer + '&proposeName=' + proposeName + '&trxid=' + trxid + '&nodeName=' + se[node-name] + '&sign=' + sign;;
              httpClient.requestAsync(URL);
            }
            //wait 3 second then exec
            await sleep(3000);

            //exec
            result = await eosapi.exec(proposer,proposeName,proposer,ProvidedKey);
            
            response.writeHead(200);
            if(result && result.transaction_id){
              handledTrx.add(trxid);
              response.write(result.transaction_id);
            }
            else{
              response.write(result);
            }
            response.end();
          })();
        }
      }
    }

    //confirm
    else if(pathName == '/eos/confirm'){
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
          throw new Error('invalid params');
        }
        else{
          (async function(){
            //do some check
            checkData = proposer + proposeName + trxid; 
            checkResult = await checkRequestSIgn(nodeName,sgin,checkData);
            if(!checkResult){
              throw new Error('invalid sign');
            }
            if(redis.client.get(trxid) == "done"){
              throw new Error('handled trx id');
            }

            var check = await eosapi.getProposeAction(proposer,proposeName);
            if(!check.result) throw new Error(check.data);
            var action = JSON.parse(check.data);
            var to = action.to;
            var from = action.from;
            var quantity = action.quantity;

            if(from != se["eos-mulAccount"])  throw new Error('transcation rejected');
            var checkResult = await checkDATXWithdraw(trxid,to,quantity)
            if(!checkResult)  throw new Error('transcation rejected');

            var result = await eosapi.confirm(proposer,proposeName,se["eos-account"],se["eos-privateKey"]);

            redis.client.set(trxid,"done");
            response.writeHead(200,{'Content-Type': 'application/json'});
            response.write(JSON.stringify(result));
            response.end();
          })();
        } 
      }
    }
  }
  catch (e) {
    response.writeHead(200);
    response.write('Error:' + e.toString());
    response.end();
  }
}

function sleep(time) {
  return new Promise(function (resolve, reject) {
      setTimeout(function () {
          resolve();
      }, time);
  })
};

async function checkDATXWithdraw(trxid,to,value){
  try{
    var URL = se["datx-endpoint"] + '/v1/history/get_transaction';
    var result = await httpClient.requestWithOvertime(URL,5000,'POST','{"id":"' + trxid +'"}');
    result = json.parse(result);
    if(result.result){
      var trx = JSON.parse(result.data);
      if(trx.trx.trx.actions[0].data.quantity == value){
        if(trx.trx.trx.actions[0].data.memo == to){
          return true;
        }
      }
    }
    return false;
  }
  catch(e){
    return false
  }
}

async function checkRequestSIgn(nodeName,sign,data){
  if(!nodeName) return false;

  var verifyNodes = await httpClient.requestWithOvertime(se["datx-endpoint"] + '/v1/chain/get_verifier_schedule',3000,'POST','');
  verifyNodes = JSON.parse(verifyNodes).rows;
  for(var i = 0;i < verifyNodes.length; i++){
    if(nodeName == verifyNodes[i].owner){
      var pubkey = verifyNodes[i].verifier_key;
      return eosapi.EccVerify(sign,data,pubkey);
    }
  }
  return false;
}

// Creating a bunyan logger (optional)
var log = require('../test/util').createLogger('server');

// Creating the server in plain or TLS mode (TLS mode is the default)
var server;
if (process.env.HTTP2_PLAIN) {
  server = http2.raw.createServer({log: log}, onRequest);
} else {
  server = http2.createServer(
      {
        log: log,
        key: fs.readFileSync(path.join(__dirname, '/localhost.key')),
        cert: fs.readFileSync(path.join(__dirname, '/localhost.crt'))
      },
      onRequest);
}
server.listen(process.env.HTTP2_PORT || 8080);

//监听以太坊合约提币成功消息
ethContractInstance = ethapi.getContractInstanceWithSocket();
ethContractInstance.events.MultiTransact({fromBlock: 0, toBlock: 'latest'},function(error,result){
  if(result != undefined){
    console.log(ethapi.toAscii(result.returnValues.data));//0x6173646164616461   //asdadada
    console.log(result.transactionHash); //0x3a903723c6a31eda9616228d1c0343f45bbd9e3d11dcb079fd7a3122f74466cd
  }
});




