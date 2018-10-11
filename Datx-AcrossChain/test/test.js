redis = require("../lib/redis.js");

(async function(){

    if(await redis.sismemberAsync('qqqq','bbbbbb') == 11){
    console.log('success');
    }
})();

