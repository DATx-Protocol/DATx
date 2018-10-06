var spawn = require('child_process').spawn;

function start(){
    nw = new spawn('node',['./lib/server.js']);
    nw.on('close',function(code, signal){
        nw.kill(signal);
        nw = start();
    });
    nw.on('error',function(code, signal){
        nw.kill(signal);
        nw = start();
    });
    return nw;
}
start();