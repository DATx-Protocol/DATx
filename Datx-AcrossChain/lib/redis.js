var redis = require("redis");
const {promisify} = require('util');

var client = redis.createClient(); //127.0.0.1:6379

const getAsync = promisify(client.get).bind(client);
const sismemberAsync = promisify(client.sismember).bind(client);
module.exports = {
    client,
    getAsync,
    sismemberAsync
}

