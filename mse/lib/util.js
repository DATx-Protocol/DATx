const exec = require('child_process').exec;

Date.prototype.format = function (fmt) {
    var o = {
        "M+": this.getMonth() + 1, //月份
        "d+": this.getDate(), //日
        "h+": this.getHours(), //小时
        "m+": this.getMinutes(), //分
        "s+": this.getSeconds(), //秒
        "q+": Math.floor((this.getMonth() + 3) / 3), //季度
        "S": this.getMilliseconds() //毫秒
    };
    if (/(y+)/.test(fmt)) {
      fmt = fmt.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length));
    }
    for (var k in o) {
      if (new RegExp("(" + k + ")").test(fmt)) {
        fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ?
          (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
      }
    }
    return fmt;
  }

function log(msg){
    console.log(new Date().format("yyyy-MM-dd hh:mm:ss.S") + "====" + msg);
  }


function runDatxCmd(cli){
    exec(cli,{encoding:'utf8'},function (err,stdout,stderr){
      if (err){
          if(err.message.indexOf('Locked wallet') !== -1){
            let {wname,wpassword} = INI.getWallet();
            exec('cldatx wallet unlock -n ' + wname + ' --password ' + wpassword,{encoding:'utf8'},function (err,stdout,stderr){
              if(err){
                  console.log(err.message);
                  return err.message;
              }else{
                  exec(cli,{encoding:'utf8'},function (err,stdout,stderr){
                      if(err){
                          console.log(err.message);
                          return err.message;
                      }
                      else{
                          return stdout;
                      }
                  });
              }
            });
          }else{
              console.log(err.message);
              return err.message;
          }
      }else{
          return stdout;
      }
    })
  }


module.exports = {
    log,
    runDatxCmd
  }
  