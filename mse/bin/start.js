const fs = require('fs');
const path = require('path');
const http2 = require('../lib');
const url = require('url');
const handler = require('../lib/handler').handler;
const ethapi = require('../lib/ethapi.js');

const INI = require("../lib/ini-file-loader.js");
const httpClient = require('../lib/client');

const domain = require('domain');

// Creating the server in TLS mode
var server = http2.createServer(
{
    key: fs.readFileSync(path.join(__dirname, '../config/localhost.key')),
    cert: fs.readFileSync(path.join(__dirname, '../config/localhost.crt'))
},
(req,res) =>{
    const d = domain.create();

    d.on('error',(er) =>{
        console.error(`error ${er.stack}`);
        try
        {
            const killtimer = setTimeout(() => {
                process.exit(1);
            }, 2000);
            killtimer.unref();
            
            server.close();
            // res.statusCode = 500;
            // res.setHeader('content-type', 'text/plain');
            // res.end('Oops, there was a problem!\n');
        }catch (er2) {
        console.error(`Error sending 500! ${er2.stack}`);
    }
    });

    d.add(req);
    d.add(res);

    // Now run the handler function in the domain.
    d.run(() => {
        handler(req, res);
    });

});

let se = INI.getConfigFile();
let listenPoint = se["across-chain-endpoint"];
server.listen(url.parse(listenPoint).port); 


//listen ethereum withdraw event
ethContractInstance = ethapi.getContractInstanceWithSocket();
ethContractInstance.events.MultiTransact({fromBlock: 0, toBlock: 'latest'},function(error,result){
    if(result != undefined){
        let body = {
            category : 'ETH',
            transactionid : result.transactionHash,
            from : se["eth-muladdress"],
            to : result.returnValues.to,
            amount : parseFloat((parseFloat(result.returnValues.value)/1e18).toFixed(4)),
            time : null,
            blocknum : result.blockNumber,
            isirreversible : false,
            memo : ethapi.toAscii(result.returnValues.data)
        };
        console.log(JSON.stringify(body));
        httpClient.requestAsync(se['listen-server-endpoint'] + '/eth_extract','POST',JSON.stringify(body),function (error, response, body) {
            if (error) console.log('ERROR:' + error.toString());
            else console.log('SUCCESS:' + body);
          });
    }else{
        console.log('ERR:' + error.toString());
    }
});



