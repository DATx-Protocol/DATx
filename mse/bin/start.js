const cluster = require('cluster');
const numCPUs = require('os').cpus().length;
const fs = require('fs');
const path = require('path');
const http2 = require('../lib');
const url = require('url');
const handler = require('../lib/handler').handler;
const ethapi = require('../lib/ethapi.js');

const INI = require("../lib/ini-file-loader.js");
const httpClient = require('../lib/client');

// if(cluster.isMaster)
// {
//     INI.getConfigFile();
//     for (let i = 0; i < numCPUs; i++) {
//         cluster.fork();
//     }
  
//     cluster.on('disconnect', (worker) => {
//       console.error('disconnect!');
//       cluster.fork();
//     });
// }
// else
// {
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

                //cluster.worker.disconnect();

                res.statusCode = 500;
                res.setHeader('content-type', 'text/plain');
                res.end('Oops, there was a problem!\n');
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


    //监听以太坊合约提币成功消息
    ethContractInstance = ethapi.getContractInstanceWithSocket();
    ethContractInstance.events.MultiTransact({fromBlock: 0, toBlock: 'latest'},function(error,result){
        if(result != undefined){
            let body = {
                category : 'ETH',
                transactionid : result.transactionHash,
                from : se["eth-muladdress"],
                to : result.returnValues.to,
                amount : result.returnValues.value,
                time : null,
                blocknum : result.blockNumber,
                isirreversible : false,
                memo : ethapi.toAscii(result.returnValues.data)
            };
            httpClient.requestAsync(se['listen-server-endpoint'] + '/eth_extract','POST',body);
        }
    });
// }


