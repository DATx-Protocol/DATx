var redis = require("redis");
var client = redis.createClient(); //127.0.0.1:6379

module.exports = {
    client
}

