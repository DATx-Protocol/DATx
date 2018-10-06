
var path = require('path');
var INI = require("../lib/ini-file-loader");
var confPath = path.resolve(__dirname, '../config/config.ini');
var ini___ = INI.loadFileSync(confPath);
var se = ini___.getOrCreateSection("node config");

var url = require('url');
var request = require("request");

function requestWithOvertime(url, interval,method,body) {
  return new Promise(function(resolve, reject) {
    var options = { method: method ||'GET',
                    url: url,
                    insecure: true,
                    rejectUnauthorized: false};
    if(body){
      options.body = body;
    }

    request(options, function (error, response, body) {
      if (error) resolve({result: false, errmsg: error.message});
      resolve({result: true, data: body});
    });

    setTimeout(_=>{
      resolve({result: false, errmsg: 'request time out'});
     }, interval);
  });
}

function requestAsync(url,method,body,callback){
  var options = { method: method ||'GET',
                  url: url,
                  insecure: true,
                  rejectUnauthorized: false};
  if(body){
    options.body = body;
  }
  request(options,callback);
}


module.exports = {
  requestWithOvertime,
  requestAsync
}