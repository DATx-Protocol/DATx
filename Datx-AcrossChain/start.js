const cluster = require('cluster');
const numCPUs = require('os').cpus().length;
const fs = require('fs');
const path = require('path');
const http2 = require('./lib');
const handler = require('./lib/handler').handler;
const ethapi = require('./lib/ethapi.js');

// if(cluster.isMaster)
// {
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
//     const domain = require('domain');

//     // Creating the server in TLS mode
//     var server = http2.createServer(
//     {
//         key: fs.readFileSync(path.join(__dirname, './lib/localhost.key')),
//         cert: fs.readFileSync(path.join(__dirname, './lib/localhost.crt'))
//     },
//     (req,res) =>{
//         const d = domain.create();

//         d.on('error',(er) =>{
//             console.error(`error ${er.stack}`);
//             try
//             {
//                 const killtimer = setTimeout(() => {
//                     process.exit(1);
//                 }, 30000);
//                 killtimer.unref();
                
//                 server.close();

//                 cluster.worker.disconnect();

//                 res.statusCode = 500;
//                 res.setHeader('content-type', 'text/plain');
//                 res.end('Oops, there was a problem!\n');
//             }catch (er2) {
//             console.error(`Error sending 500! ${er2.stack}`);
//         }
//         });

//         d.add(req);
//         d.add(res);

//         // Now run the handler function in the domain.
//         d.run(() => {
//             handler(req, res);
//         });

//     });

//     server.listen(process.env.HTTP2_PORT || 8080); 


//     //监听以太坊合约提币成功消息
//     ethContractInstance = ethapi.getContractInstanceWithSocket();
//     ethContractInstance.events.MultiTransact({fromBlock: 0, toBlock: 'latest'},function(error,result){
//     if(result != undefined){
//         console.log(ethapi.toAscii(result.returnValues.data));//0x6173646164616461   //asdadada
//         console.log(result.transactionHash); //0x3a903723c6a31eda9616228d1c0343f45bbd9e3d11dcb079fd7a3122f74466cd
//     }
//     });
// }

    var server = http2.createServer(
    {
        key: fs.readFileSync(path.join(__dirname, './lib/localhost.key')),
        cert: fs.readFileSync(path.join(__dirname, './lib/localhost.crt'))
    },handler);
    server.listen(8080);