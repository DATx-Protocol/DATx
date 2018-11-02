const exec = require('child_process').exec;
const INI = require("../lib/ini-file-loader.js");


//需要执行的命令字符串
var cli = 'cldatx wallet keys';

exec(cli,{encoding:'utf8'},function (err,stdout,stderr){
    if (err){
        if(err.message.indexOf('Locked wallet') !== -1){
          let {wname,wpassword} = INI.getWallet();
          exec('cldatx wallet unlock -n ' + wname + ' --password ' + wpassword,{encoding:'utf8'},function (err,stdout,stderr){
            if(err){
                return err.message;
            }else{
                exec(cli,{encoding:'utf8'},function (err,stdout,stderr){
                    if(err){
                        return err.message;
                    }
                    else{
                        return stdout;
                    }
                });
            }
          });
        }else{
            return err.message;
        }
    }else{
        return stdout;
    }
    
})