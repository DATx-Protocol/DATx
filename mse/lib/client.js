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
      else if (response.statusCode != 200) resolve({result: false, errmsg: body});
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